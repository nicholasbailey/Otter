package interpreter

import (
	"fmt"

	"github.com/nicholasbailey/becca/exception"
)

func ConstructFloat(interpreter *Interpreter, values []*BeccaValue) (*BeccaValue, exception.Exception) {
	v := values[0]
	if v.IsInstanceOf(TFloat) {
		return v, nil
	} else {
		// TODO - make lines and cols work
		return nil, exception.New(exception.TypeError, fmt.Sprintf("cannot convert %v to float", v.Type.Value), 0, 0)
	}
}

func (interpreter *Interpreter) NewFloat(f float64) *BeccaValue {
	return &BeccaValue{
		Type:  interpreter.MustResolveType(TFloat),
		Value: f,
	}
}
