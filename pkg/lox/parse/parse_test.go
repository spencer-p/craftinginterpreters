package parse

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/spencer-p/craftinginterpreters/pkg/lox/expr"
	"github.com/spencer-p/craftinginterpreters/pkg/lox/scan"
	. "github.com/spencer-p/craftinginterpreters/pkg/lox/tok"
)

func TestParse(t *testing.T) {
	table := []struct {
		in   string
		want expr.Type
	}{{
		in:   "1",
		want: &expr.Literal{1},
	}, {
		in: "1 == 2",
		want: &expr.Binary{
			Left:  &expr.Literal{1},
			Right: &expr.Literal{2},
			Op: Token{
				Typ: EQUAL_EQUAL,
			},
		},
	}}

	for _, row := range table {
		tokens := scan.New(row.in).Tokens() // not too happy about dependency. writing tokens is hard.
		got := NewParser(tokens).AST()
		if diff := cmp.Diff(got, row.want); diff != "" {
			t.Errorf("Parse %q failed (-got, +want): %s", row.in, diff)
		}
	}
}
