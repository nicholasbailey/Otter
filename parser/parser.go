package parser

import (
	"io"

	"github.com/nicholasbailey/becca/exception"
)

// A Parser is an object that generates a Becca abstract
// syntax tree from a io.Reader stream of source code.
// It exposes a single method, Statements which converts
// the entire source stream into a slice of ASTs representing
// the statments in the source file. If the source code is
// syntactically incorrect in any way, it returns an exception
// value.
type Parser interface {
	Statements() ([]*Token, exception.Exception)
}

type BeccaParser struct {
	BaseParser  Parser
	Unsweetener Unsweetener
}

func NewParser(source io.Reader) Parser {
	language := NewBeccaLanguage()
	lexer := NewLexer(source, language)
	baseParser := NewTDOPParser(lexer)
	unsweetener := NewUnsweetener()

	return &BeccaParser{
		BaseParser:  baseParser,
		Unsweetener: unsweetener,
	}
}

func (beccaParser *BeccaParser) Statements() ([]*Token, exception.Exception) {
	statements, err := beccaParser.BaseParser.Statements()
	if err != nil {
		return nil, err
	}
	newStatements := []*Token{}
	for _, statement := range statements {
		unsweetened, err := beccaParser.Unsweetener.Unsweeten(statement)
		if err != nil {
			return nil, err
		}
		newStatements = append(newStatements, unsweetened)
	}
	return newStatements, nil
}
