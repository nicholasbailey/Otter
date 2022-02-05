package interpreter

import (
	"github.com/nicholasbailey/becca/common"
	"github.com/nicholasbailey/becca/parser"
)

type BeccaValue struct {
	Type     *BeccaValue
	Value    interface{}
	Callable *Callable
	Methods  map[string]*Callable
}

type BuiltInFunction func(interpreter *Interpreter, values []*BeccaValue) (*BeccaValue, common.Exception)

type Callable struct {
	Name                string
	Arity               int
	UserDefinedFunction *parser.Token
	BuiltInFunction     BuiltInFunction
}
