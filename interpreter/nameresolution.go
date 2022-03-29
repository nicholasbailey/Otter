package interpreter

import (
	"fmt"

	"github.com/nicholasbailey/otter/exception"
	"github.com/nicholasbailey/otter/parser"
)

func (intepreter *Interpreter) resolveName(name *parser.Token) (*OtterValue, exception.Exception) {
	val, found := intepreter.CallStack.ResolveVariable(name.Value)
	if found {
		return val, nil
	} else {
		return nil, exception.New(exception.NameError, fmt.Sprintf("%v is not defined", name.Value), name.Line, name.Col)
	}
}
