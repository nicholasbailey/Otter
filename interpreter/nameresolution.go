package interpreter

import (
	"fmt"

	"github.com/nicholasbailey/becca/exception"
	"github.com/nicholasbailey/becca/parser"
)

func (intepreter *Interpreter) resolveName(name *parser.Token) (*BeccaValue, exception.Exception) {
	val, found := intepreter.CallStack.ResolveVariable(name.Value)
	if found {
		return val, nil
	} else {
		return nil, exception.New(exception.NameError, fmt.Sprintf("%v is not defined", name.Value), name.Line, name.Col)
	}
}
