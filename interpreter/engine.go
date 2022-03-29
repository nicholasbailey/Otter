package interpreter

import (
	"io"

	"github.com/nicholasbailey/otter/exception"
	"github.com/nicholasbailey/otter/parser"
)

func NewEngine() *Engine {
	interpreter := NewInterpreter()
	parserFactory := func(source io.Reader) parser.Parser {
		return parser.NewParser(source)
	}
	return &Engine{
		ParserFactory: parserFactory,
		Interpreter:   *interpreter,
	}
}

type Engine struct {
	ParserFactory func(io.Reader) parser.Parser
	Interpreter   Interpreter
}

func (engine *Engine) Execute(source io.Reader) (*OtterValue, exception.Exception) {
	parser := engine.ParserFactory(source)
	trees, err := parser.Statements()
	if err != nil {
		return nil, err
	}
	return engine.Interpreter.Execute(trees)
}
