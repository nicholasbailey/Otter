package interpreter

import (
	"strconv"

	"github.com/nicholasbailey/becca/common"
)

func String(values []*BeccaValue) (*BeccaValue, common.Exception) {
	if len(values) > 1 || len(values) == 0 {
		// TODO - get call stack info for builtins
		return nil, common.NewException("ArgumentError", "", 0, 0)
	}
	value := values[0]
	var strVal string
	switch value.Type {
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
		// TODO - write a util for this
		strVal = value.FunctionDefinition.Children[0].Value
	default:
		strVal = "[Object]"
	}
	return &BeccaValue{
		Type:  TString,
		Value: strVal,
	}, nil
}

func GoStringToBeccaString(s string) *BeccaValue {
	return &BeccaValue{
		Type:  TString,
		Value: s,
	}
}
