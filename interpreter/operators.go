package interpreter

import (
	"fmt"

	"github.com/nicholasbailey/becca/common"
	"github.com/nicholasbailey/becca/parser"
)

func resolveBinaryOperands(interpreter *Interpreter, tree *parser.Token) (*BeccaValue, *BeccaValue, error) {
	if len(tree.Children) != 2 {
		return nil, nil, common.NewException(common.SyntaxError, fmt.Sprintf("invalid symbol %v", tree.Value), tree.Line, tree.Col)
	}
	left := tree.Children[0]
	right := tree.Children[1]
	leftValue, leftErr := interpreter.Evaluate(left)

	rightValue, rightErr := interpreter.Evaluate(right)
	if leftErr != nil {
		return leftValue, rightValue, leftErr
	} else {
		return leftValue, rightValue, rightErr
	}
}

func (interpreter *Interpreter) doLessThan(tree *parser.Token) (*BeccaValue, error) {
	leftValue, rightValue, err := resolveBinaryOperands(interpreter, tree)
	if err != nil {
		return nil, err
	}
	if leftValue.IsInstanceOf(TInt) && rightValue.IsInstanceOf(TInt) {
		return interpreter.NewBool(leftValue.Value.(int64) < rightValue.Value.(int64)), nil
	} else if leftValue.IsInstanceOf(TFloat) && rightValue.IsInstanceOf(TFloat) {
		return interpreter.NewBool(leftValue.Value.(float64) < rightValue.Value.(float64)), nil
	} else if leftValue.IsInstanceOf(TString) && rightValue.IsInstanceOf(TString) {
		return interpreter.NewBool(leftValue.Value.(string) < rightValue.Value.(string)), nil
	} else if leftValue.Type == rightValue.Type {
		return nil, common.NewException(common.TypeError, fmt.Sprintf("type %v cannot be compared with <", rightValue.Type), tree.Line, tree.Col)
	}
	return nil, common.NewException(common.TypeError, "attempted to compare incomparable types with <", tree.Line, tree.Col)
}

func (interpreter *Interpreter) doGreaterThan(tree *parser.Token) (*BeccaValue, error) {
	leftValue, rightValue, err := resolveBinaryOperands(interpreter, tree)
	if err != nil {
		return nil, err
	}
	if leftValue.IsInstanceOf(TInt) && rightValue.IsInstanceOf(TInt) {
		return interpreter.NewBool(leftValue.Value.(int64) > rightValue.Value.(int64)), nil
	} else if leftValue.IsInstanceOf(TFloat) && rightValue.IsInstanceOf(TFloat) {
		return interpreter.NewBool(leftValue.Value.(float64) > rightValue.Value.(float64)), nil
	} else if leftValue.IsInstanceOf(TString) && rightValue.IsInstanceOf(TString) {
		return interpreter.NewBool(leftValue.Value.(string) > rightValue.Value.(string)), nil
	} else if leftValue.Type == rightValue.Type {
		return nil, common.NewException(common.TypeError, fmt.Sprintf("type %v cannot be compared with >", rightValue.Type), tree.Line, tree.Col)
	}
	return nil, common.NewException(common.TypeError, "attempted to compare incomparable types with >", tree.Line, tree.Col)
}

func (interpreter *Interpreter) doLessThanOrEqualTo(tree *parser.Token) (*BeccaValue, error) {
	equal, err := interpreter.doEqualityCheck(tree)
	if err != nil {
		return nil, err
	}
	if equal.Value == true {
		return interpreter.True(), nil
	}
	lessThan, err := interpreter.doLessThan(tree)
	if err != nil {
		return nil, err
	}
	if lessThan.Value == true {
		return interpreter.True(), nil
	}
	return interpreter.False(), nil
}

func (interpreter *Interpreter) doGreaterThanOrEqualTo(tree *parser.Token) (*BeccaValue, error) {
	equal, err := interpreter.doEqualityCheck(tree)
	if err != nil {
		return nil, err
	}
	if equal.Value == true {
		return interpreter.True(), nil
	}
	lessThan, err := interpreter.doGreaterThan(tree)
	if err != nil {
		return nil, err
	}
	if lessThan.Value == true {
		return interpreter.True(), nil
	}
	return interpreter.False(), nil
}

func (interpreter *Interpreter) doEqualityCheck(tree *parser.Token) (*BeccaValue, error) {
	leftValue, rightValue, err := resolveBinaryOperands(interpreter, tree)
	if err != nil {
		return nil, err
	}
	if leftValue.Type != rightValue.Type {
		return interpreter.False(), nil
	}
	if leftValue.Value == rightValue.Value {
		return interpreter.True(), nil
	}
	return interpreter.False(), nil
}

func (interpreter *Interpreter) doInequalityCheck(tree *parser.Token) (*BeccaValue, error) {
	result, err := interpreter.doEqualityCheck(tree)
	if err != nil {
		return nil, err
	}
	if result.Value == false {
		return interpreter.True(), nil
	} else {
		return interpreter.False(), nil
	}
}

func (interpreter *Interpreter) doAnd(tree *parser.Token) (*BeccaValue, error) {

	leftValue, rightValue, err := resolveBinaryOperands(interpreter, tree)

	if err != nil {
		return nil, err
	}
	leftTruthy := interpreter.Truthiness(leftValue)

	if leftTruthy.Value == false {
		return leftValue, nil
	}

	return rightValue, nil
}

func (interpreter *Interpreter) doOr(tree *parser.Token) (*BeccaValue, error) {
	leftValue, rightValue, err := resolveBinaryOperands(interpreter, tree)
	if err != nil {
		return nil, err
	}
	leftTruthy := interpreter.Truthiness(leftValue)
	if leftTruthy.Value == true {
		return leftValue, nil
	}
	return rightValue, nil
}

func (interpreter *Interpreter) doAssigment(tree *parser.Token) (*BeccaValue, error) {
	if len(tree.Children) != 2 {
		// TODO make more detailed
		return nil, fmt.Errorf("syntaxerror: invald assignment expression at line %v, col %v", tree.Line, tree.Col)
	}
	left := tree.Children[0]
	right := tree.Children[1]
	if left.Symbol != parser.Name {
		return nil, common.NewException(common.SyntaxError, "invalid assigment expression at line %v, col %v", tree.Line, tree.Col)
	}
	rightValue, err := interpreter.Evaluate(right)
	if err != nil {
		return nil, err
	}

	// TODO - handle colisions with builtins
	err = interpreter.CallStack.AssignVariable(left.Value, rightValue)
	if err != nil {
		return nil, err
	}
	return rightValue, nil
}

func (interpreter *Interpreter) doAddition(tree *parser.Token) (*BeccaValue, error) {
	leftValue, rightValue, err := resolveBinaryOperands(interpreter, tree)
	if err != nil {
		return nil, err
	}
	if leftValue.IsInstanceOf(TInt) && rightValue.IsInstanceOf(TInt) {
		newValue := leftValue.Value.(int64) + rightValue.Value.(int64)
		return interpreter.NewInt(newValue), nil
	}
	if leftValue.IsInstanceOf(TFloat) && rightValue.IsInstanceOf(TFloat) {
		newValue := leftValue.Value.(float64) + rightValue.Value.(float64)
		return interpreter.NewFloat(newValue), nil
	}
	if leftValue.IsInstanceOf(TString) && rightValue.IsInstanceOf(TString) {
		newValue := leftValue.Value.(string) + rightValue.Value.(string)
		return interpreter.NewString(newValue), nil
	}
	if leftValue.Type == rightValue.Type {
		return nil, fmt.Errorf("typeerror: type %v does not support operator + at line %v, col %v", leftValue.Type, tree.Line, tree.Col)
	}
	return nil, fmt.Errorf("typerror: incompatable types %v and %v with operator + at line %v, col %v", leftValue.Type, rightValue.Type, tree.Line, tree.Col)
}

func (interpreter *Interpreter) doSubtraction(tree *parser.Token) (*BeccaValue, error) {
	leftValue, rightValue, err := resolveBinaryOperands(interpreter, tree)
	if err != nil {
		return nil, err
	}
	if leftValue.IsInstanceOf(TInt) && rightValue.IsInstanceOf(TInt) {
		newValue := leftValue.Value.(int64) - rightValue.Value.(int64)
		return interpreter.NewInt(newValue), nil
	}
	if leftValue.IsInstanceOf(TFloat) && rightValue.IsInstanceOf(TFloat) {
		newValue := leftValue.Value.(float64) - rightValue.Value.(float64)
		return interpreter.NewFloat(newValue), nil
	}
	if leftValue.Type == rightValue.Type {
		return nil, fmt.Errorf("typeerror: type %v does not support operator - at line %v, col %v", leftValue.Type, tree.Line, tree.Col)
	}
	return nil, fmt.Errorf("typerror: incompatable types %v and %v with operator - at line %v, col %v", leftValue.Type, rightValue.Type, tree.Line, tree.Col)
}

func (interpreter *Interpreter) doMultiplication(tree *parser.Token) (*BeccaValue, error) {
	leftValue, rightValue, err := resolveBinaryOperands(interpreter, tree)
	if err != nil {
		return nil, err
	}
	if leftValue.IsInstanceOf(TInt) && rightValue.IsInstanceOf(TInt) {
		newValue := leftValue.Value.(int64) * rightValue.Value.(int64)
		return interpreter.NewInt(newValue), nil
	}
	if leftValue.IsInstanceOf(TFloat) && rightValue.IsInstanceOf(TFloat) {
		newValue := leftValue.Value.(float64) * rightValue.Value.(float64)
		return interpreter.NewFloat(newValue), nil
	}
	if leftValue.Type == rightValue.Type {
		return nil, fmt.Errorf("typeerror: type %v does not support operator * at line %v, col %v", leftValue.Type, tree.Line, tree.Col)
	}
	return nil, fmt.Errorf("typerror: incompatable types %v and %v with operator * at line %v, col %v", leftValue.Type, rightValue.Type, tree.Line, tree.Col)
}

func (interpreter *Interpreter) doDivision(tree *parser.Token) (*BeccaValue, error) {
	leftValue, rightValue, err := resolveBinaryOperands(interpreter, tree)
	if err != nil {
		return nil, err
	}
	if leftValue.IsInstanceOf(TInt) && rightValue.IsInstanceOf(TInt) {
		if rightValue.Value.(int64) == 0 {
			return nil, fmt.Errorf("dividebyzeroerror: integer division by zero at line %v, col %v", tree.Line, tree.Col)
		}
		newValue := leftValue.Value.(int64) / rightValue.Value.(int64)
		return interpreter.NewInt(newValue), nil
	}
	if leftValue.IsInstanceOf(TFloat) && rightValue.IsInstanceOf(TFloat) {
		if rightValue.Value.(int64) == 0.0 {
			return nil, fmt.Errorf("dividebyzeroerror: float division by zero at line %v, col %v", tree.Line, tree.Col)
		}
		newValue := leftValue.Value.(float64) / rightValue.Value.(float64)
		return interpreter.NewFloat(newValue), nil
	}
	if leftValue.Type == rightValue.Type {
		return nil, fmt.Errorf("typeerror: type %v does not support operator / at line %v, col %v", leftValue.Type, tree.Line, tree.Col)
	}
	return nil, fmt.Errorf("typerror: incompatable types %v and %v with operator / at line %v, col %v", leftValue.Type, rightValue.Type, tree.Line, tree.Col)
}

func (interpreter *Interpreter) doModulo(tree *parser.Token) (*BeccaValue, error) {
	leftValue, rightValue, err := resolveBinaryOperands(interpreter, tree)
	if err != nil {
		return nil, err
	}
	if leftValue.IsInstanceOf(TInt) && rightValue.IsInstanceOf(TInt) {
		if rightValue.Value.(int64) == 0 {
			return nil, fmt.Errorf("dividebyzeroerror: integer modulo by zero at line %v, col %v", tree.Line, tree.Col)
		}
		newValue := leftValue.Value.(int64) % rightValue.Value.(int64)
		return interpreter.NewInt(newValue), nil
	}
	if leftValue.Type == rightValue.Type {
		return nil, fmt.Errorf("typeerror: type %v does not support operator %% at line %v, col %v", leftValue.Type, tree.Line, tree.Col)
	}
	return nil, fmt.Errorf("typerror: incompatable types %v and %v with operator %% at line %v, col %v", leftValue.Type, rightValue.Type, tree.Line, tree.Col)
}
