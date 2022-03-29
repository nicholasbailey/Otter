package parser

func NewOtterLanguage() *LanguageSpecification {

	spec := NewLanguage()

	spec.DefineWhile("while")
	spec.DefineForIn("for", "in")
	spec.DefineIf("if", "else")
	spec.DefineAccess(".")
	spec.DefineComment("//")
	spec.DefineParens("(", ")")
	spec.DefineQuotes('"', '"', StringLiteral)
	spec.DefineQuotes('\'', '\'', StringLiteral)
	spec.DefineReturn("return")
	spec.DefineFunctionDefinition("def")

	spec.DefinePrefix("!", 80)
	spec.DefineInfix("&&", "&&", 30)
	spec.DefineInfix("||", "||", 20)
	spec.DefineInfix("=", Assignment, 10)
	spec.DefineInfix("==", "==", 50)
	spec.DefineInfix("!=", "!=", 50)
	spec.DefineInfix("<", "<", 50)
	spec.DefineInfix(">", ">", 50)
	spec.DefineInfix("<=", "<=", 50)
	spec.DefineInfix(">=", ">=", 50)
	spec.DefineInfix("+", "+", 60)
	spec.DefineInfix("-", "-", 60)
	spec.DefineInfix("*", "*", 70)
	spec.DefineInfix("/", "/", 70)
	spec.DefineInfix("%", "%", 70)
	spec.DefineStatementTerminator(";")
	spec.DefineEmpty(",")
	spec.DefineBlock("{", "}")
	spec.DefineValue("true")
	spec.DefineValue("false")

	return spec
}
