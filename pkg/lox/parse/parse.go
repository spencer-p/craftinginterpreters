package parse

import (
	"fmt"

	"github.com/spencer-p/craftinginterpreters/pkg/lox/expr"
	. "github.com/spencer-p/craftinginterpreters/pkg/lox/tok"
)

type Parser struct {
	tokens  []Token
	current int
}

func New(toks []Token) *Parser {
	return &Parser{
		tokens:  toks,
		current: 0,
	}
}

func (p *Parser) AST() (e expr.Type, err error) {
	defer func() {
		// TODO - accumulate multiple errors somehow
		// Type assert will have no effect if it fails
		err, _ = recover().(error)
	}()
	e = p.expression()
	return
}

func (p *Parser) expression() expr.Type {
	return p.equality()
}

func (p *Parser) equality() expr.Type {
	e := p.comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		op := p.previous()
		right := p.comparison()
		e = &expr.Binary{
			Left:  e,
			Right: right,
			Op:    op,
		}
	}

	return e
}

func (p *Parser) comparison() expr.Type {
	e := p.addition()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		op := p.previous()
		right := p.addition()
		e = &expr.Binary{
			Left:  e,
			Right: right,
			Op:    op,
		}
	}

	return e
}

func (p *Parser) addition() expr.Type {
	e := p.multiplication()

	for p.match(MINUS, PLUS) {
		op := p.previous()
		right := p.multiplication()
		e = &expr.Binary{
			Left:  e,
			Right: right,
			Op:    op,
		}
	}

	return e
}

func (p *Parser) multiplication() expr.Type {
	e := p.unary()

	for p.match(SLASH, STAR) {
		op := p.previous()
		right := p.unary()
		e = &expr.Binary{
			Left:  e,
			Right: right,
			Op:    op,
		}
	}

	return e
}

func (p *Parser) unary() expr.Type {
	if p.match(BANG, MINUS) {
		op := p.previous()
		right := p.unary()
		return &expr.Unary{
			Op:    op,
			Right: right,
		}
	}

	return p.primary()
}

func (p *Parser) primary() expr.Type {
	if p.match(TRUE) {
		return &expr.Literal{true}
	} else if p.match(FALSE) {
		return &expr.Literal{false}
	} else if p.match(NIL) {
		return &expr.Literal{nil}
	} else if p.match(NUMBER, STRING) {
		return &expr.Literal{p.previous().Lit}
	} else if p.match(LEFT_PAREN) {
		e := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return &expr.Grouping{e}
	}

	panic(err(p.peek(), "Expected expression."))
}

func (p *Parser) match(types ...TokenType) bool {
	for _, typ := range types {
		if p.check(typ) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(typ TokenType) bool {
	if p.atEnd() {
		return false
	}
	return p.tokens[p.current].Typ == typ
}

func (p *Parser) atEnd() bool {
	return p.current >= len(p.tokens) || p.peek().Typ == EOF
}

func (p *Parser) advance() Token {
	if p.atEnd() == false {
		p.current += 1
	}
	return p.previous()
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(typ TokenType, msg string) Token {
	if p.check(typ) {
		return p.advance()
	}

	panic(err(p.peek(), msg))
}

func err(tok Token, msg string) error {
	spot := tok.Lexeme
	if tok.Typ == EOF {
		spot = `EOF`
	}

	return fmt.Errorf("[line %d:%d] at %q: %s", tok.Line, tok.Char, spot, msg)
}
