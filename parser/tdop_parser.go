package parser

import (
	"fmt"

	"github.com/nicholasbailey/otter/exception"
)

// A TDOPParser is the core parser implentation for the Otter language
// it uses Top-Down Operator Precedence Parsing to obtain a sequence
// of statement values.
type TDOPParser struct {
	Lexer *Lexer
}

// Factory function for a TDOPParser
func NewTDOPParser(lexer *Lexer) Parser {
	return &TDOPParser{
		Lexer: lexer,
	}
}

func (parser *TDOPParser) IsStatementTerminator(token *Token) bool {
	return parser.Lexer.IsStatementTerminator(token)
}

func (parser *TDOPParser) Next() (*Token, error) {
	return parser.Lexer.Next()
}

func (parser *TDOPParser) Peek() (*Token, error) {
	return parser.Lexer.Peek()
}

func (parser *TDOPParser) Block() (*Token, error) {
	token, err := parser.Lexer.Next()
	if err != nil {
		return nil, err
	}
	if !parser.Lexer.IsBlockStart(token) {
		return nil, fmt.Errorf("syntaxerror: expected block start, but got %v at line %v, col %v", token.Value, token.Line, token.Col)
	}

	block, err := token.Std(token, parser)
	if err != nil {
		return nil, err
	}
	parser.Lexer.Next()
	return block, nil
}

func (parser *TDOPParser) Statement() (*Token, error) {

	tok, err := parser.Lexer.Peek()
	if err != nil {
		return nil, err
	}
	if tok.Std != nil {
		tok, err = parser.Lexer.Next()
		if err != nil {
			return nil, err
		}
		return tok.Std(tok, parser)
	}
	res, err := parser.Expression(0)
	if err != nil {
		return nil, err
	}
	terminator, err := parser.Lexer.Next()
	if !parser.Lexer.IsStatementTerminator(terminator) {
		return nil, exception.New(exception.SyntaxError, fmt.Sprintf("unterminated statement with %v", terminator.Value), terminator.Line, terminator.Col)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (parser *TDOPParser) Statements() ([]*Token, exception.Exception) {
	statements := []*Token{}
	next, err := parser.Lexer.Peek()
	if err != nil {
		return nil, err
	}
	for next.Symbol != EOF && !parser.Lexer.IsAnyBlockEnd(next) {
		statement, err := parser.Statement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, statement)
		next, err = parser.Lexer.Peek()
		if err != nil {
			return nil, err
		}
	}
	return statements, nil
}

func (parser *TDOPParser) Expression(rightBindingPower int) (*Token, exception.Exception) {
	var left *Token

	t, err := parser.Lexer.Next()
	if err != nil {
		return nil, err
	}
	if t.Nud == nil {
		return nil, exception.New(exception.SyntaxError, fmt.Sprintf("%v is not a valid prefix symbol", t.Value), t.Line, t.Col)
	}
	left, err = t.Nud(t, parser)
	if err != nil {
		return nil, err
	}
	for {
		peek, err := parser.Lexer.Peek()
		if err != nil {
			return nil, err
		}
		if rightBindingPower >= peek.BindingPower {
			break
		}
		t, err := parser.Lexer.Next()
		if err != nil {
			return nil, err
		}
		if t.Led == nil {
			return nil, fmt.Errorf("%v is not a valid infix symbol", t.Symbol)
		}
		left, err = t.Led(t, parser, left)
		if err != nil {
			return nil, err
		}
	}
	return left, nil
}
