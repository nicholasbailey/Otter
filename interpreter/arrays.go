package interpreter

import (
	"fmt"

	"github.com/nicholasbailey/becca/exception"
)

func ConstructArray(interpreter *Interpreter, values []*BeccaValue) (*BeccaValue, exception.Exception) {
	newSlice := make([]*BeccaValue, len(values))
	copy(newSlice, values)
	return &BeccaValue{
		Type:  interpreter.MustResolveType(TArray),
		Value: newSlice,
	}, nil
}

func ArrayLength(interpreter *Interpreter, values []*BeccaValue) (*BeccaValue, exception.Exception) {
	array := values[0]
	underlyingSlice := array.Value.([]*BeccaValue)
	length := len(underlyingSlice)
	return interpreter.NewInt(int64(length)), nil
}

func ArrayAdd(interpreter *Interpreter, values []*BeccaValue) (*BeccaValue, exception.Exception) {
	array := values[0]
	underlyingSlice := array.Value.([]*BeccaValue)
	underlyingSlice = append(underlyingSlice, values[1:]...)
	array.Value = underlyingSlice
	return interpreter.NewNull(), nil
}

func ArrayGetItem(interpreter *Interpreter, values []*BeccaValue) (*BeccaValue, exception.Exception) {
	array := values[0]
	index := values[1]
	if !index.IsInstanceOf(TInt) {
		return nil, exception.New(exception.ArgumentError, fmt.Sprintf("Array Index must be int, got %v", index.Type.Value), 0, 0)
	}
	underlyingSlice := array.Value.([]*BeccaValue)
	indexValue := index.Value.(int64)

	if indexValue < 0 || len(underlyingSlice) < int(indexValue) {
		return nil, exception.New(exception.IndexError, "Array index out of range", 0, 0)
	}
	return underlyingSlice[indexValue], nil
}

func DefineArrayType(interpreter *Interpreter) {
	interpreter.DefineType(TArray, NewBuiltInConstructor(TArray, Variadic, ConstructArray))
	interpreter.DefineBuiltinMethod(TArray, "length", 1, ArrayLength)
	interpreter.DefineBuiltinMethod(TArray, "append", Variadic, ArrayAdd)
	interpreter.DefineBuiltinMethod(TArray, "getItem", 2, ArrayGetItem)
}
