package interpret

import (
	"bytes"
	"testing"

	"github.com/spencer-p/craftinginterpreters/pkg/lox/errtrack"
	"github.com/spencer-p/craftinginterpreters/pkg/lox/parse"
	"github.com/spencer-p/craftinginterpreters/pkg/lox/scan"

	"github.com/google/go-cmp/cmp"
)

func TestInterpret(t *testing.T) {
	table := map[string]struct {
		in      string
		want    string
		wanterr bool
	}{
		"string":             {in: "print \"hello, world\";", want: `hello, world`},
		"number":             {in: "print 42;", want: `42`},
		"negative number":    {in: "print -42;", want: `-42`},
		"not true":           {in: "print !true;", want: `false`},
		"1+2*3":              {in: "print 1+2*3;", want: `7`},
		"!true":              {in: "print !true;", want: `false`},
		"greater equal":      {in: "print 1 >= 2;", want: `false`},
		"compose bool":       {in: "print !(1 >= 2);", want: `true`},
		"types not equal":    {in: "print 3 == \"three\";", want: `false`},
		"nil comp":           {in: "print nil == nil;", want: `true`},
		"string comp":        {in: `print "one" == "one";`, want: `true`},
		"-true":              {in: "print -true;", wanterr: true},
		"bad add":            {in: "print 1 + \"two\";", wanterr: true},
		"bad add 2":          {in: "print \"two\" + 1;", wanterr: true},
		"two stmt":           {in: "print 1; print 2;", want: "1\n2"},
		"assignment":         {in: "var x = 2; print x;", want: "2"},
		"lookup fail":        {in: "print x;", wanterr: true},
		"block scope":        {in: "{var x = 1; print x;}", want: "1"},
		"block out of scope": {in: "{var x = 1;} print x;", wanterr: true},
		"shadowing":          {in: "var x = 2; {var x = 1;} print x;", want: "2"},
		"assign to outer":    {in: "var x = 2; {x = 1;} print x;", want: "1"},
		"assign to descoped": {in: "{var x = 1;} x = 2;", wanterr: true},
		"use uninitialized":  {in: "var x; print x;", wanterr: true},
	}

	for name, tc := range table {
		t.Run(name, func(t *testing.T) {
			var fakeOut bytes.Buffer
			fake := errtrack.NewFake()

			toks := scan.New(fake.Tracker, tc.in).Tokens()
			if fake.Tracker.HadError() {
				t.Fatalf(string(fake.Errors()))
			}

			ast := parse.New(fake.Tracker, toks).AST()
			if fake.Tracker.HadError() {
				t.Fatalf(string(fake.Errors()))
			}

			interpreter := New(fake.Tracker)
			interpreter.SetOutput(&fakeOut)
			interpreter.Interpret(ast)
			if fake.Tracker.HadError() {
				if tc.wanterr {
					return // successfully failed, move on
				} else {
					t.Errorf("unexpected error: %q", fake.Errors())
				}
			}

			got := fakeOut.String()
			if diff := cmp.Diff(got, tc.want+"\n"); diff != "" {
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
