package parser

import (
	"io"

	"github.com/nicholasbailey/otter/exception"
)

// A Parser is an object that generates a Otter abstract
// syntax tree from a io.Reader stream of source code.
// It exposes a single method, Statements which converts
// the entire source stream into a slice of ASTs representing
// the statments in the source file. If the source code is
// syntactically incorrect in any way, it returns an exception
// value.
type Parser interface {
	Statements() ([]*Token, exception.Exception)
}

type OtterParser struct {
	BaseParser  Parser
	Unsweetener Unsweetener
}

func NewParser(source io.Reader) Parser {
	language := NewOtterLanguage()
	lexer := NewLexer(source, language)
	baseParser := NewTDOPParser(lexer)
	unsweetener := NewUnsweetener()

	return &OtterParser{
		BaseParser:  baseParser,
		Unsweetener: unsweetener,
	}
}

func (otterParser *OtterParser) Statements() ([]*Token, exception.Exception) {
	statements, err := otterParser.BaseParser.Statements()
	if err != nil {
		return nil, err
	}
	newStatements := []*Token{}
	for _, statement := range statements {
		unsweetened, err := otterParser.Unsweetener.Unsweeten(statement)
		if err != nil {
			return nil, err
		}
		newStatements = append(newStatements, unsweetened)
	}
	return newStatements, nil
}
