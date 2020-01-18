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
	TokenPlus TokenType = iota
	TokenMinus
	TokenStar
	TokenSlash
	TokenIntLiteral
)

// Token structure
type Token struct {
	token TokenType
	value int
}

var (
	Line    int  = 1
	Putback rune = '\n'
)

const (
	EOF rune = -1
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

// Scan and return the next token found in the input.
// Return 1 if token valid, 0 if no tokens left.
func scan(t *Token) bool {
	// Skip whitespace
	c := skip()
	// Determine the token based on
	// the input character
	switch c {
	case EOF:
		return false
	case '+':
		t.token = TokenPlus
	case '-':
		t.token = TokenMinus
	case '*':
		t.token = TokenStar
	case '/':
		t.token = TokenSlash
	default:
		// If it's a digit, scan the
		// literal integer value in
		if unicode.IsDigit(c) {
			t.value = scanint(c)
			t.token = TokenIntLiteral
			break
		}
		fatal("unrecognized character %c on line %d\n", c, Line)
	}
	// We found a token
	return true
}
