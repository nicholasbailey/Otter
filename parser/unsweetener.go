package parser

import "github.com/nicholasbailey/otter/exception"

// An Unsweeter takes a AST and removes the 'Syntactic Sugar' by transpiling
// rich ASTs to a simpler set of operations
type Unsweetener interface {
	Unsweeten(tree *Token) (*Token, exception.Exception)
}

func NewUnsweetener() Unsweetener {
	unsweeteningRules := map[Symbol]UnsweetingRule{}
	unsweeteningRules[ForIn] = UnsweetenForIn
	return &SimpleUnsweeter{
		UnsweeteningRules: unsweeteningRules,
	}
}

type UnsweetingRule func(tree *Token) (*Token, exception.Exception)

type SimpleUnsweeter struct {
	UnsweeteningRules map[Symbol]UnsweetingRule
}

func (unsweetener *SimpleUnsweeter) Unsweeten(tree *Token) (*Token, exception.Exception) {
	rule, found := unsweetener.UnsweeteningRules[tree.Symbol]
	if found {
		return rule(tree)
	}
	return tree, nil
}

// Converts the syntax tree for a for-in loop to a while
// loop
func UnsweetenForIn(tree *Token) (*Token, exception.Exception) {

	iterationVariableToken := tree.Children[0]
	iterableToken := tree.Children[1]
	originalBlockToken := tree.Children[2]

	iterationVariableName := iterationVariableToken.Value
	iteratorVariableName := "~" + iterationVariableName + "Iterator"

	iterationVariableInitalizer := BuildAssignment(
		BuildName(iterationVariableName, iterationVariableToken.Line, iterationVariableToken.Col),
		BuildName("null", iterationVariableToken.Line, iterationVariableToken.Col),
		iterationVariableToken.Line,
		iterationVariableToken.Col,
	)
	iteratorInitializer := BuildAssignment(
		BuildName(iteratorVariableName, iterableToken.Line, iterableToken.Col),
		BuildAccess(iterableToken, "iterator", []*Token{}, iterableToken.Line, iterableToken.Col),
		iterableToken.Line,
		iterableToken.Col,
	)

	whileConditionalExpression := BuildAccess(
		BuildName(iteratorVariableName, iterableToken.Line, iterableToken.Col),
		"hasNext",
		[]*Token{},
		iterableToken.Line,
		iterableToken.Col,
	)

	newBlockChildren := []*Token{
		BuildAssignment(
			BuildName(iterationVariableName, iterationVariableToken.Line, iterationVariableToken.Col),
			BuildAccess(
				BuildName(iteratorVariableName, iterableToken.Line, iterableToken.Col),
				"getNext",
				[]*Token{},
				iterableToken.Line,
				iterableToken.Col,
			),
			iterableToken.Line,
			iterableToken.Col,
		),
	}
	newBlockChildren = append(newBlockChildren, originalBlockToken.Children...)

	newBlock := BuildBlock(newBlockChildren, originalBlockToken.Line, originalBlockToken.Col)

	whileStatement := BuildWhile(
		whileConditionalExpression,
		newBlock,
		tree.Line,
		tree.Col,
	)
	unsweetenedTree := BuildBlock([]*Token{iterationVariableInitalizer, iteratorInitializer, whileStatement}, tree.Line, tree.Col)
	return unsweetenedTree, nil
}
