package parser

import (
	"fmt"

	"github.com/nicholasbailey/otter/exception"
)

func (spec *LanguageSpecification) DefineParens(openParens Symbol, closeParens Symbol) {

	nud := func(t *Token, p *TDOPParser) (*Token, exception.Exception) {
		expressionToken, err := p.Expression(0)
		if err != nil {
			return nil, err
		}
		next, err := p.Next()
		if err != nil {
			return nil, err
		}
		if next.Symbol != closeParens {
			return nil, fmt.Errorf("SyntaxError: Unterminated braces at line %v, col %v", next.Line, next.Col)
		}
		return expressionToken, nil
	}
	spec.Define(openParens, 0, 0, nud, nil, nil)
	spec.DefineValue(closeParens)

	openParensLed := func(right *Token, parser *TDOPParser, left *Token) (*Token, exception.Exception) {
		if left.Symbol != Name && left.Symbol != Symbol("(") {
			return nil, fmt.Errorf("syntaxerror: unexpected ( at line %v, col %v", right.Line, right.Col)
		}
		right.Children = append(right.Children, left)
		t, err := parser.Peek()
		if err != nil {
			return nil, err
		}
		if t.Symbol != ")" {
			for {
				expressionResult, err := parser.Expression(0)
				if err != nil {
					return nil, err
				}
				right.Children = append(right.Children, expressionResult)
				right, err := parser.Peek()
				if err != nil {
					return nil, err
				}
				if right.Symbol != "," {
					break
				}
				_, err = parser.Next()
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
		right.Symbol = FunctionInvocation
		return right, nil
	}

	spec.Define("(", 110, 0, nil, openParensLed, nil)

}
