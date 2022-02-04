package common

import "fmt"

type Exception error

type ExceptionType string

const (
	SyntaxError       ExceptionType = "SyntaxError"
	DivideByZeroError ExceptionType = "DivideByZero"
	TypeError         ExceptionType = "TypeError"
	NameError         ExceptionType = "NameError"
)

func NewException(
	exceptionType ExceptionType,
	message string,
	line int,
	col int) Exception {
	return fmt.Errorf("%v: %v at %v:%v", exceptionType, message, line, col)
}
