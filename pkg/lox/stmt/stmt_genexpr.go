package stmt

import (
	"github.com/spencer-p/craftinginterpreters/pkg/lox/expr"
	"github.com/spencer-p/craftinginterpreters/pkg/lox/tok"
)

type Type interface {
	Accept(Visitor) interface{}
}

type Visitor interface {
	VisitExpression(*Expression) interface{}
	VisitPrint(*Print) interface{}
	VisitVar(*Var) interface{}
}

type Expression struct {
	Expr expr.Type
}

func (e *Expression) Accept(v Visitor) interface{} {
	return v.VisitExpression(e)
}

type Print struct {
	Expr expr.Type
}

func (e *Print) Accept(v Visitor) interface{} {
	return v.VisitPrint(e)
}

type Var struct {
	Name tok.Token
	Initializer expr.Type
}

func (e *Var) Accept(v Visitor) interface{} {
	return v.VisitVar(e)
}

