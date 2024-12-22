package calc

import (
	"strconv"
	"strings"
)

func Calc(expression string) (float64, error) {
	tokens := tokenize(expression)
	if len(tokens) == 0 {
		return 0, ErrEmptyExpression
	}

	rpn, err := toRPN(tokens)
	if err != nil {
		return 0, err
	}

	return evaluateRPN(rpn)
}

func tokenize(expression string) []string {
	expression = strings.ReplaceAll(expression, " ", "")
	var tokens []string
	var number strings.Builder

	for _, ch := range expression {
		switch ch {
		case '+', '-', '*', '/':
			if number.Len() > 0 {
				tokens = append(tokens, number.String())
				number.Reset()
			}
			tokens = append(tokens, string(ch))
		case '(', ')':
			if number.Len() > 0 {
				tokens = append(tokens, number.String())
				number.Reset()
			}
			tokens = append(tokens, string(ch))
		default:
			number.WriteRune(ch)
		}
	}
	if number.Len() > 0 {
		tokens = append(tokens, number.String())
	}

	return tokens
}

func toRPN(tokens []string) ([]string, error) {
	var output []string
	var stack []string

	for _, token := range tokens {
		switch token {
		case "+", "-":
			for len(stack) > 0 && (stack[len(stack)-1] == "+" || stack[len(stack)-1] == "-" || stack[len(stack)-1] == "*" || stack[len(stack)-1] == "/") {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		case "*", "/":
			for len(stack) > 0 && (stack[len(stack)-1] == "*" || stack[len(stack)-1] == "/") {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		case "(":
			stack = append(stack, token)
		case ")":
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return nil, ErrInvalidExpression
			}
			stack = stack[:len(stack)-1]
		default:
			output = append(output, token)
		}
	}

	for len(stack) > 0 {
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return output, nil
}

func evaluateRPN(rpn []string) (float64, error) {
	var stack []float64

	for _, token := range rpn {
		switch token {
		case "+":
			if len(stack) < 2 {
				return 0, ErrInvalidExpression
			}
			b := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			a := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			stack = append(stack, a+b)
		case "-":
			if len(stack) < 2 {
				return 0, ErrInvalidExpression
			}
			b := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			a := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			stack = append(stack, a-b)
		case "*":
			if len(stack) < 2 {
				return 0, ErrInvalidExpression
			}
			b := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			a := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			stack = append(stack, a*b)
		case "/":
			if len(stack) < 2 {
				return 0, ErrInvalidExpression
			}
			b := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			a := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if b == 0 {
				return 0, ErrDivisionByZero
			}
			stack = append(stack, a/b)
		default:
			num, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, ErrInvalidExpression
			}
			stack = append(stack, num)
		}
	}

	if len(stack) != 1 {
		return 0, ErrInvalidExpression
	}
	return stack[0], nil
}
