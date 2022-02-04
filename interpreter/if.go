package interpreter

import (
	"github.com/nicholasbailey/becca/common"
	"github.com/nicholasbailey/becca/parser"
)

func (interpreter *Interpreter) doIf(tree *parser.Token) (*BeccaValue, error) {
	if len(tree.Children) < 2 {
		common.NewException(common.SyntaxError, "invalid if expression", tree.Line, tree.Col)
	}
	condition := tree.Children[0]
	conditionValue, err := interpreter.Evaluate(condition)
	if err != nil {
		return nil, err
	}
	executeCondition := Truthiness(conditionValue)
	if executeCondition.Value == true {
		block := tree.Children[1]
		return interpreter.Evaluate(block)
	}
	for _, child := range tree.Children[2:] {
		if child.Symbol == parser.ElseIf {
			return interpreter.doIf(child)
		} else {
			return interpreter.Evaluate(child)
		}
	}
	return Null(), nil
}
