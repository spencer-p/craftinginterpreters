package interpret

import (
	"testing"

	"github.com/spencer-p/craftinginterpreters/pkg/lox/errtrack"
	"github.com/spencer-p/craftinginterpreters/pkg/lox/expr"
	"github.com/spencer-p/craftinginterpreters/pkg/lox/parse"
	"github.com/spencer-p/craftinginterpreters/pkg/lox/scan"

	"github.com/google/go-cmp/cmp"
)

func getExpr(in string) expr.Type {
	fake := errtrack.NewFake()

	toks := scan.New(fake.Tracker, in).Tokens()
	if fake.Tracker.HadError() {
		panic(string(fake.Errors()))
	}

	ast := parse.New(fake.Tracker, toks).AST()
	if fake.Tracker.HadError() {
		panic(string(fake.Errors()))
	}

	return ast
}

func TestInterpret(t *testing.T) {
	table := map[string]struct {
		in      expr.Type
		want    interface{}
		wantErr bool
	}{
		"string":          {in: getExpr("\"hello, world\""), want: "hello, world"},
		"number":          {in: getExpr("42"), want: 42.0},
		"negative number": {in: getExpr("-42"), want: -42.0},
		"not true":        {in: getExpr("!true"), want: false},
		"1+2*3":           {in: getExpr("1+2*3"), want: 7.0},
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
			fake := errtrack.NewFake()
			interpreter := New(fake.Tracker)
			got := interpreter.Eval(tc.in)
			if fake.Tracker.HadError() {
				if tc.wantErr {
					return // successfully failed, move on
				} else {
					t.Errorf("unexpected error: %q", fake.Errors())
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
	got := Stringify(4.0)
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
