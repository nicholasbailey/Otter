package interpreter

import (
	"github.com/nicholasbailey/becca/exception"
	"github.com/nicholasbailey/becca/parser"
)

func (interpreter *Interpreter) doFor(tree *parser.Token) (*BeccaValue, error) {
	return nil, nil
}

func (interpreter *Interpreter) doWhile(tree *parser.Token) (*BeccaValue, error) {
	if len(tree.Children) != 2 {
		return nil, exception.New(exception.SyntaxError, "invalid while block", tree.Line, tree.Col)
	}
	expression := tree.Children[0]
	block := tree.Children[1]
	retVal := interpreter.NewNull()
	for {
		expressionRes, err := interpreter.Evaluate(expression)
		if err != nil {
			return nil, err
		}
		expressionTruthiness := interpreter.Truthiness(expressionRes)
		if expressionTruthiness.Value == false {
			break
		}
		retVal, err = interpreter.Evaluate(block)
		if err != nil {
			return nil, err
		}
	}
	return retVal, nil
}
