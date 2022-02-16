package parser

import (
	"fmt"

	"github.com/nicholasbailey/becca/exception"
)

func (spec *LanguageSpecification) DefineFunctionDefinition(defSymbol Symbol) {
	defStd := func(token *Token, parser *TDOPParser) (*Token, exception.Exception) {
		token.Symbol = FunctionDefinition
		functionName, err := parser.Next()
		if err != nil {
			return nil, err
		}
		if functionName.Symbol != Name {
			return nil, exception.New(exception.SyntaxError, fmt.Sprintf("expected identifier, got %v", functionName.Value), functionName.Line, functionName.Col)
		}
		token.Children = append(token.Children, functionName)

		openParens, err := parser.Next()
		if err != nil {
			return nil, err
		}
		if openParens.Symbol != "(" {
			return nil, exception.New(exception.SyntaxError, fmt.Sprintf("expected (, got %v", openParens.Value), openParens.Line, openParens.Col)
		}
		parameters := []*Token{}
		next, err := parser.Next()
		if err != nil {
			return nil, err
		}
		if next.Symbol != ")" {
			for {
				if next.Symbol != Name {
					return nil, err
				}
				parameters = append(parameters, next)
				further, err := parser.Peek()
				if err != nil {
					return nil, err
				}
				if further.Symbol != "," {
					break
				}
				_, err = parser.Next()
				if err != nil {
					return nil, err
				}
				next, err = parser.Next()
				if err != nil {
					return nil, err
				}
			}
			close, err := parser.Next()
			if err != nil {
				return nil, err
			}
			if close.Symbol != ")" {
				return nil, fmt.Errorf("syntaxerror: unterminated parentheses with symbol %v at line %v, col %v", close.Value, close.Line, close.Col)
			}
		} else {
			_, err = parser.Next()
			if err != nil {
				return nil, err
			}
		}
		parameterToken := &Token{
			Symbol:   FunctionParameters,
			Value:    "(",
			Line:     openParens.Line,
			Col:      openParens.Col,
			Children: parameters,
		}
		token.Children = append(token.Children, parameterToken)
		block, err := parser.Block()
		if err != nil {
			return nil, err
		}
		token.Children = append(token.Children, block)
		return token, nil
	}

	spec.DefineStatment(defSymbol, defStd)
}

func (spec *LanguageSpecification) DefineReturn(returnSymbol Symbol) {
	returnStd := func(token *Token, parser *TDOPParser) (*Token, exception.Exception) {
		expression, err := parser.Expression(0)
		if err != nil {
			return nil, err
		}
		token.Children = append(token.Children, expression)
		// Hack, something is wonky here
		next, err := parser.Peek()
		if err != nil {
			return nil, err
		}
		if parser.IsStatementTerminator(next) {
			_, err = parser.Next()
			if err != nil {
				return nil, err
			}
		}
		return token, nil
	}

	spec.DefineStatment(returnSymbol, returnStd)
}
