package interpreter

import (
	"fmt"

	"github.com/nicholasbailey/otter/exception"
)

func ConstructArray(interpreter *Interpreter, values []*OtterValue) (*OtterValue, exception.Exception) {
	newSlice := make([]*OtterValue, len(values))
	copy(newSlice, values)
	return &OtterValue{
		Type:  interpreter.MustResolveType(TArray),
		Value: newSlice,
	}, nil
}

func ArrayLength(interpreter *Interpreter, values []*OtterValue) (*OtterValue, exception.Exception) {
	array := values[0]
	underlyingSlice := array.Value.([]*OtterValue)
	length := len(underlyingSlice)
	return interpreter.NewInt(int64(length)), nil
}

func ArrayAdd(interpreter *Interpreter, values []*OtterValue) (*OtterValue, exception.Exception) {
	array := values[0]
	underlyingSlice := array.Value.([]*OtterValue)
	underlyingSlice = append(underlyingSlice, values[1:]...)
	array.Value = underlyingSlice
	return interpreter.NewNull(), nil
}

func ArrayGetItem(interpreter *Interpreter, values []*OtterValue) (*OtterValue, exception.Exception) {
	array := values[0]
	index := values[1]
	if !index.IsInstanceOf(TInt) {
		return nil, exception.New(exception.ArgumentError, fmt.Sprintf("Array Index must be int, got %v", index.Type.Value), 0, 0)
	}
	underlyingSlice := array.Value.([]*OtterValue)
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
