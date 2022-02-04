package interpreter

import (
	"fmt"

	"github.com/nicholasbailey/becca/common"
	"github.com/nicholasbailey/becca/parser"
)

func (interpreter *Interpreter) defineFunction(tree *parser.Token) (*BeccaValue, error) {

	// TODO - validate inputs
	functionName := tree.Children[0].Value
	newValue := &BeccaValue{
		Type:               TFunction,
		Value:              nil, // TODO - figure out what this should be
		FunctionDefinition: tree,
	}
	err := interpreter.CallStack.AssignVariable(functionName, newValue)
	if err != nil {
		return nil, err
	}
	return newValue, nil
}

func (interpreter *Interpreter) callUserDefinedFunction(fn *BeccaValue, arguments []*BeccaValue) (*BeccaValue, error) {

	fnTree := fn.FunctionDefinition
	functionName := fn.FunctionDefinition.Children[0].Value
	// TODO - check for well formed tree here
	parameters := fnTree.Children[1].Children
	if len(parameters) != len(arguments) {
		return nil, common.NewException(common.TypeError, fmt.Sprintf("%v takes %v arguments, got %v", functionName, len(parameters), len(arguments)), fnTree.Line, fnTree.Col)
	}
	// TODO: Could this be cleaner
	stackFrame := NewCallStackFrame(parser.Symbol(functionName))
	for index, parameter := range parameters {
		arg := arguments[index]
		stackFrame.Scope[parameter.Value] = arg
	}
	interpreter.CallStack.Push(stackFrame)
	block := fn.FunctionDefinition.Children[2]
	returnValue, err := interpreter.Evaluate(block)
	if err != nil {
		return nil, err
	}
	interpreter.CallStack.Pop()
	return returnValue, nil
}

func (interpreter *Interpreter) callFunction(tree *parser.Token) (*BeccaValue, error) {

	functionName := tree.Children[0]
	function, builtInFound := interpreter.BuiltIns[functionName.Value]
	var udf *BeccaValue
	var err error
	if !builtInFound {
		udf, err = interpreter.resolveName(functionName)
		if err != nil {
			return nil, common.NewException(common.NameError, fmt.Sprintf("%v is not defined", functionName.Value), tree.Line, tree.Col)
		}
		if udf.FunctionDefinition == nil {
			return nil, common.NewException(common.TypeError, fmt.Sprintf("%v object is not callable", udf.Type), tree.Line, tree.Col)
		}
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
	if builtInFound {
		return function(arguments)
	}
	return interpreter.callUserDefinedFunction(udf, arguments)
}
