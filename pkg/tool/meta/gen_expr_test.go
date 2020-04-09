package meta

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGenExpr(t *testing.T) {
	in := []Typ{{
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
	}}

	packagename := `dummy`

	want := `package dummy

type Expr interface {
	Visit(Visitor) interface{}
}

type Visitor interface {
	VisitMyExpr(*MyExpr) interface{}
}

type MyExpr struct {
	x int
	a string
	m bool
}

func (e *MyExpr) Visit(v Visitor) interface{} {
	return v.VisitMyExpr(e)
}

`

	var buf bytes.Buffer
	GenExpr(&buf, in, packagename)

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

	got, err := ParseTypes(in)
	if err != nil {
		t.Errorf("unexpected error: %+v", err)
		return
	}
	if diff := cmp.Diff(got, want, cmp.AllowUnexported(Typ{}, Field{})); diff != "" {
		t.Errorf("(-got,+want): %s", diff)
	}
}
