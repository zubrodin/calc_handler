package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type Task struct {
	ID            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}

type Result struct {
	ID     string  `json:"id"`
	Result float64 `json:"result"`
}

var computingPower int

func main() {

	var err error
	computingPower, err = strconv.Atoi(getEnv("COMPUTING_POWER", "4"))
	if err != nil {
		log.Fatalf("Invalid COMPUTING_POWER: %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < computingPower; i++ {
		wg.Add(1)
		go worker(&wg)
	}

	wg.Wait()
}

func worker(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		task, err := getTask()
		if err != nil {
			log.Println(err)
			time.Sleep(1 * time.Second)
			continue
		}

		result := performCalculation(task)
		if err := sendResult(task.ID, result); err != nil {
			log.Printf("Failed to send result for task %s: %v\n", task.ID, err)
		}
	}
}

func getTask() (*Task, error) {
	resp, err := http.Get("http://localhost/internal/task")
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("no task available: %s", resp.Status)
	}

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, fmt.Errorf("failed to decode task response: %v", err)
	}

	return &task, nil
}

func performCalculation(task *Task) float64 {
	time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)

	var result float64
	switch task.Operation {
	case "+":
		result = task.Arg1 + task.Arg2
	case "-":
		result = task.Arg1 - task.Arg2
	case "*":
		result = task.Arg1 * task.Arg2
	case "/":
		if task.Arg2 != 0 {
			result = task.Arg1 / task.Arg2
		} else {
			log.Println("Division by zero")
			result = 0
		}
	}

	return result
}

func sendResult(taskID string, result float64) error {
	url := "http://localhost/internal/task/result"
	resultData := Result{
		ID:     taskID,
		Result: result,
	}

	jsonData, err := json.Marshal(resultData)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send result: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to record result: %s", resp.Status)
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
