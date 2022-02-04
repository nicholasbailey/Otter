package parser

import (
	"fmt"

	"github.com/nicholasbailey/becca/common"
)

func NewBeccaLanguage() *LanguageSpecification {
	spec := NewLanguage()
	spec.DefineQuotes('"', '"', StringLiteral)
	spec.DefineQuotes('\'', '\'', StringLiteral)
	spec.DefineParens("(", ")")
	spec.DefinePrefix("!", 80)
	spec.DefineInfix("&&", 30)
	spec.DefineInfix("||", 20)
	spec.DefineInfix("=", 10)
	spec.DefineInfix("==", 50)
	spec.DefineInfix("!=", 50)
	spec.DefineInfix("<", 50)
	spec.DefineInfix(">", 50)
	spec.DefineInfix("<=", 50)
	spec.DefineInfix(">=", 50)
	spec.DefineInfix("+", 60)
	spec.DefineInfix("-", 60)
	spec.DefineInfix("*", 70)
	spec.DefineInfix("/", 70)
	spec.DefineInfix("%", 70)
	spec.DefineStatementTerminator(";")
	spec.DefineEmpty(",")
	spec.DefineBlock("{", "}")
	spec.DefineEmpty("else")
	spec.DefineValue("true")
	spec.DefineValue("false")
	ifStd := func(token *Token, parser *Parser) (*Token, common.Exception) {
		expression, err := parser.Expression(0)
		if err != nil {
			return nil, err
		}
		token.Children = append(token.Children, expression)

		block, err := parser.Block()
		if err != nil {
			return nil, err
		}
		token.Children = append(token.Children, block)
		next, err := parser.Lexer.Peek()
		if err != nil {
			return nil, err
		}
		// TODO: make this symbol references
		if next.Value == "else" {
			parser.Lexer.Next()
			next, err := parser.Lexer.Peek()
			if err != nil {
				return nil, err
			}
			if next.Value == "if" {
				stmt, err := parser.Statement()
				if err != nil {
					return nil, err
				}
				stmt.Symbol = ElseIf
				token.Children = append(token.Children, stmt)
			} else {
				block, err := parser.Block()
				if err != nil {
					return nil, err
				}
				token.Children = append(token.Children, block)
			}
		}
		return token, nil
	}

	ifNud := func(token *Token, parser *Parser) (*Token, common.Exception) {
		expression, err := parser.Expression(0)
		if err != nil {
			return nil, err
		}
		token.Children = append(token.Children, expression)

		block, err := parser.Block()
		if err != nil {
			return nil, err
		}
		token.Children = append(token.Children, block)
		next, err := parser.Lexer.Peek()
		if err != nil {
			return nil, err
		}
		// TODO: make this symbol references
		if next.Value == "else" {
			parser.Lexer.Next()
			next, err := parser.Lexer.Peek()
			if err != nil {
				return nil, err
			}
			if next.Value == "if" {
				stmt, err := parser.Statement()
				if err != nil {
					return nil, err
				}
				token.Children = append(token.Children, stmt)
			} else {
				block, err := parser.Block()
				if err != nil {
					return nil, err
				}
				token.Children = append(token.Children, block)
			}
		}
		return token, nil
	}

	spec.Define("if", 0, 0, ifNud, nil, ifStd)

	returnStd := func(token *Token, parser *Parser) (*Token, common.Exception) {
		expression, err := parser.Expression(0)
		if err != nil {
			return nil, err
		}
		token.Children = append(token.Children, expression)
		// Hack, something is wonky here
		next, err := parser.Lexer.Peek()
		if err != nil {
			return nil, err
		}
		if parser.Lexer.IsStatementTerminator(next) {
			_, err = parser.Lexer.Next()
			if err != nil {
				return nil, err
			}
		}
		return token, nil
	}

	spec.DefineStatment("return", returnStd)

	whileStd := func(token *Token, parser *Parser) (*Token, common.Exception) {
		expression, err := parser.Expression(0)
		if err != nil {
			return nil, err
		}

		token.Children = append(token.Children, expression)
		block, err := parser.Block()
		if err != nil {
			return nil, err
		}
		token.Children = append(token.Children, block)
		return token, nil
	}

	spec.DefineStatment("while", whileStd)

	openParensLed := func(right *Token, parser *Parser, left *Token) (*Token, common.Exception) {
		if left.Symbol != Name && left.Symbol != Symbol("(") {
			return nil, fmt.Errorf("syntaxerror: unexpected ( at line %v, col %v", right.Line, right.Col)
		}
		right.Children = append(right.Children, left)
		t, err := parser.Lexer.Peek()
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
				right, err := parser.Lexer.Peek()
				if err != nil {
					return nil, err
				}
				if right.Symbol != "," {
					break
				}
				_, err = parser.Lexer.Next()
				if err != nil {
					return nil, err
				}
			}
			close, err := parser.Lexer.Next()
			if err != nil {
				return nil, err
			}
			if close.Symbol != ")" {
				return nil, fmt.Errorf("syntaxerror: unterminated parentheses with symbol %v at line %v, col %v", close.Value, close.Line, close.Col)
			}
		} else {
			_, err = parser.Lexer.Next()
			if err != nil {
				return nil, err
			}
		}
		right.Symbol = FunctionInvocation
		return right, nil
	}

	spec.Define("(", 90, 0, nil, openParensLed, nil)

	defStd := func(token *Token, parser *Parser) (*Token, common.Exception) {
		token.Symbol = FunctionDefinition
		functionName, err := parser.Lexer.Next()
		if err != nil {
			return nil, err
		}
		if functionName.Symbol != Name {
			return nil, common.NewException(common.SyntaxError, fmt.Sprintf("expected identifier, got %v", functionName.Value), functionName.Line, functionName.Col)
		}
		token.Children = append(token.Children, functionName)

		openParens, err := parser.Lexer.Next()
		if err != nil {
			return nil, err
		}
		if openParens.Symbol != "(" {
			return nil, common.NewException(common.SyntaxError, fmt.Sprintf("expected (, got %v", openParens.Value), openParens.Line, openParens.Col)
		}
		parameters := []*Token{}
		next, err := parser.Lexer.Next()
		if err != nil {
			return nil, err
		}
		if next.Symbol != ")" {
			for {
				if next.Symbol != Name {
					return nil, err
				}
				parameters = append(parameters, next)
				further, err := parser.Lexer.Peek()
				if err != nil {
					return nil, err
				}
				if further.Symbol != "," {
					break
				}
				_, err = parser.Lexer.Next()
				if err != nil {
					return nil, err
				}
				next, err = parser.Lexer.Next()
				if err != nil {
					return nil, err
				}
			}
			close, err := parser.Lexer.Next()
			if err != nil {
				return nil, err
			}
			if close.Symbol != ")" {
				return nil, fmt.Errorf("syntaxerror: unterminated parentheses with symbol %v at line %v, col %v", close.Value, close.Line, close.Col)
			}
		} else {
			_, err = parser.Lexer.Next()
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

	spec.DefineStatment("def", defStd)

	return spec
}
