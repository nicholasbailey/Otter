package parser

import "github.com/nicholasbailey/becca/exception"

func (spec *LanguageSpecification) DefineAccess(accessSymbol Symbol) {
	dotLed := func(token *Token, parser *Parser, left *Token) (*Token, exception.Exception) {
		next, err := parser.Lexer.Peek()
		if err != nil {
			return nil, err
		}
		if next.Symbol != Name {
			return nil, exception.New(exception.SyntaxError, "invalid property access", token.Line, token.Col)
		}

		token.Symbol = Access
		token.Children = append(token.Children, left)
		exp, err := parser.Expression(token.BindingPower)
		if err != nil {
			return nil, err
		}
		token.Children = append(token.Children, exp)
		return token, nil
	}

	spec.Define(accessSymbol, 100, 2, nil, dotLed, nil)
}
