package errtrack

import (
	"fmt"
	"io"
	"os"

	"github.com/spencer-p/craftinginterpreters/pkg/lox/tok"
)

type LoxError struct {
	Message error
	Token   tok.Token
}

func (e LoxError) Error() string {
	return fmt.Sprintf("[line %d:%d] at %q: %v", e.Token.Line, e.Token.Char, e.Token.Lexeme, e.Message)
}

func (e LoxError) String() string {
	return e.Error()
}

func (e LoxError) Unwrap() error {
	return e.Message
}

// Tracker tracks errors that may happen deep in the call stack.
type Tracker struct {
	hadError bool
	output   io.Writer
}

func New() *Tracker {
	return &Tracker{
		hadError: false,
		output:   os.Stdout,
	}
}

// Report logs an error to output and makes a note there was an error.
func (t *Tracker) Report(err LoxError) {
	t.hadError = true
	fmt.Fprintf(t.output, "%s\n", err.Error())
}

// Fatal logs an error, notes it, and panics.
func (t *Tracker) Fatal(err LoxError) {
	t.Report(err)
	panic(err)
}

// HadError returns true if Report or Fatal were called until it is Reset.
func (t *Tracker) HadError() bool {
	return t.hadError
}

// Reset clears any errors. HadError returns false after a Reset.
func (t *Tracker) Reset() {
	t.hadError = false
}

// CatchFatal stops any calls to Tracker.Fatal from escaping a function. Must be
// deferred.
func (t *Tracker) CatchFatal() {
	if r := recover(); r != nil && !t.HadError() {
		// If we panicked from something but there was no tracker error, then
		// continue panicking.
		// Other panics are dropped.
		panic(r)
	}
}
