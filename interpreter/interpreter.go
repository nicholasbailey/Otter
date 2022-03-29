package interpreter

import (
	"fmt"
	"strconv"

	"github.com/nicholasbailey/otter/exception"
	"github.com/nicholasbailey/otter/parser"
)

// TODO - There are a lot of magic strings here

type Scope map[string]*OtterValue

func NewScope() Scope {
	return map[string]*OtterValue{}
}

type Interpreter struct {
	CallStack CallStack
}

func (interpreter *Interpreter) Execute(statements []*parser.Token) (*OtterValue, exception.Exception) {
	var value *OtterValue
	var err error = nil
	for _, statement := range statements {

		value, err = interpreter.Evaluate(statement)
		if err != nil {
			return nil, err
		}
	}
	return value, nil
}

func (interpreter *Interpreter) Evaluate(tree *parser.Token) (*OtterValue, exception.Exception) {
	switch tree.Symbol {
	case parser.StringLiteral:
		return interpreter.NewString(tree.Value), nil
	case parser.IntLiteral:
		parsedInt, err := strconv.ParseInt(tree.Value, 0, 64)
		if err != nil {
			// TODO: Make this a proper exception
			return nil, err
		}
		return interpreter.NewInt(parsedInt), nil
	case parser.FloatLiteral:
		parsedFloat, err := strconv.ParseFloat(tree.Value, 64)
		if err != nil {
			return nil, err
		}
		return interpreter.NewFloat(parsedFloat), nil
	case "true":
		return interpreter.True(), nil
	case "false":
		return interpreter.False(), nil
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
	case parser.Assignment:
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
	case parser.While:
		return interpreter.doWhile(tree)
	case parser.FunctionDefinition:
		return interpreter.defineFunction(tree)
	case parser.FunctionInvocation:
		return interpreter.callFunction(tree)
	case parser.Block:
		var result *OtterValue
		var err exception.Exception
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
			return nil, exception.New(exception.SyntaxError, "illegal return in global scope", tree.Line, tree.Col)
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
	case parser.Access:
		return interpreter.doAccess(tree)
	}

	return nil, fmt.Errorf("syntaxerror: unrecognized symbol '%v' at line %v, col %v", tree.Value, tree.Line, tree.Col)
}

func (interpreter *Interpreter) DefineGlobal(name string, value *OtterValue) {
	interpreter.CallStack.Globals().Scope[name] = value
}

func (interpreter *Interpreter) DefineMethod(typeName TypeName, methodName string, callable *Callable) {
	typeVal := interpreter.MustResolveType(typeName)
	typeVal.Methods[methodName] = callable
}

func (interpreter *Interpreter) DefineBuiltinMethod(
	typeName TypeName,
	methodName string,
	arity int,
	builtInFunction BuiltInFunction,
) {
	methodFn, _ := interpreter.NewBuiltInFunction(methodName, arity, builtInFunction)
	interpreter.DefineMethod(typeName, methodName, methodFn.Callable)
}

func NewInterpreter() *Interpreter {
	interpreter := &Interpreter{
		CallStack: *NewCallStack(),
	}
	globalFrame := NewCallStackFrame("global")
	interpreter.CallStack.Push(globalFrame)
	DefineTypeType(interpreter)

	// Define built in types
	DefineStringTypes(interpreter)
	// DefineIntType(interpreter)
	// DefineBoolType(interpreter)
	// DefineFloatType(interpreter)
	// DefineNullType(interpreter)
	// DefineFunctionType(interpreter)
	// DefineListType(interpreter)

	interpreter.DefineType(TInt, NewBuiltInConstructor(TString, 1, ConstructInt))
	interpreter.DefineType(TFloat, NewBuiltInConstructor(TFloat, 1, ConstructFloat))
	interpreter.DefineType(TBool, NewBuiltInConstructor(TBool, 1, ConstructBool))
	interpreter.DefineType(TNull, NewBuiltInConstructor(TNull, 0, ConstructNull))

	interpreter.DefineType(TFunction, NewBuiltInConstructor("function", 0, ConstructFunction))
	interpreter.DefineGlobal("true", interpreter.True())
	interpreter.DefineGlobal("false", interpreter.False())
	interpreter.DefineGlobal("null", interpreter.NewNull())
	DefineBuiltins(interpreter)

	return interpreter
}
