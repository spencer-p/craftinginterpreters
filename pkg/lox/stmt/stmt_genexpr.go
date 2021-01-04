package stmt

import (
	"github.com/spencer-p/craftinginterpreters/pkg/lox/expr"
)

type Type interface {
	Accept(Visitor) interface{}
}

type Visitor interface {
	VisitExpression(*Expression) interface{}
	VisitPrint(*Print) interface{}
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

