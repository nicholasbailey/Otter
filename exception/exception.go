package exception

import "fmt"

type Exception error

type ExceptionType string

const (
	// Represents an issue with syntax. Generally
	// thrown by the parser, but
	AssertionError    ExceptionType = "AssertionError"
	SyntaxError       ExceptionType = "SyntaxError"
	DivideByZeroError ExceptionType = "DivideByZeroError"
	TypeError         ExceptionType = "TypeError"
	NameError         ExceptionType = "NameError"
	InternalError     ExceptionType = "InternalError"
	MethodError       ExceptionType = "MethodError"
	ArgumentError     ExceptionType = "ArgumentError"
	IndexError        ExceptionType = "IndexError"
	IterationError    ExceptionType = "IterationError"
)

func New(
	exceptionType ExceptionType,
	message string,
	line int,
	col int) Exception {
	return fmt.Errorf("%v: %v at %v:%v", exceptionType, message, line, col)
}
