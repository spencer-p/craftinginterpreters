package prettyprint

import (
	"fmt"
	"io"

	"github.com/spencer-p/craftinginterpreters/pkg/lox/expr"
)

type WithTabs struct {
	indentStack []int
	out         io.Writer
}

func NewWithTabs(out io.Writer) *WithTabs {
	return &WithTabs{
		out:         out,
		indentStack: make([]int, 0),
	}
}

func (p *WithTabs) push(i int) {
	i += p.peek()
	p.indentStack = append(p.indentStack, i)
}

func (p *WithTabs) pop() int {
	if len(p.indentStack) == 0 {
		return 0
	}

	n := len(p.indentStack)
	ret := p.indentStack[n-1]
	p.indentStack = p.indentStack[:n-1]
	return ret
}

func (p *WithTabs) peek() int {
	if len(p.indentStack) == 0 {
		return 0
	}
	return p.indentStack[len(p.indentStack)-1]
}

func (p *WithTabs) tabs() {
	for i := 0; i < p.peek(); i++ {
		p.out.Write([]byte{' '})
	}
}

func (p *WithTabs) VisitBinary(e *expr.Binary) interface{} {
	p.tabs()
	fmt.Fprintf(p.out, "%s\n", e.Op.Lexeme)

	p.push(len(e.Op.Lexeme) + 1)
	defer p.pop()

	e.Left.Accept(p)
	e.Right.Accept(p)

	return nil
}

func (p *WithTabs) VisitGrouping(e *expr.Grouping) interface{} {
	p.tabs()
	fmt.Fprintf(p.out, "(\n")

	p.push(1)
	e.Expression.Accept(p)
	p.pop()

	p.tabs()
	fmt.Fprintf(p.out, ")\n")

	return nil
}

func (p *WithTabs) VisitLiteral(e *expr.Literal) interface{} {
	p.tabs()
	fmt.Fprintf(p.out, "%+v\n", e.Value)
	return nil
}

func (p *WithTabs) VisitUnary(e *expr.Unary) interface{} {
	p.tabs()
	fmt.Fprintf(p.out, "%s\n", e.Op.Lexeme)
	p.push(len(e.Op.Lexeme) + 1)
	defer p.pop()
	e.Right.Accept(p)
	return nil
}
