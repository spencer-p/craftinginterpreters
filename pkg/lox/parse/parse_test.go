package parse

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/spencer-p/craftinginterpreters/pkg/lox/errtrack"
	"github.com/spencer-p/craftinginterpreters/pkg/lox/expr"
	"github.com/spencer-p/craftinginterpreters/pkg/lox/scan"
	"github.com/spencer-p/craftinginterpreters/pkg/lox/stmt"
	. "github.com/spencer-p/craftinginterpreters/pkg/lox/tok"
)

func TestParse(t *testing.T) {
	table := []struct {
		in       string
		want     []stmt.Type
		wantExpr expr.Type
		wanterr  bool
	}{{
		in:       `1`,
		wantExpr: &expr.Literal{1.0},
	}, {
		in: `1 == 2`,
		wantExpr: &expr.Binary{
			Left:  &expr.Literal{1.0},
			Right: &expr.Literal{2.0},
			Op: Token{
				Typ: EQUAL_EQUAL,
			},
		},
	}, {
		in: `2 != 1`,
		wantExpr: &expr.Binary{
			Left:  &expr.Literal{2.0},
			Right: &expr.Literal{1.0},
			Op: Token{
				Typ: BANG_EQUAL,
			},
		},
	}, {
		in: `2 != 1 == true`,
		wantExpr: &expr.Binary{
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
		wantExpr: &expr.Binary{
			Left:  &expr.Literal{1.0},
			Right: &expr.Literal{2.0},
			Op:    Token{Typ: LESS},
		},
	}, {
		in: `1 + 2`,
		wantExpr: &expr.Binary{
			Left:  &expr.Literal{1.0},
			Right: &expr.Literal{2.0},
			Op:    Token{Typ: PLUS},
		},
	}, {
		in: `1 + 2 * 3`,
		wantExpr: &expr.Binary{
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
		wantExpr: &expr.Unary{
			Op:    Token{Typ: MINUS},
			Right: &expr.Literal{12.0},
		},
	}, {
		in: `!false`,
		wantExpr: &expr.Unary{
			Op:    Token{Typ: BANG},
			Right: &expr.Literal{false},
		},
	}, {
		in:       `"hello world!"`,
		wantExpr: &expr.Literal{"hello world!"},
	}, {
		in:       `nil`,
		wantExpr: &expr.Literal{nil},
	}, {
		in: `(1 + 2)`,
		wantExpr: &expr.Grouping{&expr.Binary{
			Left:  &expr.Literal{1.0},
			Right: &expr.Literal{2.0},
			Op:    Token{Typ: PLUS},
		}},
	}, {
		in:      `(1 + 2`,
		wanterr: true,
	}, {
		in: `print "hello world";`,
		want: []stmt.Type{&stmt.Print{
			Expr: &expr.Literal{"hello world"},
		}},
	}, {
		in:      `print "hello world"`,
		wanterr: true,
	}, {
		in: `1;
		print "hello world";`,
		want: []stmt.Type{
			&stmt.Expression{Expr: &expr.Literal{1.0}},
			&stmt.Print{Expr: &expr.Literal{"hello world"}},
		},
	}, {
		in:      `1 = 2;`,
		wanterr: true,
	}, {
		in: `myVar = 2;`,
		want: []stmt.Type{&stmt.Expression{
			&expr.Assign{Token{}, &expr.Literal{2.0}},
		}},
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
			// quick hack to box expression tests into statements
			if row.wantExpr != nil {
				row.want = []stmt.Type{&stmt.Expression{row.wantExpr}}
				row.in = row.in + ";"
			}

			fake := errtrack.NewFake()
			tokens := scan.New(fake.Tracker, row.in).Tokens() // not too happy about dependency. writing tokens is hard.
			got := New(fake.Tracker, tokens).AST()

			if row.wanterr {
				if fake.Tracker.HadError() == false {
					t.Errorf("Wanted an error but got none")
				}
				return // successfully errored
			}

			if fake.Tracker.HadError() && row.wanterr == false {
				t.Errorf("Parse %q unexpected error %q", row.in, fake.Errors())
			} else if diff := cmp.Diff(got, row.want, ignoreTokenTypeFields); diff != "" {
				t.Errorf("Parse %q failed (-got, +want): %s", row.in, diff)
			}
		})
	}
}
