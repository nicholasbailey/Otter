package main

import (
	"fmt"
	"os"

	"github.com/nicholasbailey/becca/interpreter"
	"github.com/nicholasbailey/becca/parser"
)

func main() {
	// TODO: fix this

	path := os.Args[len(os.Args)-1]
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	if os.Args[1] == "--raw-syntax" {
		spec := parser.NewBeccaLanguage()
		lexer := parser.NewLexer(file, spec)
		parser := parser.NewTDOPParser(lexer)
		tokens, err := parser.Statements()
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		} else {
			for _, token := range tokens {
				fmt.Printf("%v\n", token.TreeString(0))
			}
			os.Exit(0)
		}
	} else if os.Args[1] == "--unsweetened-syntax" {
		parser := parser.NewParser(file)
		tokens, err := parser.Statements()
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		} else {
			for _, token := range tokens {
				fmt.Printf("%v\n", token.TreeString(0))
			}
			os.Exit(0)
		}
	}
	engine := interpreter.NewEngine()
	_, err = engine.Execute(file)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
