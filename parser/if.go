package parser

import "github.com/nicholasbailey/becca/exception"

// Contains parser logic for if statements

func (spec *LanguageSpecification) DefineIf(ifKeyword Symbol, elseKeyword Symbol) {

	ifStd := func(token *Token, parser *TDOPParser) (*Token, exception.Exception) {
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
		next, err := parser.Peek()
		if err != nil {
			return nil, err
		}
		// TODO: make this symbol references
		if next.Value == string(elseKeyword) {
			parser.Next()
			next, err := parser.Peek()
			if err != nil {
				return nil, err
			}
			if next.Value == string(ifKeyword) {
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

	ifNud := func(token *Token, parser *TDOPParser) (*Token, exception.Exception) {
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
		next, err := parser.Peek()
		if err != nil {
			return nil, err
		}
		// TODO: make this symbol references
		if next.Value == string(elseKeyword) {
			parser.Next()
			next, err := parser.Peek()
			if err != nil {
				return nil, err
			}
			if next.Value == string(ifKeyword) {
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

	spec.Define(ifKeyword, 0, 0, ifNud, nil, ifStd)
	spec.DefineEmpty(elseKeyword)
}
