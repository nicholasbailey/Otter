package interpreter

import (
	"fmt"

	"github.com/nicholasbailey/becca/common"
	"github.com/nicholasbailey/becca/parser"
)

func (intepreter *Interpreter) resolveName(name *parser.Token) (*BeccaValue, error) {
	val, found := intepreter.CallStack.ResolveVariable(name.Value)
	if found {
		return val, nil
	} else {
		return nil, common.NewException(common.NameError, fmt.Sprintf("%v is not defined", name.Value), name.Line, name.Col)
	}
}
