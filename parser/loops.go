package parser

import (
	"fmt"

	"github.com/nicholasbailey/otter/exception"
)

func (spec *LanguageSpecification) DefineForIn(forKeyword Symbol, inKeyword Symbol) {
	forStd := func(token *Token, parser *TDOPParser) (*Token, exception.Exception) {
		loopVarToken, err := parser.Next()
		if err != nil {
			return nil, err
		}
		if loopVarToken.Symbol != Name {
			errorMsg := fmt.Sprintf("Unexpected symbol %v in for expression", loopVarToken.Value)
			return nil, exception.New(exception.SyntaxError, errorMsg, loopVarToken.Line, loopVarToken.Col)
		}
		// TODO support tuply syntax here
		inToken, err := parser.Next()
		if err != nil {
			return nil, err
		}
		if inToken.Symbol != inKeyword {
			return nil, exception.New(exception.SyntaxError, "Expected 'in'", inToken.Line, inToken.Col)
		}
		rangeToken, err := parser.Expression(0)
		if err != nil {
			return nil, err
		}
		blockToken, err := parser.Block()
		if err != nil {
			return nil, err
		}

		token.Symbol = ForIn
		token.Children = append(token.Children, loopVarToken)
		token.Children = append(token.Children, rangeToken)
		token.Children = append(token.Children, blockToken)
		return token, nil
	}
	spec.DefineEmpty(inKeyword)
	spec.DefineStatment(forKeyword, forStd)
}

func (spec *LanguageSpecification) DefineWhile(whileKeyword Symbol) {
	whileStd := func(token *Token, parser *TDOPParser) (*Token, exception.Exception) {
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

	spec.DefineStatment(whileKeyword, whileStd)
}
