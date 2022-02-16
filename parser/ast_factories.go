package parser

// Factory functions for AST components.
// Note that these are used by the unsweetener, but generally
// not used by the parser itself which relies on the TDOP mechanics

func BuildWhile(
	conditional *Token,
	block *Token,
	line int,
	col int,
) *Token {

	newTree := &Token{
		Symbol:   While,
		Value:    "while",
		Line:     line,
		Col:      col,
		Children: []*Token{},
	}
	newTree.Children = append(newTree.Children, conditional)
	newTree.Children = append(newTree.Children, block)
	return newTree
}

func BuildBlock(
	statements []*Token,
	line int,
	col int,
) *Token {

	newTree := &Token{
		Symbol:   Block,
		Value:    "{",
		Line:     line,
		Col:      col,
		Children: statements,
	}

	return newTree
}

func BuildAssignment(
	left *Token,
	right *Token,
	line int,
	col int,
) *Token {

	newTree := &Token{
		Symbol:   Assignment,
		Value:    "=",
		Line:     line,
		Col:      col,
		Children: []*Token{left, right},
	}
	return newTree
}

func BuildAccess(
	target *Token,
	methodName string,
	parameters []*Token,
	line int,
	col int,
) *Token {

	functionInvocation := &Token{
		Symbol: FunctionInvocation,
		// TODO - this is probably not quite what we want here,
		// but it matches the behavior of the base parser
		Value:    "(",
		Line:     line,
		Col:      col + 1,
		Children: []*Token{},
	}
	methodNameToken := &Token{
		Symbol: Name,
		Value:  methodName,
		Line:   line,
		Col:    col + 1,
	}

	functionInvocation.Children = append(functionInvocation.Children, methodNameToken)
	functionInvocation.Children = append(functionInvocation.Children, parameters...)
	newTree := &Token{
		Symbol: Access,
		// TODO, this is also just a compatability thing
		// Got to decide if value even has meaning post parse
		// Probably makes sense to do some sort of AST/Token separation
		// at some point
		Value:    ".",
		Line:     line,
		Col:      col,
		Children: []*Token{target, functionInvocation},
	}
	return newTree
}

func BuildName(name string, line int, col int) *Token {
	return &Token{
		Symbol: Name,
		Value:  name,
		Line:   line,
		Col:    col,
	}
}
