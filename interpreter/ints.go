package interpreter

import (
	"fmt"
	"strconv"

	"github.com/nicholasbailey/becca/exception"
)

func ConstructInt(interpreter *Interpreter, values []*BeccaValue) (*BeccaValue, exception.Exception) {
	v := values[0]
	if v.IsInstanceOf(TInt) {
		return v, nil
	} else if v.IsInstanceOf(TString) {
		parsedInt, err := strconv.ParseInt(v.Value.(string), 0, 64)
		if err != nil {
			return nil, err
		}
		return interpreter.NewInt(parsedInt), nil
	} else {
		// TODO - make lines and cols work
		return nil, exception.New(exception.TypeError, fmt.Sprintf("cannot convert %v to int", v.Type.Value), 0, 0)
	}
}

func (interpreter *Interpreter) NewInt(i int64) *BeccaValue {
	return &BeccaValue{
		Type:     interpreter.MustResolveType(TInt),
		Value:    i,
		Callable: nil,
	}
}
