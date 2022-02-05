package interpreter

import (
	"fmt"

	"github.com/nicholasbailey/becca/common"
	"github.com/nicholasbailey/becca/parser"
)

const Variadic = -1

func ConstructFunction(interpreter *Interpreter, values []*BeccaValue) (*BeccaValue, common.Exception) {
	return nil, common.NewException(common.NameError, "function is not callable", 0, 0)
}

func ValidateFunctionDefinition(tree *parser.Token) common.Exception {
	if tree == nil {
		return common.NewException(common.InternalError, "null token passed to NewUserDefinedFunction", 0, 0)
	}
	if tree.Symbol != parser.FunctionDefinition {
		return common.NewException(common.InternalError, fmt.Sprintf("non function definition token %v passed to NewUserDefinedFunction", tree.Symbol), tree.Line, tree.Col)
	}
	return nil
}

func (interpreter *Interpreter) NewBuiltInFunction(name string, arity int, builtIn BuiltInFunction) (*BeccaValue, common.Exception) {
	// TODO - santize inputs
	callable := &Callable{
		Name:                name,
		Arity:               arity,
		BuiltInFunction:     builtIn,
		UserDefinedFunction: nil,
	}
	return &BeccaValue{
		Type:     interpreter.MustResolveType(TFunction),
		Value:    nil,
		Callable: callable,
	}, nil
}

// Gott a come up with a better name here
func NewBuiltInConstructor(typeName TypeName, arity int, builtIn BuiltInFunction) *Callable {
	return &Callable{
		Name:                string(typeName),
		Arity:               arity,
		BuiltInFunction:     builtIn,
		UserDefinedFunction: nil,
	}
}

func (interpreter *Interpreter) NewUserDefinedFunction(tree *parser.Token) (*BeccaValue, common.Exception) {
	err := ValidateFunctionDefinition(tree)
	if err != nil {
		return nil, err
	}
	functionName := tree.Children[0].Value
	parameters := tree.Children[1].Children

	callable := &Callable{
		UserDefinedFunction: tree,
		Arity:               len(parameters),
		BuiltInFunction:     nil,
		Name:                functionName,
	}

	return &BeccaValue{
		Type:     interpreter.MustResolveType(TFunction),
		Value:    nil, // TODO - figure out what this should be
		Callable: callable,
	}, nil
}

func (interpreter *Interpreter) defineFunction(tree *parser.Token) (*BeccaValue, error) {

	udf, err := interpreter.NewUserDefinedFunction(tree)
	if err != nil {
		return nil, err
	}
	// TODO - prevent overriding builtins
	err = interpreter.CallStack.AssignVariable(udf.Callable.Name, udf)
	if err != nil {
		return nil, err
	}
	return udf, nil
}

func (interpreter *Interpreter) invokeCallable(callable *Callable, arguments []*BeccaValue, line int, col int) (*BeccaValue, common.Exception) {
	arity := callable.Arity
	if arity != Variadic && len(arguments) != arity {
		return nil, common.NewException(common.TypeError, fmt.Sprintf("%v takes exactly %v arguments, found %v", callable.Name, callable.Arity, len(arguments)), line, col)
	}
	if callable.BuiltInFunction != nil {
		return callable.BuiltInFunction(interpreter, arguments)
	}
	udf := callable.UserDefinedFunction
	parameters := udf.Children[1].Children
	if len(parameters) != len(arguments) {
		return nil, common.NewException(common.TypeError, fmt.Sprintf("%v takes %v arguments, got %v", callable.Name, len(parameters), len(arguments)), line, col)
	}
	// TODO: Could this be cleaner
	stackFrame := NewCallStackFrame(callable.Name)
	for index, parameter := range parameters {
		arg := arguments[index]
		stackFrame.Scope[parameter.Value] = arg
	}
	interpreter.CallStack.Push(stackFrame)
	block := udf.Children[2]
	var err error
	for _, child := range block.Children {
		_, err = interpreter.Evaluate(child)
		if err != nil {
			break
		}
		stackFrame := interpreter.CallStack.Peek()
		// TODO handle nil stack frame

		if stackFrame.ReturnValue != nil {
			break
		}
	}
	frame := interpreter.CallStack.Pop()
	if err != nil {
		return nil, err
	}
	if frame.ReturnValue == nil {
		frame.ReturnValue = interpreter.NewNull()
	}
	return frame.ReturnValue, nil
}

// Should probably not be called call function, as it is also the syntax for other calls
func (interpreter *Interpreter) callFunction(tree *parser.Token) (*BeccaValue, common.Exception) {
	// TODO - check inputs
	functionName := tree.Children[0]
	functionValue, err := interpreter.resolveName(functionName)
	if err != nil {
		return nil, err
	}

	if functionValue.Callable == nil {
		return nil, common.NewException(common.TypeError, fmt.Sprintf("%v is not callable", functionName.Value), tree.Line, tree.Col)
	}

	// TODO - optimize memory allocation here
	arguments := []*BeccaValue{}
	for _, childToken := range tree.Children[1:] {
		childValue, err := interpreter.Evaluate(childToken)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, childValue)
	}

	return interpreter.invokeCallable(functionValue.Callable, arguments, tree.Line, tree.Col)
}
