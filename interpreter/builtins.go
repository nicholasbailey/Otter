package interpreter

import (
	"fmt"

	"github.com/nicholasbailey/becca/common"
)

func Print(interpreter *Interpreter, values []*BeccaValue) (*BeccaValue, common.Exception) {
	for _, value := range values {
		switch value.Type.Value {
		case TString:
			fmt.Print(value.Value.(string))
		case TInt:
			// TODO - move away from builtin
			fmt.Print(value.Value.(int64))
		case TBool:
			fmt.Print(value.Value.(bool))
		case TFloat:
			fmt.Print(value.Value.(float64))
		case TNull:
			fmt.Print("<null>")
		}
		fmt.Print(" ")
	}
	fmt.Print("\n")
	return interpreter.NewNull(), nil
}
