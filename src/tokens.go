package main

import (
	"fmt"
)

//go:generate stringer -type=TokenType
type TokenType int

const (
	// Reserve the zero token as invalid
	INVALID TokenType = iota

	// Single character tokens.
	LEFT_PAREN
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// One or two character tokens.
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literals.
	IDENTIFIER
	STRING
	NUMBER

	// Keywords.
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	// Denote end of file
	EOF
)

type Token struct {
	typ    TokenType
	lexeme string
	lit    interface{}
	line   int
}

func (t Token) String() string {
	return fmt.Sprintf("{%s %q %v}", t.typ.String(), t.lexeme, t.lit)
}
