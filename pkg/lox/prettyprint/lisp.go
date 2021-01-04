package prettyprint

import (
	"fmt"

	"github.com/spencer-p/craftinginterpreters/pkg/lox/expr"
)

type Lisp struct{}

func (p *Lisp) VisitBinary(e *expr.Binary) interface{} {
	return fmt.Sprintf("(%s %s %s)", e.Op.Lexeme, e.Left.Accept(p).(string), e.Right.Accept(p).(string))
}

func (p *Lisp) VisitGrouping(e *expr.Grouping) interface{} {
	return fmt.Sprintf("(grp %s)", e.Expr.Accept(p).(string))
}

func (p *Lisp) VisitLiteral(e *expr.Literal) interface{} {
	return fmt.Sprintf("%+v", e.Value)
}

func (p *Lisp) VisitUnary(e *expr.Unary) interface{} {
	return fmt.Sprintf("(%s %s)", e.Op.Lexeme, e.Right.Accept(p).(string))
}
