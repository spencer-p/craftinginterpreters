package meta

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGenExpr(t *testing.T) {
	in := Info{
		Package: "dummy",
		Imports: []string{"fmt"},
		Types: []Typ{{
			name: "MyExpr",
			fields: []Field{{
				name: "x",
				typ:  "int",
			}, {
				name: "a",
				typ:  "string",
			}, {
				name: "m",
				typ:  "bool",
			}},
		}}}

	want := `package dummy

import (
	"fmt"
)

type Expr interface {
	Accept(Visitor) interface{}
}

type Visitor interface {
	VisitMyExpr(*MyExpr) interface{}
}

type MyExpr struct {
	x int
	a string
	m bool
}

func (e *MyExpr) Accept(v Visitor) interface{} {
	return v.VisitMyExpr(e)
}

`

	var buf bytes.Buffer
	GenExpr(&buf, &in)

	if diff := cmp.Diff(buf.Bytes(), []byte(want)); diff != "" {
		t.Errorf("error in genexpr (-got,+want): %s", diff)
	}
}

func TestParseTypes(t *testing.T) {
	in := `MyExpr: x int, a string, m bool
	MyOther: a int`

	want := []Typ{{
		name: "MyExpr",
		fields: []Field{{
			name: "x",
			typ:  "int",
		}, {
			name: "a",
			typ:  "string",
		}, {
			name: "m",
			typ:  "bool",
		}},
	}, {
		name: "MyOther",
		fields: []Field{{
			name: "a",
			typ:  "int",
		}},
	}}

	var got Info
	err := ParseTypes(&got, in)
	if err != nil {
		t.Errorf("unexpected error: %+v", err)
		return
	}
	if diff := cmp.Diff(got.Types, want, cmp.AllowUnexported(Typ{}, Field{})); diff != "" {
		t.Errorf("(-got,+want): %s", diff)
	}
}
