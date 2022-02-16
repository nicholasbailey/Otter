package interpreter

import (
	"strconv"
	"strings"

	"github.com/nicholasbailey/becca/exception"
)

func ConstructString(interpreter *Interpreter, values []*BeccaValue) (*BeccaValue, exception.Exception) {
	if len(values) != 1 {
		// TODO - get call stack info for builtins
		return nil, exception.New(exception.ArgumentError, "", 0, 0)
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

func StringLength(interpreter *Interpreter, values []*BeccaValue) (*BeccaValue, exception.Exception) {
	// TODO - error handling
	v := values[0]
	s := v.Value.(string)
	length := int64(len(s))
	return interpreter.NewInt(length), nil
}

func StringToUpperCase(interpreter *Interpreter, values []*BeccaValue) (*BeccaValue, exception.Exception) {
	v := values[0]
	s := v.Value.(string)
	newVal := strings.ToUpper(s)
	return interpreter.NewString(newVal), nil
}

func StringToLowerCase(interpreter *Interpreter, values []*BeccaValue) (*BeccaValue, exception.Exception) {
	v := values[0]
	s := v.Value.(string)
	newVal := strings.ToLower(s)
	return interpreter.NewString(newVal), nil
}

func StringReplace(interpreter *Interpreter, values []*BeccaValue) (*BeccaValue, exception.Exception) {
	base := values[0]
	target := values[1]
	replacement := values[2]

	baseStr := base.Value.(string)
	targetStr := target.Value.(string)
	replacementStr := replacement.Value.(string)

	newStr := strings.Replace(baseStr, targetStr, replacementStr, -1)
	return interpreter.NewString(newStr), nil
}

func StringIterator(interpreter *Interpreter, values []*BeccaValue) (*BeccaValue, exception.Exception) {
	return ConstructStringIterator(interpreter, values)
}

type StringIteratorInternals struct {
	String string
	Index  int
}

func ConstructStringIterator(interpreter *Interpreter, values []*BeccaValue) (*BeccaValue, exception.Exception) {
	if len(values) > 1 || len(values) == 0 {
		// TODO - get call stack info for builtins
		return nil, exception.New(exception.ArgumentError, "", 0, 0)
	}
	value := values[0]
	if !value.IsInstanceOf(TString) {
		return nil, exception.New(exception.ArgumentError, "argument str to constructor StringIterator must be a string", 0, 0)
	}
	iteratorValue := StringIteratorInternals{
		String: value.Value.(string),
		Index:  0,
	}
	return &BeccaValue{
		Type:  interpreter.MustResolveType(TStringIterator),
		Value: &iteratorValue,
	}, nil
}

func StringIteratorHasNext(interpreter *Interpreter, values []*BeccaValue) (*BeccaValue, exception.Exception) {
	iterator := values[0]
	internals := iterator.Value.(*StringIteratorInternals)
	if internals.Index >= len(internals.String) {
		return interpreter.False(), nil
	} else {
		return interpreter.True(), nil
	}
}

func StringIteratorGetNext(interpreter *Interpreter, values []*BeccaValue) (*BeccaValue, exception.Exception) {
	iterator := values[0]
	internals := iterator.Value.(*StringIteratorInternals)
	if internals.Index >= len(internals.String) {
		return nil, exception.New(exception.IterationError, "iterable has no more elements", 0, 0)
	}
	value := string([]rune(internals.String)[internals.Index])
	internals.Index = internals.Index + 1
	return interpreter.NewString(value), nil
}

func DefineStringTypes(interpreter *Interpreter) {

	interpreter.DefineType(TString, NewBuiltInConstructor(TString, 1, ConstructString))
	interpreter.DefineBuiltinMethod(TString, "length", 1, StringLength)
	interpreter.DefineBuiltinMethod(TString, "toUpperCase", 1, StringToUpperCase)
	interpreter.DefineBuiltinMethod(TString, "toLowerCase", 1, StringToLowerCase)
	interpreter.DefineBuiltinMethod(TString, "replace", 3, StringReplace)
	interpreter.DefineBuiltinMethod(TString, "iterator", 1, StringIterator)

	interpreter.DefineType(TStringIterator, NewBuiltInConstructor(TStringIterator, 1, ConstructStringIterator))
	interpreter.DefineBuiltinMethod(TStringIterator, "hasNext", 1, StringIteratorHasNext)
	interpreter.DefineBuiltinMethod(TStringIterator, "getNext", 1, StringIteratorGetNext)

}
