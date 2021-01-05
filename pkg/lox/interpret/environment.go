package interpret

import (
	"errors"

	"github.com/spencer-p/craftinginterpreters/pkg/lox/errtrack"
	"github.com/spencer-p/craftinginterpreters/pkg/lox/tok"
)

type Env struct {
	table     map[string]interface{}
	tracker   *errtrack.Tracker
	enclosing *Env
}

type Uninitialized struct{}

// NewEnv constructs an empty environment with enclosing parent.
// Enclosing may be nil.
func NewEnv(tracker *errtrack.Tracker, enclosing *Env) *Env {
	// If there is no enclosing environment, then we construct a dummy
	// environment that can report errors when it is used.
	if enclosing == nil {
		enclosing = &Env{
			tracker: tracker,
		}
	}
	return &Env{
		table:     make(map[string]interface{}),
		tracker:   tracker,
		enclosing: enclosing,
	}
}

func (e *Env) Define(name string, val interface{}) {
	e.table[name] = val
}

func (e *Env) Get(name tok.Token) interface{} {
	if e.table == nil {
		e.tracker.Fatal(errtrack.ErrorUndefined(name))
	}

	val, ok := e.table[name.Lexeme]
	if !ok {
		return e.enclosing.Get(name)
	} else if _, ok := val.(Uninitialized); ok {
		e.tracker.Fatal(errtrack.LoxError{
			Message: errors.New("Variable uninitialized."),
			Token:   name,
		})
	}

	return val
}

func (e *Env) Assign(name tok.Token, val interface{}) {
	if e.table == nil {
		// We can't assign to an undeclared variable.
		e.tracker.Fatal(errtrack.ErrorUndefined(name))
	}

	if _, ok := e.table[name.Lexeme]; !ok {
		e.enclosing.Assign(name, val)
	}

	e.table[name.Lexeme] = val
}
