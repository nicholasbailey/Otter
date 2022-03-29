package interpreter

import (
	"github.com/nicholasbailey/otter/exception"
	"github.com/nicholasbailey/otter/parser"
)

func (interpreter *Interpreter) doIf(tree *parser.Token) (*OtterValue, error) {
	if len(tree.Children) < 2 {
		exception.New(exception.SyntaxError, "invalid if expression", tree.Line, tree.Col)
	}
	condition := tree.Children[0]
	conditionValue, err := interpreter.Evaluate(condition)
	if err != nil {
		return nil, err
	}
	executeCondition := interpreter.Truthiness(conditionValue)
	if executeCondition.Value == true {
		block := tree.Children[1]
		return interpreter.Evaluate(block)
	}
	for _, child := range tree.Children[2:] {
		if child.Symbol == parser.ElseIf {
			return interpreter.doIf(child)
		} else {
			// TODO - more checks here
			return interpreter.Evaluate(child)
		}
	}
	return interpreter.NewNull(), nil
}
