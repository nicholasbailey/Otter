package interpreter

import "github.com/nicholasbailey/becca/exception"

func ConstructNull(interpreter *Interpreter, values []*BeccaValue) (*BeccaValue, exception.Exception) {
	return interpreter.NewNull(), nil
}

func (interpreter *Interpreter) NewNull() *BeccaValue {
	return &BeccaValue{
		Type:  interpreter.MustResolveType(TNull),
		Value: nil,
	}
}
