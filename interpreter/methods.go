package interpreter

import (
	"fmt"

	"github.com/nicholasbailey/otter/exception"
	"github.com/nicholasbailey/otter/parser"
)

func (interpreter *Interpreter) doAccess(tree *parser.Token) (*OtterValue, exception.Exception) {
	valueTree := tree.Children[0]
	targetTree := tree.Children[1]
	value, err := interpreter.Evaluate(valueTree)
	if err != nil {
		return nil, err
	}
	var methodName string
	arguments := []*OtterValue{}
	if targetTree.Symbol == parser.Name {
		methodName = targetTree.Value
	} else if targetTree.Symbol == parser.FunctionInvocation {
		methodName = targetTree.Children[0].Value
		for _, childToken := range targetTree.Children[1:] {
			childValue, err := interpreter.Evaluate(childToken)
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, childValue)
		}
	}
	return interpreter.callMethod(value, methodName, arguments)
}

func (interpreter *Interpreter) callMethod(value *OtterValue, methodName string, arguments []*OtterValue) (*OtterValue, exception.Exception) {

	method, found := value.Type.Methods[methodName]
	if !found {
		// TODO - handle line and col
		return nil, exception.New(exception.MethodError, fmt.Sprintf("%v has no method %v", value.Type.Value, methodName), 0, 0)
	}
	fullArguments := []*OtterValue{value}

	fullArguments = append(fullArguments, arguments...)
	return interpreter.invokeCallable(method, fullArguments, 0, 0)
}
