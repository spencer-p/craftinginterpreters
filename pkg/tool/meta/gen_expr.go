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

	exprInterface = "type Type interface {\n\tAccept(Visitor) interface{}\n}\n\n"

	visitorMethod = "func (e *%s) Accept(v Visitor) interface{} {\n\treturn v.Visit%s(e)\n}\n\n"

	importPrefix = "import "
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

type Info struct {
	Package string
	Types   []Typ
	Imports []string
}

func GenExpr(out io.Writer, i *Info) {
	writeHeader(out, i.Package)
	writeImports(out, i.Imports)
	writeExprInterface(out)
	writeVisitorInterface(out, i.Types)
	for _, typ := range i.Types {
		writeType(out, typ)
		writeVisitorMethod(out, typ)
	}
}

func writeHeader(out io.Writer, packagename string) {
	out.Write([]byte(`package `))
	out.Write([]byte(packagename))
	out.Write([]byte{'\n', '\n'})
}

func writeImports(out io.Writer, imports []string) {
	fmt.Fprintf(out, "import (\n")
	for i := range imports {
		fmt.Fprintf(out, "\t\"%s\"\n", imports[i])
	}
	fmt.Fprintf(out, ")\n\n")
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

func ParseTypes(info *Info, in string) error {
	// Sorry about this.

	if info.Types == nil {
		info.Types = make([]Typ, 0)
	}

	lines := strings.Split(in, "\n")
	for i := range lines {
		if strings.HasPrefix(lines[i], importPrefix) {
			info.Imports = append(info.Imports, strings.Trim(lines[i][len(importPrefix):], " "))
			continue
		}

		nameAndFields := strings.SplitN(strings.Trim(lines[i], " \t"), ":", 2)

		if len(nameAndFields) != 2 {
			return ParseError
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
				return ParseError
			}
			t.fields = append(t.fields, Field{
				name: strings.Trim(nameAndTyp[0], " "),
				typ:  strings.Trim(nameAndTyp[1], " "),
			})
		}

		info.Types = append(info.Types, t)
	}
	return nil
}
