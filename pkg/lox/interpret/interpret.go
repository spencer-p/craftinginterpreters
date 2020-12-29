package interpret

import (
	"github.com/spencer-p/craftinginterpreters/pkg/lox/errtrack"
	"github.com/spencer-p/craftinginterpreters/pkg/lox/expr"
	"github.com/spencer-p/craftinginterpreters/pkg/lox/tok"
)

// Interpreter executes code with the visitor pattern.
type Interpreter struct {
	tracker *errtrack.Tracker
}

// Verify it satisfies the type
var _ expr.Visitor = &Interpreter{}

func New(tracker *errtrack.Tracker) *Interpreter {
	return &Interpreter{tracker}
}

func (i *Interpreter) Eval(e expr.Type) interface{} {
	defer i.tracker.CatchFatal()
	return e.Accept(i)
}

func (i *Interpreter) eval(e expr.Type) interface{} {
	return e.Accept(i)
}

func (i *Interpreter) VisitBinary(e *expr.Binary) interface{} {
	left := i.eval(e.Left)
	right := i.eval(e.Right)

	switch e.Op.Typ {
	case tok.MINUS:
		i.checkNumbers(e.Op, right, left)
		return left.(float64) - right.(float64)
	case tok.SLASH:
		i.checkNumbers(e.Op, right, left)
		return left.(float64) / right.(float64)
	case tok.STAR:
		i.checkNumbers(e.Op, right, left)
		return left.(float64) * right.(float64)
	case tok.PLUS:
		if leftActual, ok := left.(float64); ok {
			if rightActual, ok := right.(float64); ok {
				return leftActual + rightActual
			}
			i.tracker.Fatal(errtrack.LoxError{
				Message: ErrorNotANumber,
				Token:   e.Op,
			})
		} else if leftActual, ok := left.(string); ok {
			if rightActual, ok := right.(string); ok {
				return leftActual + rightActual
			}
			i.tracker.Fatal(errtrack.LoxError{
				Message: ErrorNotAString,
				Token:   e.Op,
			})
		}
	case tok.GREATER:
		i.checkNumbers(e.Op, right, left)
		return left.(float64) > right.(float64)
	case tok.GREATER_EQUAL:
		i.checkNumbers(e.Op, right, left)
		return left.(float64) >= right.(float64)
	case tok.LESS:
		i.checkNumbers(e.Op, right, left)
		return left.(float64) < right.(float64)
	case tok.LESS_EQUAL:
		i.checkNumbers(e.Op, right, left)
		return left.(float64) <= right.(float64)
	case tok.BANG_EQUAL:
		return !equal(left, right)
	case tok.EQUAL_EQUAL:
		return equal(left, right)
	default:
		i.tracker.Fatal(errtrack.LoxError{
			Message: ErrorUnknownOp,
			Token:   e.Op,
		})
	}
	i.tracker.Fatal(errtrack.LoxError{
		Message: ErrorNotANumber,
		Token:   e.Op,
	})
	return nil // unreachable
}

func (i *Interpreter) VisitGrouping(e *expr.Grouping) interface{} {
	return i.eval(e.Expression)
}

func (i *Interpreter) VisitLiteral(e *expr.Literal) interface{} {
	return e.Value
}

func (i *Interpreter) VisitUnary(e *expr.Unary) interface{} {
	right := i.eval(e.Right)
	switch e.Op.Typ {
	case tok.MINUS:
		i.checkNumbers(e.Op, right)
		right = -1 * right.(float64)
	case tok.BANG:
		right = !truthy(right)
	}
	return right
}
