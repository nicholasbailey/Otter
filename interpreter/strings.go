package interpreter

import (
	"strconv"

	"github.com/nicholasbailey/becca/common"
)

func ConstructString(interpreter *Interpreter, values []*BeccaValue) (*BeccaValue, common.Exception) {
	if len(values) > 1 || len(values) == 0 {
		// TODO - get call stack info for builtins
		return nil, common.NewException("ArgumentError", "", 0, 0)
	}
	value := values[0]
	var strVal string
	switch value.Type.Value {
	case TString:
		strVal = value.Value.(string)
	case TInt:
		// TODO - move away from builtin
		strVal = strconv.FormatInt(value.Value.(int64), 10)
	case TBool:
		if value.Value == true {
			strVal = "true"
		} else {
			strVal = "false"
		}
	case TFloat:
		strVal = strconv.FormatFloat(value.Value.(float64), 'f', -1, 64)
	case TNull:
		strVal = "<null>"
	case TFunction:
		strVal = value.Callable.Name
	default:
		strVal = "[Object]"
	}
	return &BeccaValue{
		Type:  interpreter.MustResolveType(TString),
		Value: strVal,
	}, nil
}

func (interpreter *Interpreter) NewString(s string) *BeccaValue {
	return &BeccaValue{
		Type:  interpreter.MustResolveType(TString),
		Value: s,
	}
}

func StringLength(interpreter *Interpreter, values []*BeccaValue) (*BeccaValue, common.Exception) {
	// TODO - error handling
	v := values[0]
	s := v.Value.(string)
	length := int64(len(s))
	return interpreter.NewInt(length), nil
}

func DefineStringType(interpreter *Interpreter) {
	interpreter.DefineType(TString, NewBuiltInConstructor(TString, 1, ConstructString))
	length, _ := interpreter.NewBuiltInFunction("length", 1, StringLength)

	interpreter.DefineMethod(TString, "length", length.Callable)
}
