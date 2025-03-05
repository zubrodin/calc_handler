package orchestrator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/zubrodin/calc_handler/internal/models"
)

var (
	expressions = make(map[string]*models.Expression)
	tasks       = make(chan models.Task, 100)
	mu          sync.Mutex
)

func AddExpression(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Expression == "" {
		http.Error(w, "Invalid data", http.StatusUnprocessableEntity)
		return
	}

	exp := models.Expression{
		ID:     uuid.New().String(),
		Status: "pending",
		Result: 0,
	}
	mu.Lock()
	expressions[exp.ID] = &exp
	mu.Unlock()

	go processExpression(req.Expression, exp.ID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": exp.ID})
}

func GetExpressions(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var expList []models.Expression
	for _, exp := range expressions {
		expList = append(expList, *exp)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"expressions": expList,
	})
}

func GetExpressionByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/expressions/")
	mu.Lock()
	defer mu.Unlock()

	exp, exists := expressions[id]
	if !exists {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"expression": exp,
	})
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	select {
	case task := <-tasks:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(task)
	default:
		http.Error(w, "No task available", http.StatusNotFound)
	}
}

func HandleResult(w http.ResponseWriter, r *http.Request) {
	var result models.Result
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, "Invalid result data", http.StatusUnprocessableEntity)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if exp, exists := expressions[result.TaskID]; exists {
		exp.Result += result.Result
		exp.Status = "completed"
	} else {
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "result recorded"})
}
func processExpression(expression string, expID string) {
	tokens := tokenize(expression)
	var stack []float64
	var ops []string

	fmt.Printf("Processing expression: %s\n", expression)

	for _, token := range tokens {
		if token == "(" {
			ops = append(ops, token)
		} else if token == ")" {
			for len(ops) > 0 && ops[len(ops)-1] != "(" {
				if len(stack) < 2 {
					mu.Lock()
					exp := expressions[expID]
					exp.Status = "failed"
					mu.Unlock()
					fmt.Println("Failed to process expression: not enough operands")
					return
				}

				arg2 := stack[len(stack)-1]
				arg1 := stack[len(stack)-2]
				stack = stack[:len(stack)-2]

				op := ops[len(ops)-1]
				ops = ops[:len(ops)-1]

				time := getOperationTime(op)
				fmt.Printf("Operation %s will take %d ms\n", op, time)

				result := performOperation(arg1, arg2, op)
				stack = append(stack, result)
			}
			if len(ops) > 0 {
				ops = ops[:len(ops)-1]
			}
		} else if isOperator(token) {
			ops = append(ops, token)
		} else {

			arg, err := strconv.ParseFloat(token, 64)
			if err != nil {
				mu.Lock()
				exp := expressions[expID]
				exp.Status = "failed"
				mu.Unlock()
				fmt.Println("Failed to parse token as float:", token)
				return
			}
			stack = append(stack, arg)
		}
	}

	for len(ops) > 0 {
		if len(stack) < 2 {
			mu.Lock()
			exp := expressions[expID]
			exp.Status = "failed"
			mu.Unlock()
			fmt.Println("Failed to process expression: not enough operands")
			return
		}

		arg2 := stack[len(stack)-1]
		arg1 := stack[len(stack)-2]
		stack = stack[:len(stack)-2]

		op := ops[len(ops)-1]
		ops = ops[:len(ops)-1]

		time := getOperationTime(op)
		fmt.Printf("Operation %s will take %d ms\n", op, time)

		result := performOperation(arg1, arg2, op)
		stack = append(stack, result)
	}

	mu.Lock()
	exp := expressions[expID]
	exp.Status = "completed"
	if len(stack) > 0 {
		exp.Result = stack[0]
	}
	mu.Unlock()
	fmt.Printf("Expression completed: ID=%s, Result=%f\n", expID, exp.Result)
}
func performOperation(arg1, arg2 float64, operator string) float64 {
	switch operator {
	case "+":
		return arg1 + arg2
	case "-":
		return arg1 - arg2
	case "*":
		return arg1 * arg2
	case "/":
		if arg2 == 0 {

			fmt.Println("Error: Division by zero")
			return 0
		}
		return arg1 / arg2
	default:
		return 0
	}
}

func getOperationTime(operation string) int {
	switch operation {
	case "+":
		return getEnvAsInt("TIME_ADDITION_MS")
	case "-":
		return getEnvAsInt("TIME_SUBTRACTION_MS")
	case "*":
		return getEnvAsInt("TIME_MULTIPLICATIONS_MS")
	case "/":
		return getEnvAsInt("TIME_DIVISIONS_MS")
	default:
		return 0
	}
}
func getEnvAsInt(env string) int {
	val, _ := os.LookupEnv(env)
	if intVal, err := strconv.Atoi(val); err == nil {
		return intVal
	}
	return 0
}
func isOperator(token string) bool {
	return token == "+" || token == "-" || token == "*" || token == "/"
}
func tokenize(expression string) []string {

	re := regexp.MustCompile(`\d+\.?\d*|[+\-*/()]|\s+`)
	matches := re.FindAllString(expression, -1)

	var tokens []string
	for _, match := range matches {
		if match != " " && match != "\n" {
			tokens = append(tokens, match)
		}
	}
	return tokens
}
