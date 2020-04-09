package meta

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	interfacePreamble = "type Visitor interface {\n"
	interfaceFunc     = "	Visit%s(*%s) interface{}\n"
	interfaceSuffix   = "}\n\n"

	exprInterface = "type Expr interface {\n\tVisit(Visitor) interface{}\n}\n\n"

	visitorMethod = "func (e *%s) Visit(v Visitor) interface{} {\n\treturn v.Visit%s(e)\n}\n\n"
)

var (
	ParseError = errors.New("parse error")
)

type Field struct {
	name string
	typ  string
}

type Typ struct {
	name   string
	fields []Field
}

func GenExpr(out io.Writer, types []Typ, packagename string) {
	writeHeader(out, packagename)
	writeExprInterface(out)
	writeVisitorInterface(out, types)
	for _, typ := range types {
		writeType(out, typ)
		writeVisitorMethod(out, typ)
	}
}

func writeHeader(out io.Writer, packagename string) {
	out.Write([]byte(`package `))
	out.Write([]byte(packagename))
	out.Write([]byte{'\n', '\n'})
}

func writeExprInterface(out io.Writer) {
	out.Write([]byte(exprInterface))
}

func writeVisitorInterface(out io.Writer, types []Typ) {
	out.Write([]byte(interfacePreamble))
	for _, t := range types {
		fmt.Fprintf(out, interfaceFunc, t.name, t.name)
	}
	out.Write([]byte(interfaceSuffix))
}

func writeType(out io.Writer, typ Typ) {
	fmt.Fprintf(out, "type %s struct {\n", typ.name)
	for _, f := range typ.fields {
		fmt.Fprintf(out, "\t%s %s\n", f.name, f.typ)
	}
	fmt.Fprintf(out, "}\n\n")
}

func writeVisitorMethod(out io.Writer, typ Typ) {
	fmt.Fprintf(out, visitorMethod, typ.name, typ.name)
}

func ParseTypes(in string) ([]Typ, error) {
	// Sorry about this.

	types := make([]Typ, 0)
	lines := strings.Split(in, "\n")
	for i := range lines {
		nameAndFields := strings.SplitN(strings.Trim(lines[i], " \t"), ":", 2)

		if len(nameAndFields) != 2 {
			return nil, ParseError
		}

		name := strings.Trim(nameAndFields[0], " ")
		fields := strings.Trim(nameAndFields[1], " ")

		t := Typ{
			name:   name,
			fields: make([]Field, 0),
		}

		for _, field := range strings.Split(fields, ",") {
			nameAndTyp := strings.SplitN(strings.Trim(field, " "), " ", 2)
			if len(nameAndTyp) != 2 {
				return nil, ParseError
			}
			t.fields = append(t.fields, Field{
				name: strings.Trim(nameAndTyp[0], " "),
				typ:  strings.Trim(nameAndTyp[1], " "),
			})
		}

		types = append(types, t)
	}
	return types, nil
}
