package interpreter

import (
	"fmt"

	"github.com/nicholasbailey/otter/exception"
	"github.com/nicholasbailey/otter/parser"
)

type OtterValue struct {
	Type     *OtterValue
	Value    interface{}
	Callable *Callable
	Methods  map[string]*Callable
}

func (v *OtterValue) String() string {
	return fmt.Sprintf("%v", v.Value)
}

type BuiltInFunction func(interpreter *Interpreter, values []*OtterValue) (*OtterValue, exception.Exception)

type Callable struct {
	Name                string
	Arity               int
	UserDefinedFunction *parser.Token
	BuiltInFunction     BuiltInFunction
}

func (left *OtterValue) isEqualTo(right *OtterValue) bool {
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
