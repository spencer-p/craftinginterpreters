package prettyprint

import (
	"fmt"

	"github.com/spencer-p/craftinginterpreters/pkg/lox/expr"
	"github.com/spencer-p/craftinginterpreters/pkg/lox/tok"
)

func ExampleLisp() {
	e := expr.Binary{
		Op: tok.Token{Lexeme: "*"},
		Left: &expr.Unary{
			Op:    tok.Token{Lexeme: "-"},
			Right: &expr.Literal{123},
		},
		Right: &expr.Grouping{
			&expr.Literal{45.67},
		},
	}

	result := e.Accept(&Lisp{})
	fmt.Println(result)

	// Output:
	// (* (- 123) (grp 45.67))
}
