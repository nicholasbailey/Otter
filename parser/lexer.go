package parser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type LexerState int

const (
	unknown       LexerState = 0
	stringLiteral            = 1
	intLiteral               = 2
	floatLiteral             = 3
	name                     = 4
	whiteSpace               = 5
	operator                 = 6
	eof                      = 7
)

func (state LexerState) String() string {
	switch state {
	case unknown:
		return "unknown"
	case stringLiteral:
		return "stringLiteral"
	case intLiteral:
		return "intLiteral"
	case floatLiteral:
		return "floatLiteral"
	case name:
		return "name"
	case whiteSpace:
		return "whiteSpace"
	case operator:
		return "operator"
	case eof:
		return "eof"
	}
	return fmt.Sprint(int(state))
}

func NewLexer(reader io.Reader, symbolTable *LanguageSpecification) *Lexer {
	return &Lexer{
		reader:       bufio.NewReader(reader),
		languageSpec: symbolTable,
		cachedToken:  nil,
		line:         1,
		col:          0,
	}
}

type Lexer struct {
	reader            *bufio.Reader
	languageSpec      *LanguageSpecification
	cachedToken       *Token
	line              int
	col               int
	builder           strings.Builder
	currentState      LexerState
	tokenStartCol     int
	currentQuoteStart rune
}

func (lexer *Lexer) IsBlockStart(token *Token) bool {
	return lexer.languageSpec.IsBlockStart(token.Symbol)
}

func (lexer *Lexer) IsBlockEnd(token *Token, blockStart *Token) bool {
	return lexer.languageSpec.IsBlockEnd(token.Symbol, blockStart.Symbol)
}
func (lexer *Lexer) IsAnyBlockEnd(token *Token) bool {
	return lexer.languageSpec.IsAnyBlockEnd(token.Symbol)
}

func (lexer *Lexer) IsStatementTerminator(token *Token) bool {
	return lexer.languageSpec.IsStatementTerminator(token.Symbol)
}

func (lexer *Lexer) syntaxError(msg string) error {
	return fmt.Errorf("syntaxerror: "+msg+" at line %v, col %v", lexer.line, lexer.col)
}

func (lexer *Lexer) startOfToken(char rune) {
	if quoteSpec := lexer.languageSpec.GetQuoteSpec(char); quoteSpec != nil {
		lexer.currentState = stringLiteral
		lexer.tokenStartCol = lexer.col
		lexer.currentQuoteStart = char
		// Don't write the quote character into the string literal
	} else if lexer.languageSpec.IsIdentifierStartChararacter(char) {
		lexer.currentState = name
		lexer.tokenStartCol = lexer.col
		lexer.builder.WriteRune(char)
	} else if unicode.IsDigit(char) {
		lexer.currentState = intLiteral
		lexer.tokenStartCol = lexer.col
		lexer.builder.WriteRune(char)
	} else if unicode.IsSpace(char) {
		lexer.currentState = whiteSpace
	} else {
		lexer.currentState = operator
		lexer.tokenStartCol = lexer.col
		lexer.builder.WriteRune(char)
	}
}

func (lexer *Lexer) endOfToken() (*Token, error) {
	switch lexer.currentState {
	case stringLiteral:
		quoteSpec := lexer.languageSpec.GetQuoteSpec(lexer.currentQuoteStart)
		if quoteSpec == nil {
			return nil, fmt.Errorf("syntaxerror: invalid quoted literal with quote %v at line %v, col %v", string(lexer.currentQuoteStart), lexer.line, lexer.col)
		}
		stringVal := lexer.builder.String()
		token := lexer.languageSpec.GenerateToken(StringLiteral, stringVal, lexer.line, lexer.tokenStartCol)
		lexer.tokenStartCol = lexer.col
		lexer.builder = strings.Builder{}
		return token, nil
	case intLiteral:
		stringVal := lexer.builder.String()
		token := lexer.languageSpec.GenerateToken(IntLiteral, stringVal, lexer.line, lexer.tokenStartCol)
		lexer.builder = strings.Builder{}
		lexer.tokenStartCol = lexer.col
		return token, nil
	case floatLiteral:
		stringVal := lexer.builder.String()
		token := lexer.languageSpec.GenerateToken(FloatLiteral, stringVal, lexer.line, lexer.tokenStartCol)
		lexer.builder = strings.Builder{}
		lexer.tokenStartCol = lexer.col
		return token, nil
	case name:
		stringVal := lexer.builder.String()
		var token *Token
		if lexer.languageSpec.IsDefined(Symbol(stringVal)) {
			token = lexer.languageSpec.GenerateToken(Symbol(stringVal), stringVal, lexer.line, lexer.tokenStartCol)
		} else {
			token = lexer.languageSpec.GenerateToken(Name, stringVal, lexer.line, lexer.tokenStartCol)
		}
		lexer.builder = strings.Builder{}
		lexer.tokenStartCol = lexer.col
		return token, nil
	case operator:
		stringVal := lexer.builder.String()
		var token *Token
		if lexer.languageSpec.IsDefined(Symbol(stringVal)) {
			token = lexer.languageSpec.GenerateToken(Symbol(stringVal), stringVal, lexer.line, lexer.tokenStartCol)
		} else {
			return nil, fmt.Errorf("syntaxerror: unidentified operator %v at line %v, col %v", stringVal, lexer.line, lexer.col)
		}
		lexer.builder = strings.Builder{}
		lexer.tokenStartCol = lexer.col
		return token, nil
	case whiteSpace:
		return nil, fmt.Errorf("syntaxerror: attempted to resolve token in whitespace at line %v, col %v", lexer.line, lexer.col)
	default:
		return nil, fmt.Errorf("syntaxerror: attempted to resolve token in unkown parse state at line %v, col %v", lexer.line, lexer.col)
	}
}

func (lexer *Lexer) Peek() (*Token, error) {
	token, err := lexer.Next()
	if err != nil {
		return nil, err
	}
	lexer.cachedToken = token
	return token, nil
}

func (lexer *Lexer) readRune() (rune, int, error) {

	char, size, err := lexer.reader.ReadRune()
	if char == '\n' {
		lexer.line++
		lexer.col = 1
	} else {
		lexer.col++
	}
	return char, size, err
}

func (lexer *Lexer) Next() (*Token, error) {
	// First check to see if we have a cached token
	// from a call to Peek
	if lexer.cachedToken != nil {
		token := lexer.cachedToken
		lexer.cachedToken = nil
		return token, nil
	}

	char, size, err := lexer.readRune()
	for size > 0 && err == nil {
		var token *Token
		switch lexer.currentState {
		case unknown:
			if !unicode.IsSpace(char) {
				lexer.startOfToken(char)
			}
		case whiteSpace:
			if !unicode.IsSpace(char) {
				lexer.startOfToken(char)
			}
		case intLiteral:
			if unicode.IsDigit(char) {
				lexer.builder.WriteRune(char)
			} else if char == '.' {
				lexer.currentState = floatLiteral
				lexer.builder.WriteRune(char)
			} else {
				token, err = lexer.endOfToken()
				if err != nil {
					return nil, err
				}
				lexer.startOfToken(char)
			}
		case floatLiteral:
			if unicode.IsDigit(char) {
				lexer.builder.WriteRune(char)
			} else {
				token, err = lexer.endOfToken()
				if err != nil {
					return nil, err
				}
				lexer.startOfToken(char)
			}
		case stringLiteral:
			quoteSpecification := lexer.languageSpec.GetQuoteSpec(lexer.currentQuoteStart)
			if quoteSpecification == nil {
				// This should never happen,
				return nil, fmt.Errorf("syntaxerror: unrecognized quote character '%v'", string(lexer.currentQuoteStart))
			}
			if char == quoteSpecification.closeQuote {
				token, err = lexer.endOfToken()
				if err != nil {
					return nil, err
				}
				lexer.currentState = unknown
			} else if char == '\n' {
				return nil, lexer.syntaxError("new line in middle of string literal")
			} else {
				lexer.builder.WriteRune(char)
			}
		case name:
			if lexer.languageSpec.IsIdentifierCharacter(char) {
				lexer.builder.WriteRune(char)
			} else {
				token, err = lexer.endOfToken()
				if err != nil {
					return nil, err
				}
				lexer.startOfToken(char)
			}
		case operator:
			currentString := lexer.builder.String()
			stringWithNewChar := currentString + string(char)
			if lexer.languageSpec.IsDefined(Symbol(stringWithNewChar)) {
				lexer.builder.WriteRune(char)
			} else if lexer.languageSpec.IsDefined(Symbol(currentString)) {
				token, err = lexer.endOfToken()
				if err != nil {
					return nil, err
				}
				lexer.startOfToken(char)
			} else {
				return nil, fmt.Errorf("syntaxerror: unrecognized operator %v at line %v, col %v", string(char), lexer.line, lexer.col)
			}
		default:
			return nil, fmt.Errorf("syntaxerror: invalid lexer state state %v at line %v, col %v", lexer.currentState, lexer.line, lexer.col)
		}

		if token != nil {
			return token, nil
		}
		char, size, err = lexer.readRune()
	}

	if errors.Is(err, io.EOF) {
		switch lexer.currentState {
		case stringLiteral:
			return nil, fmt.Errorf("syntaxerror: unexpected EOF in string literal at line %v, col %v", lexer.line, lexer.col)
		case eof:
			return lexer.languageSpec.Eof(lexer.line, lexer.col), nil
		default:
			if lexer.builder.Len() > 0 {
				token, err := lexer.endOfToken()
				if err != nil {
					return nil, err
				} else {
					lexer.currentState = eof
					return token, nil
				}
			} else {
				return lexer.languageSpec.Eof(lexer.line, lexer.col), nil
			}
		}
	}
	return nil, fmt.Errorf("syntaxerror: unreadable character %v at line %v, col %v", string(char), lexer.line, lexer.col)
}
