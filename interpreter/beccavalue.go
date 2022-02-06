package interpreter

import (
	"github.com/nicholasbailey/becca/exception"
	"github.com/nicholasbailey/becca/parser"
)

type BeccaValue struct {
	Type     *BeccaValue
	Value    interface{}
	Callable *Callable
	Methods  map[string]*Callable
}

type BuiltInFunction func(interpreter *Interpreter, values []*BeccaValue) (*BeccaValue, exception.Exception)

type Callable struct {
	Name                string
	Arity               int
	UserDefinedFunction *parser.Token
	BuiltInFunction     BuiltInFunction
}

func (left *BeccaValue) isEqualTo(right *BeccaValue) bool {
	if left.Type.Value != right.Type.Value {
		return false
	}
	switch left.Type.Value {
	case TString, TFloat, TBool, TInt:
		return left.Value == right.Value
	case TType:
		return areTypesEqual(left, right)
	case TFunction:
		return areFunctionsEqual(left, right)
	case TNull:
		return true
	default:
		return false
	}
}
