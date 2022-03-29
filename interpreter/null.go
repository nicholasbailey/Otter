package interpreter

import "github.com/nicholasbailey/otter/exception"

func ConstructNull(interpreter *Interpreter, values []*OtterValue) (*OtterValue, exception.Exception) {
	return interpreter.NewNull(), nil
}

func (interpreter *Interpreter) NewNull() *OtterValue {
	return &OtterValue{
		Type:  interpreter.MustResolveType(TNull),
		Value: nil,
	}
}
