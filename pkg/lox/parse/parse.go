package parse

import (
	"github.com/spencer-p/craftinginterpreters/pkg/lox/expr"
	. "github.com/spencer-p/craftinginterpreters/pkg/lox/tok"
)

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(toks []Token) *Parser {
	return &Parser{
		tokens:  toks,
		current: 0,
	}
}

func (p *Parser) AST() expr.Type {
	return p.expression()
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
	return p.primary()
}

func (p *Parser) primary() expr.Type {
	if p.match(NUMBER) {
		return &expr.Literal{p.previous().Lit}
	}

	if p.match(TRUE) {
		return &expr.Literal{true}
	}

	if p.match(FALSE) {
		return &expr.Literal{false}
	}

	return nil
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
