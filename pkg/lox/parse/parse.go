package parse

import (
	"github.com/spencer-p/craftinginterpreters/pkg/lox/expr"
	"github.com/spencer-p/craftinginterpreters/pkg/lox/tok"
)

type Parser struct {
	tokens  []tok.Token
	current int
}

func NewParser(toks []tok.Token) *Parser {
	return &Parser{
		tokens:  toks,
		current: 0,
	}
}

func (p *Parser) AST() expr.Type {
	return p.expression()
}

func (p *Parser) expression() expr.Type {
	// TODO
	return nil
}
