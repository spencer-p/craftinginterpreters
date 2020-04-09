package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestScanner(t *testing.T) {
	table := []struct {
		in   string
		want []TokenType
	}{{
		in:   "(( )){}",
		want: []TokenType{LEFT_PAREN, LEFT_PAREN, RIGHT_PAREN, RIGHT_PAREN, LEFT_BRACE, RIGHT_BRACE, EOF},
	}, {
		in:   "!*+-/=<>",
		want: []TokenType{BANG, STAR, PLUS, MINUS, SLASH, EQUAL, LESS, GREATER, EOF},
	}, {
		in:   "// this is a comment +=<=(){}",
		want: []TokenType{EOF},
	}, {
		in:   "<= >= !===",
		want: []TokenType{LESS_EQUAL, GREATER_EQUAL, BANG_EQUAL, EQUAL_EQUAL, EOF},
	}}

	for _, test := range table {
		tokens := NewScanner(test.in).Tokens()
		got := make([]TokenType, len(tokens))
		for i := 0; i < len(tokens); i++ {
			got[i] = tokens[i].typ
		}
		if diff := cmp.Diff(got, test.want); diff != "" {
			t.Errorf("Bad scan of %q (-got,+want): %s", test.in, diff)
		}
	}
}
