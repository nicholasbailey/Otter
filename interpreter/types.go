package interpreter

import "github.com/nicholasbailey/becca/exception"

type TypeName string

const (
	TString   TypeName = "string"
	TInt      TypeName = "int"
	TBool     TypeName = "bool"
	TFloat    TypeName = "float"
	TNull     TypeName = "null"
	TFunction TypeName = "function"
	TType     TypeName = "type"
)

func ConstructType(interpreter *Interpreter, values []*BeccaValue) (*BeccaValue, exception.Exception) {
	v := values[0]
	return v.Type, nil
}

func (interpreter *Interpreter) ResolveType(typeName TypeName) (*BeccaValue, exception.Exception) {
	val, _ := interpreter.CallStack.ResolveVariable(string(typeName))
	// TODO - error handling and make this more efficient
	return val, nil
}

func (interpreter *Interpreter) MustResolveType(typeName TypeName) *BeccaValue {
	val, err := interpreter.ResolveType(typeName)
	if err != nil {
		panic(err)
	}
	return val
}

func (interpreter *Interpreter) DefineType(t TypeName, constructor *Callable) (*BeccaValue, exception.Exception) {
	value := &BeccaValue{
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

func (value *BeccaValue) IsInstanceOf(typeName TypeName) bool {
	return value.Type.Value == typeName
}

// Tests if two objects of type 'type' are equal
func areTypesEqual(left *BeccaValue, right *BeccaValue) bool {
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
	typeVal := BeccaValue{
		Value:    TType,
		Callable: typeConstructor,
	}

	typeVal.Type = &typeVal

	interpreter.CallStack.Globals().Scope["type"] = &typeVal
}
