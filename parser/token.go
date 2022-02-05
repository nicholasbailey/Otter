package parser

import (
	"fmt"
	"strings"

	"github.com/nicholasbailey/becca/common"
)

// Represents the 'type' of a token. Symbols determine
// how the interpreter interacts with a node in the abstract
// syntax tree
type Symbol string

const (
	// Symbol for names (variables, function names)
	Name Symbol = "(NAME)"
	// Symbol for the token returned when the lexer reaches the end of a file
	EOF Symbol = "(EOF)"
	// Symbol for a string literal
	StringLiteral Symbol = "(STRING)"
	// Symbol for an integer literal
	IntLiteral Symbol = "(INT)"
	// Symbol for an float literal
	FloatLiteral Symbol = "(FLOAT)"
	// Symbol for a block
	Block Symbol = "(BLOCK)"
	// Symbol for an if statement
	If Symbol = "(IF)"
	// Symbol for an else if statement
	ElseIf Symbol = "(ELSEIF)"
	// Symbol for an else statement
	Else Symbol = "(ELSE)"
	// Symbol for a function definition
	FunctionDefinition Symbol = "(FUNCTIONDEFINITION)"
	// Symbol for a function parameters
	FunctionParameters Symbol = "(FUNCTIONPARAMETERS)"
	// Symbol for a function invocation
	FunctionInvocation Symbol = "(FUNCTIONINVOCATION)"
	Access             Symbol = "(ACCESS)"
)

type NudFunction func(right *Token, parser *Parser) (*Token, common.Exception)
type LedFunction func(right *Token, parser *Parser, left *Token) (*Token, common.Exception)
type StdFunction func(*Token, *Parser) (*Token, common.Exception)

// The 'Token' is the core data type of the parser
// A token is overloaded - it's both a token emitted
// by the lexer/tokenizer and a node in the post-parse
// abstract syntax tree. This double duty gives us a lot of flexibility,
// because large amounts of information about the lexing/parsing process
// remain embedded in the final AST
type Token struct {
	// The Symbol of the token
	Symbol Symbol
	// The raw string parsed into this token
	Value string
	// The binding power of this token. See parser/parser.go for how
	// this works
	BindingPower int
	// The line on which this token started
	Line int
	// The column at which this token started
	Col int
	// The children of this token in the abstract syntax tree
	Children []*Token
	// The Null Denotation function of the token. This describes
	// the token's behavior when it is a prefix of an expression
	Nud NudFunction
	// The Left Denotation function of the token. This describes the
	// token's behavior when it is an infix of an expression
	Led LedFunction
	// The Statement Denotation function of the token. This describes the
	// token's behavior when it is a
	Std StdFunction
}

// Provides a friendly, human readable version
// of the AST for a token
func (token *Token) TreeString(indentLevel int) string {
	var builder strings.Builder
	for i := 0; i < indentLevel; i++ {
		builder.WriteString("  ")
	}

	builder.WriteString(fmt.Sprintf("{symbol:%v,value:%v,bindingPower:%v}:\n", token.Symbol, token.Value, token.BindingPower))
	for _, child := range token.Children {
		builder.WriteString(child.TreeString(indentLevel + 1))
	}
	return builder.String()
}
