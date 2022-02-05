package interpreter

import "github.com/nicholasbailey/becca/common"

func ConstructNull(interpreter *Interpreter, values []*BeccaValue) (*BeccaValue, common.Exception) {
	return interpreter.NewNull(), nil
}

func (interpreter *Interpreter) NewNull() *BeccaValue {
	return &BeccaValue{
		Type:  interpreter.MustResolveType(TNull),
		Value: nil,
	}
}
