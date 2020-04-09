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
		in:   `(( )){}`,
		want: []TokenType{LEFT_PAREN, LEFT_PAREN, RIGHT_PAREN, RIGHT_PAREN, LEFT_BRACE, RIGHT_BRACE, EOF},
	}, {
		in:   `!*+-/=<>`,
		want: []TokenType{BANG, STAR, PLUS, MINUS, SLASH, EQUAL, LESS, GREATER, EOF},
	}, {
		in:   `// this is a comment +=<=(){}`,
		want: []TokenType{EOF},
	}, {
		in:   `<= >= !===`,
		want: []TokenType{LESS_EQUAL, GREATER_EQUAL, BANG_EQUAL, EQUAL_EQUAL, EOF},
	}, {
		in:   `"你好, world!"`, // unicode support!
		want: []TokenType{STRING, EOF},
	}, {
		in: `"newlines can be
		in a string"`,
		want: []TokenType{STRING, EOF},
	}, {
		in:   `123.4`,
		want: []TokenType{NUMBER, EOF},
	}, {
		in:   `123`,
		want: []TokenType{NUMBER, EOF},
	}, {
		in:   `0.123`,
		want: []TokenType{NUMBER, EOF},
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
