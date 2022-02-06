package parser

import "github.com/nicholasbailey/becca/exception"

func (spec *LanguageSpecification) DefineWhile(whileKeyword Symbol) {
	whileStd := func(token *Token, parser *Parser) (*Token, exception.Exception) {
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
		token.Symbol = While
		return token, nil
	}

	spec.DefineStatment(Symbol(whileKeyword), whileStd)
}
