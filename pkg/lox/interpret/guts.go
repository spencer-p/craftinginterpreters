package interpret

import (
	"fmt"
	"strings"

	"github.com/spencer-p/craftinginterpreters/pkg/lox/tok"
)

func truthy(value interface{}) bool {
	switch actual := value.(type) {
	case nil:
		return false
	case bool:
		return actual
	default:
		return true
	}
}

func equal(a, b interface{}) (result bool) {
	defer func() {
		// Catch failed type casting and simply return false.
		if err := recover(); err != nil {
			result = false
		}
	}()

	switch actual := a.(type) {
	case nil:
		return b == nil
	case bool:
		return actual == b.(bool)
	case string:
		return actual == b.(string)
	case float64:
		return actual == b.(float64)
	default:
		return false
	}
}

func checkNumber(op tok.Token, value interface{}) {
	if _, ok := value.(float64); !ok {
		err := RuntimeError{
			Message: ErrorNotANumber,
			Token:   op,
		}
		panic(err)
	}
}

func checkNumbers(op tok.Token, values ...interface{}) {
	for _, v := range values {
		checkNumber(op, v)
	}
}

func Stringify(result interface{}, err error) string {
	if err != nil {
		return err.Error()
	}
	s := fmt.Sprintf("%v", result)
	if _, ok := result.(float64); ok && strings.HasSuffix(s, ".0") {
		// Remove .0 only from floats
		s = s[:len(s)-2]
	}
	return s
}
