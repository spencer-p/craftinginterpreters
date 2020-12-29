package interpret

import (
	"errors"
	"fmt"

	"github.com/spencer-p/craftinginterpreters/pkg/lox/tok"
)

var (
	ErrorNotANumber = errors.New("Operand must be number.")
	ErrorNotAString = errors.New("Operand must be string.")
)

type RuntimeError struct {
	Message error
	Token   tok.Token
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("[line %d:%d] at %q: %v", e.Token.Line, e.Token.Char, e.Token.Lexeme, e.Message)
}

func (e RuntimeError) String() string {
	return e.Error()
}

func (e RuntimeError) Unwrap() error {
	return e.Message
}
