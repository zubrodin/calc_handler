package calc

import "errors"

var (
	ErrEmptyExpression   = errors.New("empty expression or invalid request")
	ErrInvalidExpression = errors.New("invalid expression")
	ErrDivisionByZero    = errors.New("division by zero")
)
