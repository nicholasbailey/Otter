package interpreter

import (
	"fmt"
	"strconv"

	"github.com/nicholasbailey/becca/common"
	"github.com/nicholasbailey/becca/parser"
)

// TODO - There are a lot of magic strings here

type BeccaValue struct {
	Type               BeccaType
	Value              interface{}
	FunctionDefinition *parser.Token
}

func Null() *BeccaValue {
	return &BeccaValue{
		Type:  TNull,
		Value: nil,
	}
}

type Scope map[string]*BeccaValue

func NewScope() Scope {
	return map[string]*BeccaValue{}
}

type BuiltInFunction func(values []*BeccaValue) (*BeccaValue, common.Exception)

type Interpreter struct {
	CallStack CallStack
	BuiltIns  map[string]BuiltInFunction
}

func (interpreter *Interpreter) Execute(statements []*parser.Token) (*BeccaValue, common.Exception) {
	var value *BeccaValue
	var err error = nil
	for _, statement := range statements {

		value, err = interpreter.Evaluate(statement)
		if err != nil {
			return nil, err
		}
	}
	return value, nil
}

func (interpreter *Interpreter) Evaluate(tree *parser.Token) (*BeccaValue, common.Exception) {
	switch tree.Symbol {
	case parser.StringLiteral:
		return &BeccaValue{
			Type:  TString,
			Value: tree.Value,
		}, nil
	case parser.IntLiteral:
		parsedInt, err := strconv.ParseInt(tree.Value, 0, 64)
		if err != nil {
			// TODO: Make this a proper exception
			return nil, err
		}
		return &BeccaValue{
			Type:  TInt,
			Value: parsedInt,
		}, nil
	case parser.FloatLiteral:
		parsedFloat, err := strconv.ParseFloat(tree.Value, 64)
		if err != nil {
			return nil, err
		}
		return &BeccaValue{
			Type:  TFloat,
			Value: parsedFloat,
		}, nil
	case "true":
		return True(), nil
	case "false":
		return False(), nil
	case parser.Name:
		value, found := interpreter.CallStack.ResolveVariable(tree.Value)
		if !found {
			return nil, fmt.Errorf("syntaxerror: unbound variable %v at line %v, col %v", tree.Value, tree.Line, tree.Col)
		}
		return value, nil
	// Handle Variable assignment
	case "&&":
		return interpreter.doAnd(tree)
	case "||":
		return interpreter.doOr(tree)
	case "!=":
		return interpreter.doInequalityCheck(tree)
	case "==":
		return interpreter.doEqualityCheck(tree)
	case "=":
		return interpreter.doAssigment(tree)
	case "+":
		return interpreter.doAddition(tree)
	case "-":
		return interpreter.doSubtraction(tree)
	case "*":
		return interpreter.doMultiplication(tree)
	case "/":
		return interpreter.doDivision(tree)
	case "%":
		return interpreter.doModulo(tree)
	case "<":
		return interpreter.doLessThan(tree)
	case ">":
		return interpreter.doGreaterThan(tree)
	case "<=":
		return interpreter.doLessThanOrEqualTo(tree)
	case ">=":
		return interpreter.doGreaterThanOrEqualTo(tree)
	case "while":
		return interpreter.doWhile(tree)
	case parser.FunctionDefinition:
		return interpreter.defineFunction(tree)
	case parser.FunctionInvocation:
		return interpreter.callFunction(tree)
	case parser.Block:
		var result *BeccaValue
		var err common.Exception
		for _, child := range tree.Children {
			result, err = interpreter.Evaluate(child)
			if err != nil {
				return nil, err
			}
		}
		return result, nil
	case "return":
		stackFrame := interpreter.CallStack.Peek()
		if stackFrame.FunctionName == "global" {
			return nil, common.NewException(common.SyntaxError, "illegal return in global scope", tree.Line, tree.Col)
		}
		child := tree.Children[0]
		value, err := interpreter.Evaluate(child)
		// TODO - stack handle errors
		if err != nil {
			return nil, err
		}
		stackFrame.ReturnValue = value
		return value, nil
	case "if":
		return interpreter.doIf(tree)
	}

	return nil, fmt.Errorf("syntaxerror: unrecognized symbol '%v' at line %v, col %v", tree.Value, tree.Line, tree.Col)
}

func NewInterpreter() *Interpreter {
	print := func(values []*BeccaValue) (*BeccaValue, common.Exception) {
		for _, value := range values {
			switch value.Type {
			case TString:
				fmt.Print(value.Value.(string))
			case TInt:
				// TODO - move away from builtin
				fmt.Print(value.Value.(int64))
			case TBool:
				fmt.Print(value.Value.(bool))
			case TFloat:
				fmt.Print(value.Value.(float64))
			case TNull:
				fmt.Print("<null>")
			}
			fmt.Print(" ")
		}
		fmt.Print("\n")
		return Null(), nil
	}
	interpreter := &Interpreter{
		BuiltIns:  map[string]BuiltInFunction{},
		CallStack: *NewCallStack(),
	}
	globalFrame := NewCallStackFrame("global")
	globalFrame.Scope["true"] = True()
	globalFrame.Scope["false"] = False()
	interpreter.CallStack.Push(globalFrame)
	interpreter.BuiltIns["print"] = print
	interpreter.BuiltIns["string"] = String
	return interpreter
}
