package interpret

import (
	"fmt"
	"io"
	"os"

	"github.com/spencer-p/craftinginterpreters/pkg/lox/errtrack"
	"github.com/spencer-p/craftinginterpreters/pkg/lox/expr"
	"github.com/spencer-p/craftinginterpreters/pkg/lox/stmt"
	"github.com/spencer-p/craftinginterpreters/pkg/lox/tok"
)

// Interpreter executes code with the visitor pattern.
type Interpreter struct {
	tracker *errtrack.Tracker
	out     io.Writer
	env     *Env
}

// Verify it satisfies the visitor types
var _ expr.Visitor = &Interpreter{}
var _ stmt.Visitor = &Interpreter{}

func New(tracker *errtrack.Tracker) *Interpreter {
	return &Interpreter{
		tracker: tracker,
		out:     os.Stdout,
		env:     NewEnv(tracker, nil),
	}
}

func (i *Interpreter) Interpret(stmts []stmt.Type) {
	defer i.tracker.CatchFatal()
	for _, st := range stmts {
		i.execute(st)
	}
}

func (i *Interpreter) SetOutput(w io.Writer) {
	i.out = w
}

func (i *Interpreter) execute(st stmt.Type) {
	st.Accept(i)
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
	return i.eval(e.Expr)
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

func (i *Interpreter) VisitExpression(st *stmt.Expression) interface{} {
	i.eval(st.Expr)
	return nil
}

func (i *Interpreter) VisitPrint(st *stmt.Print) interface{} {
	val := i.eval(st.Expr)
	fmt.Fprintln(i.out, Stringify(val))
	return nil
}

func (i *Interpreter) VisitVariable(e *expr.Variable) interface{} {
	return i.env.Get(e.Name)
}

func (i *Interpreter) VisitVar(st *stmt.Var) interface{} {
	var val interface{}
	if st.Initializer != nil {
		val = i.eval(st.Initializer)
	}

	i.env.Define(st.Name.Lexeme, val)
	return nil
}

func (i *Interpreter) VisitAssign(e *expr.Assign) interface{} {
	val := i.eval(e.Value)
	i.env.Assign(e.Name, val)
	return val
}
