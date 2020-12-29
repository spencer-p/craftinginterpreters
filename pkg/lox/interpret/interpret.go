package interpret

import (
	"fmt"

	"github.com/spencer-p/craftinginterpreters/pkg/lox/expr"
	"github.com/spencer-p/craftinginterpreters/pkg/lox/tok"
)

func Do(e expr.Type) (result interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	result = e.Accept(&Interpreter{})
	return
}

// Interpreter executes code with the visitor pattern.
type Interpreter struct{}

// Verify it satisfies the type
var _ expr.Visitor = &Interpreter{}

func (i *Interpreter) eval(e expr.Type) interface{} {
	return e.Accept(i)
}

func (i *Interpreter) VisitBinary(e *expr.Binary) interface{} {
	left := i.eval(e.Left)
	right := i.eval(e.Right)

	switch e.Op.Typ {
	case tok.MINUS:
		checkNumbers(e.Op, right, left)
		return left.(float64) - right.(float64)
	case tok.SLASH:
		checkNumbers(e.Op, right, left)
		return left.(float64) / right.(float64)
	case tok.STAR:
		checkNumbers(e.Op, right, left)
		return left.(float64) * right.(float64)
	case tok.PLUS:
		if leftActual, ok := left.(float64); ok {
			if rightActual, ok := right.(float64); ok {
				return leftActual + rightActual
			}
			panic(RuntimeError{
				Message: ErrorNotANumber,
				Token:   e.Op,
			})
		} else if leftActual, ok := left.(string); ok {
			if rightActual, ok := right.(string); ok {
				return leftActual + rightActual
			}
			panic(RuntimeError{
				Message: ErrorNotAString,
				Token:   e.Op,
			})
		}
	case tok.GREATER:
		checkNumbers(e.Op, right, left)
		return left.(float64) > right.(float64)
	case tok.GREATER_EQUAL:
		checkNumbers(e.Op, right, left)
		return left.(float64) >= right.(float64)
	case tok.LESS:
		checkNumbers(e.Op, right, left)
		return left.(float64) < right.(float64)
	case tok.LESS_EQUAL:
		checkNumbers(e.Op, right, left)
		return left.(float64) <= right.(float64)
	case tok.BANG_EQUAL:
		return !equal(left, right)
	case tok.EQUAL_EQUAL:
		return equal(left, right)
	default:
		panic(fmt.Errorf("unknown binary operator %#v", e.Op))
	}
	// unreachable
	return nil
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
		checkNumbers(e.Op, right)
		right = -1 * right.(float64)
	case tok.BANG:
		right = !truthy(right)
	}
	return right
}
