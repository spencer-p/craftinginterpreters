package parse

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/spencer-p/craftinginterpreters/pkg/lox/expr"
	"github.com/spencer-p/craftinginterpreters/pkg/lox/scan"
	. "github.com/spencer-p/craftinginterpreters/pkg/lox/tok"
)

func TestParse(t *testing.T) {
	table := []struct {
		in      string
		want    expr.Type
		wanterr bool
	}{{
		in:   `1`,
		want: &expr.Literal{1.0},
	}, {
		in: `1 == 2`,
		want: &expr.Binary{
			Left:  &expr.Literal{1.0},
			Right: &expr.Literal{2.0},
			Op: Token{
				Typ: EQUAL_EQUAL,
			},
		},
	}, {
		in: `2 != 1`,
		want: &expr.Binary{
			Left:  &expr.Literal{2.0},
			Right: &expr.Literal{1.0},
			Op: Token{
				Typ: BANG_EQUAL,
			},
		},
	}, {
		in: `2 != 1 == true`,
		want: &expr.Binary{
			Left: &expr.Binary{
				Left:  &expr.Literal{2.0},
				Right: &expr.Literal{1.0},
				Op:    Token{Typ: BANG_EQUAL},
			},
			Right: &expr.Literal{true},
			Op:    Token{Typ: EQUAL_EQUAL},
		},
	}, {
		in: `1 < 2`,
		want: &expr.Binary{
			Left:  &expr.Literal{1.0},
			Right: &expr.Literal{2.0},
			Op:    Token{Typ: LESS},
		},
	}, {
		in: `1 + 2`,
		want: &expr.Binary{
			Left:  &expr.Literal{1.0},
			Right: &expr.Literal{2.0},
			Op:    Token{Typ: PLUS},
		},
	}, {
		in: `1 + 2 * 3`,
		want: &expr.Binary{
			Left: &expr.Literal{1.0},
			Right: &expr.Binary{
				Left:  &expr.Literal{2.0},
				Right: &expr.Literal{3.0},
				Op:    Token{Typ: STAR},
			},
			Op: Token{Typ: PLUS},
		},
	}, {
		in: `-12`,
		want: &expr.Unary{
			Op:    Token{Typ: MINUS},
			Right: &expr.Literal{12.0},
		},
	}, {
		in: `!false`,
		want: &expr.Unary{
			Op:    Token{Typ: BANG},
			Right: &expr.Literal{false},
		},
	}, {
		in:   `"hello world!"`,
		want: &expr.Literal{"hello world!"},
	}, {
		in:   `nil`,
		want: &expr.Literal{nil},
	}, {
		in: `(1 + 2)`,
		want: &expr.Grouping{&expr.Binary{
			Left:  &expr.Literal{1.0},
			Right: &expr.Literal{2.0},
			Op:    Token{Typ: PLUS},
		}},
	}, {
		in:      `(1 + 2`,
		wanterr: true,
	}}

	ignoreTokenTypeFields := cmp.FilterPath(func(path cmp.Path) bool {
		possibleTokenType := path.Index(-2) // parent of value getting compared
		if possibleTokenType.Type() == reflect.TypeOf(Token{}) &&
			path.Last().String() != "Typ" {
			// parent type is Token and the member type is not Token.Typ
			return true
		} else {
			// Token.Typ is allowed along with anything not in a Token
			return false
		}
	}, cmp.Ignore())

	for _, row := range table {
		t.Run(row.in, func(t *testing.T) {
			tokens, _ := scan.New(row.in).Tokens() // not too happy about dependency. writing tokens is hard.
			got, err := New(tokens).AST()
			if err != nil && row.wanterr == false {
				t.Errorf("Parse %q unexpected error %v", row.in, err)
			} else if diff := cmp.Diff(got, row.want, ignoreTokenTypeFields); diff != "" {
				t.Errorf("Parse %q failed (-got, +want): %s", row.in, diff)
			}
		})
	}
}
