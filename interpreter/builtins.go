package interpreter

import (
	"fmt"

	"github.com/nicholasbailey/otter/exception"
)

func Print(interpreter *Interpreter, values []*OtterValue) (*OtterValue, exception.Exception) {
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

func AssertEqual(interpreter *Interpreter, values []*OtterValue) (*OtterValue, exception.Exception) {
	left := values[0]
	right := values[1]
	if left.isEqualTo(right) {
		return interpreter.NewNull(), nil
	}
	leftAsString, err := ConstructString(interpreter, []*OtterValue{left})
	if err != nil {
		return nil, err
	}
	rightAsString, err := ConstructString(interpreter, []*OtterValue{right})
	if err != nil {
		return nil, err
	}
	errorMessage := fmt.Sprintf("%v is not equal to %v", leftAsString.Value, rightAsString.Value)
	return nil, exception.New(exception.AssertionError, errorMessage, 0, 0)
}

func AssertTrue(interpreter *Interpreter, values []*OtterValue) (*OtterValue, exception.Exception) {
	val := values[0]
	if val.Type.Value == TBool && val.Value == true {
		return interpreter.NewNull(), nil
	}

	return nil, exception.New(exception.AssertionError, "Failed asssertion", 0, 0)
}

func DefineBuiltins(interpreter *Interpreter) {
	printfn, _ := interpreter.NewBuiltInFunction("print", Variadic, Print)
	assertEqualFn, _ := interpreter.NewBuiltInFunction("assertEqual", 2, AssertEqual)
	assertTrueFn, _ := interpreter.NewBuiltInFunction("assertEqual", 1, AssertTrue)
	interpreter.DefineGlobal("print", printfn)
	interpreter.DefineGlobal("assertEqual", assertEqualFn)
	interpreter.DefineGlobal("assertTrue", assertTrueFn)
}
