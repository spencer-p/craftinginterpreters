package prettyprint

import (
	"github.com/spencer-p/craftinginterpreters/pkg/lox/expr"

	"github.com/davecgh/go-spew/spew"
)

type Go struct{}

func (p *Go) VisitBinary(e *expr.Binary) interface{} {
	return spew.Sprintf("%#v", e)
}

func (p *Go) VisitGrouping(e *expr.Grouping) interface{} {
	return spew.Sprintf("%#v", e)
}

func (p *Go) VisitLiteral(e *expr.Literal) interface{} {
	return spew.Sprintf("%#v", e)
}

func (p *Go) VisitUnary(e *expr.Unary) interface{} {
	return spew.Sprintf("%#v", e)
}
