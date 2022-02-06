package parser

import (
	"unicode"

	"github.com/nicholasbailey/becca/exception"
)

// TODO: this needs a good bit of refactoring

func NewLanguage() *LanguageSpecification {
	symbols := make(map[Symbol]*Token)
	quotes := make(map[rune]*quoteSpecification)
	language := &LanguageSpecification{
		quoteDefinitions:     quotes,
		symbols:              symbols,
		statementTerminators: []Symbol{},
		blockDelimiters:      map[Symbol]Symbol{},
		commentStarts:        []Symbol{},
	}
	language.DefineValue(Name)
	language.DefineValue(IntLiteral)
	language.DefineValue(FloatLiteral)
	return language
}

type quoteSpecification struct {
	openQuote   rune
	closeQuote  rune
	literalType Symbol
}

type LanguageSpecification struct {
	quoteDefinitions     map[rune]*quoteSpecification
	symbols              map[Symbol]*Token
	statementTerminators []Symbol
	blockDelimiters      map[Symbol]Symbol
	commentStarts        []Symbol
}

func (spec *LanguageSpecification) DefineComment(symbol Symbol) {
	spec.commentStarts = append(spec.commentStarts, symbol)
}

func (spec *LanguageSpecification) IsCommentStart(symbol Symbol) bool {
	for _, start := range spec.commentStarts {
		if symbol == start {
			return true
		}
	}
	return false
}

func (spec *LanguageSpecification) IsAnyBlockEnd(symbol Symbol) bool {
	for _, v := range spec.blockDelimiters {
		if v == symbol {
			return true
		}
	}
	return false
}

func (spec *LanguageSpecification) DefineBlock(startSymbol Symbol, endSymbol Symbol) {
	spec.DefineEmpty(endSymbol)

	std := func(token *Token, parser *Parser) (*Token, exception.Exception) {
		statements, err := parser.Statements()
		if err != nil {
			return nil, err
		}
		token.Children = append(token.Children, statements...)
		token.Symbol = Block
		return token, nil
	}

	spec.DefineStatment(startSymbol, std)
	spec.blockDelimiters[startSymbol] = endSymbol
}

func (spec *LanguageSpecification) IsBlockStart(symbol Symbol) bool {
	if _, found := spec.blockDelimiters[symbol]; found {
		return true
	}
	return false
}
func (spec *LanguageSpecification) IsBlockEnd(symbol Symbol, startSymbol Symbol) bool {
	end, found := spec.blockDelimiters[startSymbol]
	if !found {
		return false
	}
	return symbol == end
}

func (spec *LanguageSpecification) DefineEmpty(symbol Symbol) {
	spec.Define(symbol, 0, 0, nil, nil, nil)
}

func (spec *LanguageSpecification) IsStatementTerminator(symbol Symbol) bool {
	for _, s := range spec.statementTerminators {
		if s == symbol {
			return true
		}
	}
	return false
}

func (spec *LanguageSpecification) DefineStatementTerminator(symbol Symbol) {
	spec.statementTerminators = append(spec.statementTerminators, symbol)
	spec.symbols[symbol] = &Token{
		Symbol:       symbol,
		BindingPower: 0,
		Nud:          nil,
		Led:          nil,
		Std:          nil,
		Children:     []*Token{},
	}
}

func (spec *LanguageSpecification) IsIdentifierStartChararacter(char rune) bool {
	return spec.IsIdentifierCharacter(char) && !unicode.IsDigit(char)
}

func (spec *LanguageSpecification) IsIdentifierCharacter(char rune) bool {
	_, found := spec.symbols[Symbol(string(char))]
	if found {
		return false
	} else {
		return !unicode.IsSpace(char)
	}
}

func (spec *LanguageSpecification) GetQuoteSpec(openQuote rune) *quoteSpecification {
	val, _ := spec.quoteDefinitions[openQuote]
	return val
}

func (spec *LanguageSpecification) IsDefined(symbol Symbol) bool {
	_, present := spec.symbols[symbol]
	return present
}

func (spec *LanguageSpecification) Define(symbol Symbol, bindingPower int, arity int, nud NudFunction, led LedFunction, std StdFunction) {
	existing, found := spec.symbols[symbol]
	if found {
		if nud != nil && existing.Nud == nil {
			existing.Nud = nud
		}
		if led != nil && existing.Led == nil {
			existing.Led = led
		}
		if std != nil && existing.Std == nil {
			existing.Std = std
		}
		if bindingPower > existing.BindingPower {
			existing.BindingPower = bindingPower
		}
	} else {
		token := Token{
			Symbol:       symbol,
			Value:        "",
			BindingPower: bindingPower,
			Nud:          nud,
			Led:          led,
			Std:          std,
			Children:     []*Token{},
			Col:          -1,
			Line:         -1,
		}
		spec.symbols[symbol] = &token
	}
}

func (spec *LanguageSpecification) GenerateToken(symbol Symbol, value string, line int, col int) *Token {
	tok, present := spec.symbols[symbol]
	if !present {
		return nil
	}
	return &Token{
		Symbol:       symbol,
		Value:        value,
		BindingPower: tok.BindingPower,
		Nud:          tok.Nud,
		Led:          tok.Led,
		Std:          tok.Std,
		Children:     []*Token{},
		Col:          col,
		Line:         line,
	}
}

func (spec *LanguageSpecification) Eof(line int, col int) *Token {
	return &Token{
		Symbol:       EOF,
		Value:        "",
		BindingPower: 0,
		Nud:          nil,
		Led:          nil,
		Children:     []*Token{},
		Col:          col,
		Line:         line,
	}
}

func (spec *LanguageSpecification) DefineQuotes(openQuote rune, closeQuote rune, literalType Symbol) {
	spec.quoteDefinitions[openQuote] = &quoteSpecification{
		openQuote:   openQuote,
		closeQuote:  closeQuote,
		literalType: literalType,
	}
	spec.DefineValue(literalType)
}

func (spec *LanguageSpecification) DefineInfix(symbol Symbol, bindingPower int) {
	led := func(t *Token, parser *Parser, left *Token) (*Token, exception.Exception) {
		t.Children = append(t.Children, left)
		exprResult, err := parser.Expression(t.BindingPower)
		if err != nil {
			return nil, err
		}
		t.Children = append(t.Children, exprResult)
		return t, nil
	}
	spec.Define(symbol, bindingPower, 2, nil, led, nil)
}

func (spec *LanguageSpecification) DefinePrefix(symbol Symbol, bindingPower int) {
	nud := func(t *Token, parser *Parser) (*Token, exception.Exception) {
		expResult, err := parser.Expression(bindingPower)
		if err != nil {
			return nil, err
		}
		t.Children = append(t.Children, expResult)
		return t, nil
	}
	spec.Define(symbol, bindingPower, 1, nud, nil, nil)
}

// come up with a better name for this
func (spec *LanguageSpecification) DefineValue(symbol Symbol) {
	nud := func(t *Token, p *Parser) (*Token, exception.Exception) {
		return t, nil
	}
	spec.Define(symbol, 0, 0, nud, nil, nil)
}

func (spec *LanguageSpecification) DefineStatment(symbol Symbol, std StdFunction) {
	spec.Define(symbol, 0, 0, nil, nil, std)
}
