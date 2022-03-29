package interpreter

import "github.com/nicholasbailey/otter/exception"

type TypeName string

const (
	TString         TypeName = "string"
	TInt            TypeName = "int"
	TBool           TypeName = "bool"
	TFloat          TypeName = "float"
	TNull           TypeName = "null"
	TFunction       TypeName = "function"
	TType           TypeName = "type"
	TArray          TypeName = "Array"
	TStringIterator TypeName = "StringIterator"
)

func ConstructType(interpreter *Interpreter, values []*OtterValue) (*OtterValue, exception.Exception) {
	v := values[0]
	return v.Type, nil
}

func (interpreter *Interpreter) ResolveType(typeName TypeName) (*OtterValue, exception.Exception) {
	val, _ := interpreter.CallStack.ResolveVariable(string(typeName))
	// TODO - error handling and make this more efficient
	return val, nil
}

func (interpreter *Interpreter) MustResolveType(typeName TypeName) *OtterValue {
	val, err := interpreter.ResolveType(typeName)
	if err != nil {
		panic(err)
	}
	return val
}

func (interpreter *Interpreter) DefineType(t TypeName, constructor *Callable) (*OtterValue, exception.Exception) {
	value := &OtterValue{
		Type:     interpreter.MustResolveType(TType),
		Value:    t,
		Callable: constructor,
		Methods:  map[string]*Callable{},
	}
	err := interpreter.CallStack.AssignVariable(string(t), value)
	if err != nil {
		return nil, err
	}
	return value, err
}

func (value *OtterValue) IsInstanceOf(typeName TypeName) bool {
	return value.Type.Value == typeName
}

// Tests if two objects of type 'type' are equal
func areTypesEqual(left *OtterValue, right *OtterValue) bool {
	// TODO - this is not safe long term
	return left.Value == right.Value
}

func DefineTypeType(interpreter *Interpreter) {
	typeConstructor := &Callable{
		Name:                "type",
		Arity:               1,
		UserDefinedFunction: nil,
		BuiltInFunction:     ConstructType,
	}

	// Boostrap the type type
	typeVal := OtterValue{
		Value:    TType,
		Callable: typeConstructor,
	}

	typeVal.Type = &typeVal

	interpreter.CallStack.Globals().Scope["type"] = &typeVal
}
