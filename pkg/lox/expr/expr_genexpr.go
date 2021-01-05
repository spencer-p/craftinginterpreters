package expr

import (
	"github.com/spencer-p/craftinginterpreters/pkg/lox/tok"
)

type Type interface {
	Accept(Visitor) interface{}
}

type Visitor interface {
	VisitBinary(*Binary) interface{}
	VisitGrouping(*Grouping) interface{}
	VisitLiteral(*Literal) interface{}
	VisitUnary(*Unary) interface{}
	VisitVariable(*Variable) interface{}
	VisitAssign(*Assign) interface{}
}

type Binary struct {
	Left Type
	Right Type
	Op tok.Token
}

func (e *Binary) Accept(v Visitor) interface{} {
	return v.VisitBinary(e)
}

type Grouping struct {
	Expr Type
}

func (e *Grouping) Accept(v Visitor) interface{} {
	return v.VisitGrouping(e)
}

type Literal struct {
	Value interface{}
}

func (e *Literal) Accept(v Visitor) interface{} {
	return v.VisitLiteral(e)
}

type Unary struct {
	Op tok.Token
	Right Type
}

func (e *Unary) Accept(v Visitor) interface{} {
	return v.VisitUnary(e)
}

type Variable struct {
	Name tok.Token
}

func (e *Variable) Accept(v Visitor) interface{} {
	return v.VisitVariable(e)
}

type Assign struct {
	Name tok.Token
	Value Type
}

func (e *Assign) Accept(v Visitor) interface{} {
	return v.VisitAssign(e)
}

