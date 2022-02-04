package interpreter

import (
	"io"

	"github.com/nicholasbailey/becca/common"
	"github.com/nicholasbailey/becca/parser"
)

func NewEngine() *Engine {
	interpreter := NewInterpreter()
	languageSpec := parser.NewBeccaLanguage()
	lexerFactory := func(reader io.Reader) *parser.Lexer {
		return parser.NewLexer(reader, languageSpec)
	}
	return &Engine{
		LexerFactory: lexerFactory,
		Interpreter:  *interpreter,
	}
}

type Engine struct {
	LexerFactory func(io.Reader) *parser.Lexer
	Interpreter  Interpreter
}

func (engine *Engine) Execute(source io.Reader) (*BeccaValue, common.Exception) {
	lexer := engine.LexerFactory(source)
	parser := parser.Parser{
		Lexer: lexer,
	}

	trees, err := parser.Statements()
	if err != nil {
		return nil, err
	}
	return engine.Interpreter.Execute(trees)
}
