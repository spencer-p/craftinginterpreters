package interpret

import (
	"fmt"
	"testing"

	"github.com/spencer-p/craftinginterpreters/pkg/lox/expr"
	"github.com/spencer-p/craftinginterpreters/pkg/lox/parse"
	"github.com/spencer-p/craftinginterpreters/pkg/lox/scan"
	"github.com/spencer-p/craftinginterpreters/pkg/lox/tok"

	"github.com/google/go-cmp/cmp"
)

func getExpr(in string) expr.Type {
	toks, err := scan.New(in).Tokens()
	if err != nil {
		panic(fmt.Errorf("%w: %q", err, in))
	}

	ast, err := parse.New(toks).AST()
	if err != nil {
		panic(fmt.Errorf("%w: %q", err, in))
	}

	return ast
}

func TestInterpret(t *testing.T) {
	table := map[string]struct {
		in      expr.Type
		want    interface{}
		wantErr bool
	}{
		"string":          {in: &expr.Literal{Value: "hello, world"}, want: "hello, world"},
		"number":          {in: &expr.Literal{Value: 42.0}, want: 42.0},
		"negative number": {in: &expr.Unary{Op: tok.Token{Typ: tok.MINUS}, Right: &expr.Literal{Value: 42.0}}, want: -42.0},
		"not true":        {in: &expr.Unary{Op: tok.Token{Typ: tok.BANG}, Right: &expr.Literal{Value: true}}, want: false},
		"1+2*3":           {in: &expr.Binary{Left: &expr.Literal{Value: 1.0}, Right: &expr.Binary{Left: &expr.Literal{Value: 2.0}, Right: &expr.Literal{Value: 3.0}, Op: tok.Token{Typ: tok.STAR}}, Op: tok.Token{Typ: tok.PLUS}}, want: 7.0},
		"!true":           {in: getExpr("!true"), want: false},
		"greater equal":   {in: getExpr("1 >= 2"), want: false},
		"compose bool":    {in: getExpr("!(1 >= 2)"), want: true},
		"types not equal": {in: getExpr("3 == \"three\""), want: false},
		"nil comp":        {in: getExpr("nil == nil"), want: true},
		"string comp":     {in: getExpr(`"one" == "one"`), want: true},
		"-true":           {in: getExpr("-true"), wantErr: true},
		"bad add":         {in: getExpr("1 + \"two\""), wantErr: true},
		"bad add 2":       {in: getExpr("\"two\" + 1"), wantErr: true},
	}

	for name, tc := range table {
		t.Run(name, func(t *testing.T) {
			got, err := Do(tc.in)
			if err != nil {
				if tc.wantErr {
					return // successfully failed, move on
				} else {
					t.Errorf("unexpected error: %v", err)
				}
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("incorrect interpretation (-got,+want): %s", diff)
			}
		})
	}
}

func TestStringify(t *testing.T) {
	want := "4"
	got := Stringify(4.0, nil)
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
