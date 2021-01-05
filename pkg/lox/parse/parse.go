package parse

import (
	"errors"

	"github.com/spencer-p/craftinginterpreters/pkg/lox/errtrack"
	"github.com/spencer-p/craftinginterpreters/pkg/lox/expr"
	"github.com/spencer-p/craftinginterpreters/pkg/lox/stmt"
	. "github.com/spencer-p/craftinginterpreters/pkg/lox/tok"
)

type Parser struct {
	tokens  []Token
	current int
	tracker *errtrack.Tracker
}

func New(tracker *errtrack.Tracker, toks []Token) *Parser {
	return &Parser{
		tokens:  toks,
		current: 0,
		tracker: tracker,
	}
}

func (p *Parser) AST() []stmt.Type {
	return p.parse()
}

func (p *Parser) parse() []stmt.Type {
	var statements []stmt.Type
	for !p.atEnd() {
		statements = append(statements, p.declaration())
	}
	return statements
}

func (p *Parser) declaration() stmt.Type {
	defer p.tracker.CatchFatal(p.synchronize)
	if p.match(VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) varDeclaration() stmt.Type {
	name := p.consume(IDENT, "Expect variable name.")

	var init expr.Type
	if p.match(EQUAL) {
		init = p.expression()
	}

	p.consume(SEMICOLON, "Expect ';' after variable declaration.")
	return &stmt.Var{name, init}
}

func (p *Parser) statement() stmt.Type {
	if p.match(PRINT) {
		return p.printStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) printStatement() stmt.Type {
	val := p.expression()
	p.consume(SEMICOLON, "Expect ';' after value.")
	return &stmt.Print{val}
}

func (p *Parser) expressionStatement() stmt.Type {
	e := p.expression()
	p.consume(SEMICOLON, "Expect ';' after value.")
	return &stmt.Expression{e}
}

func (p *Parser) expression() expr.Type {
	return p.assignment()
}

func (p *Parser) assignment() expr.Type {
	e := p.equality()

	if p.match(EQUAL) {
		equals := p.previous()
		right := p.assignment()
		switch left := e.(type) {
		case *expr.Variable:
			return &expr.Assign{left.Name, right}
		default:
			p.tracker.Report(errtrack.LoxError{
				Message: errors.New("Invalid assignment target."),
				Token:   equals,
			})
		}
	}

	return e
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
	} else if p.match(IDENT) {
		return &expr.Variable{p.previous()}
	}

	p.tracker.Fatal(errtrack.LoxError{
		Message: errors.New("Expected expression."),
		Token:   p.peek(),
	})
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

func (p *Parser) consume(typ TokenType, msg string) Token {
	if p.check(typ) {
		return p.advance()
	}

	p.tracker.Fatal(errtrack.LoxError{
		Message: errors.New(msg),
		Token:   p.peek(),
	})
	return Token{} // unreachable
}

func (p *Parser) synchronize() {
	for p.advance(); !p.atEnd(); p.advance() {
		if p.previous().Typ == SEMICOLON {
			return
		}
		switch p.peek().Typ {
		case CLASS:
			return
		case FN:
			return
		case VAR:
			return
		case FOR:
			return
		case IF:
			return
		case WHILE:
			return
		case PRINT:
			return
		case RETURN:
			return
		}
	}
}
