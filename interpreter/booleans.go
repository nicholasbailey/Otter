package interpreter

import "github.com/nicholasbailey/becca/exception"

func (interpreter *Interpreter) NewBool(x bool) *BeccaValue {
	return &BeccaValue{
		Type:  interpreter.MustResolveType(TBool),
		Value: x,
	}
}

func (interpreter *Interpreter) False() *BeccaValue {
	return interpreter.NewBool(false)
}

func (interpreter *Interpreter) True() *BeccaValue {
	return interpreter.NewBool(true)
}

func ConstructBool(interpreter *Interpreter, values []*BeccaValue) (*BeccaValue, exception.Exception) {
	v := values[0]
	return interpreter.Truthiness(v), nil
}

func (interpreter *Interpreter) Truthiness(value *BeccaValue) *BeccaValue {
	switch value.Type.Value {
	case TBool:
		return value
	case TString:
		if value.Value.(string) == "" {
			return interpreter.False()
		} else {
			return interpreter.True()
		}
	case TInt:
		if value.Value.(int64) == 0 {
			return interpreter.False()
		} else {
			return interpreter.True()
		}
	case TFloat:
		if value.Value.(float64) == 0.0 {
			return interpreter.False()
		} else {
			return interpreter.True()
		}
	case TNull:
		return interpreter.False()
	case TFunction:
		return interpreter.True()
	}
	// TODO - handle error correctly
	panic("How did we get here")
}
