package prettyprint

import (
	"os"

	"github.com/spencer-p/craftinginterpreters/pkg/lox/expr"
	"github.com/spencer-p/craftinginterpreters/pkg/lox/tok"
)

func ExampleTab() {
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

	printer := NewWithTabs(os.Stdout)
	e.Accept(printer)

	// Output:
	// *
	//   -
	//     123
	//   (
	//    45.67
	//   )
}
