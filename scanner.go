package main

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

// TokenType
type TokenType int

// Tokens
const (
	TokenEOF                TokenType = iota
	TokenPlus                         // +
	TokenMinus                        // -
	TokenStar                         // *
	TokenSlash                        // /
	TokenIntLiteral                   // 0
	TokenSemicolon                    // ;
	TokenAssign                       // =
	TokenEqual                        // ==
	TokenNotEqual                     // !=
	TokenLessThan                     // <
	TokenLessThanOrEqual              // <=
	TokenGreaterThan                  // >
	TokenGreaterThanOrEqual           // >=

	TokenAmpersand // &
	TokenAnd       // &&
	TokenComma     // ,

	TokenLeftBrace  // {
	TokenRightBrace // }
	TokenLeftParen  // (
	TokenRightParen // )

	TokenIf    // if
	TokenElse  // else
	TokenWhile // while
	TokenFor   // for
	TokenVoid  // void

	TokenIdent // x

	TokenPrint  // print
	TokenInt    // int
	TokenChar   // char
	TokenLong   // long
	TokenReturn // return
)

// Token structure
type Token struct {
	token TokenType
	value int
}

var (
	Line       int  = 1
	Putback    rune = '\n'
	Text            = ""
	FunctionId int  = -1
)

var Keywords = map[string]TokenType{
	"print":  TokenPrint,
	"int":    TokenInt,
	"if":     TokenIf,
	"else":   TokenElse,
	"while":  TokenWhile,
	"for":    TokenFor,
	"void":   TokenVoid,
	"char":   TokenChar,
	"long":   TokenLong,
	"return": TokenReturn,
}

const (
	EOF            rune = -1
	MaxIdentLength int  = 512
)

// Get the next character from the input file.
func next() rune {
	c := ' '
	if Putback != 0 { // Use the character put
		c = Putback // back if there is one
		Putback = 0
		return c
	}
	c, _, err := InFile.ReadRune()
	if err == io.EOF {
		return EOF
	}
	if err != nil {
		fatal("failed to read from file: %v\n", err)
	}
	if c == '\n' {
		Line++
	}
	return c
}

// Put back an unwanted character
func putback(c rune) {
	Putback = c
}

// Skip past input that we don't need to deal with,
// i.e. whitespace, newlines. Return the first
// character we do need to deal with.
func skip() rune {
	c := next()
	for c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f' {
		c = next()
	}
	return c
}

// Scan and return an integer literal
// value from the input file. Store
// the value as a string in Text.
func scanint(c rune) int {
	val := 0
	for {
		k := strings.Index("0123456789", fmt.Sprintf("%c", c))
		if k < 0 {
			break
		}
		val = val*10 + k
		c = next()
	}
	// We hit a non-integer character, put it back.
	putback(c)
	return val
}

// Scan an identifier from the input file and
// store it in buf[]. Return the identifier's length
func scanident(c rune, lim int) string {
	i := 0
	buf := make([]rune, 0)
	// Allow digits, alpha and underscores
	for unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_' {
		// Error if we hit the identifier length limit,
		// else append to buf[] and get next character
		if lim-1 == i {
			fatal("identifier too long on line %d\n", Line)
		} else if i < lim-1 {
			buf = append(buf, c)
			i++
		}
		c = next()
	}
	// We hit a non-valid character, put it back.
	// NUL-terminate the buf[] and return the length
	putback(c)
	return string(buf)
}

// Ensure that the current token is t,
// and fetch the next token. Otherwise
// throw an error
func match(t TokenType, s string) {
	if t == CurrentToken.token {
		scan(CurrentToken)
	} else {
		fatal("%s expected on line %d\n", s, Line)
	}
}

// Match a semicon and fetch the next token
func semi() {
	match(TokenSemicolon, ";")
}

func ident() {
	match(TokenIdent, "identifier")
}

func lbrace() {
	match(TokenLeftBrace, "{")
}

func rbrace() {
	match(TokenRightBrace, "}")
}

func lparen() {
	match(TokenLeftParen, "(")
}

func rparen() {
	match(TokenRightParen, ")")
}

// Reject the token that we just scanned
func rejectToken(t *Token) {
	if RejectedToken != nil {
		fatal("can't reject token twice\n")
	}
	RejectedToken = t
}

// Scan and return the next token found in the input.
// Return 1 if token valid, 0 if no tokens left.
func scan(t *Token) bool {
	// If we have any rejected token, return it
	if RejectedToken != nil {
		t = RejectedToken
		RejectedToken = nil
		return true
	}
	// Skip whitespace
	c := skip()
	// Determine the token based on the input character
	switch c {
	case EOF:
		t.token = TokenEOF
		return false
	case '+':
		t.token = TokenPlus
	case '-':
		t.token = TokenMinus
	case '*':
		t.token = TokenStar
	case '/':
		t.token = TokenSlash
	case ';':
		t.token = TokenSemicolon
	case '{':
		t.token = TokenLeftBrace
	case '}':
		t.token = TokenRightBrace
	case '(':
		t.token = TokenLeftParen
	case ')':
		t.token = TokenRightParen
	case ',':
		t.token = TokenComma
	case '=':
		c = next()
		if c == '=' {
			t.token = TokenEqual
		} else {
			putback(c)
			t.token = TokenAssign
		}
	case '!':
		c = next()
		if c == '=' {
			t.token = TokenNotEqual
		} else {
			fatal("unrecognized character %c\n", c)
		}
	case '<':
		c = next()
		if c == '=' {
			t.token = TokenLessThanOrEqual
		} else {
			putback(c)
			t.token = TokenLessThan
		}
		break
	case '>':
		c = next()
		if c == '=' {
			t.token = TokenGreaterThanOrEqual
		} else {
			putback(c)
			t.token = TokenGreaterThan
		}
	case '&':
		c = next()
		if c == '&' {
			t.token = TokenAnd
		} else {
			putback(c)
			t.token = TokenAmpersand
		}
	default:
		if unicode.IsDigit(c) {
			t.value = scanint(c)
			t.token = TokenIntLiteral
			break
		} else if unicode.IsLetter(c) || c == '_' {
			Text = scanident(c, MaxIdentLength)
			tokenType, exists := Keywords[Text]
			if exists {
				t.token = tokenType
			} else {
				t.token = TokenIdent
			}
		} else {
			fatal("unrecognized character %c on line %d\n", c, Line)
		}
	}
	return true
}
