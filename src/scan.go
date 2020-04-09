package main

import (
	"unicode/utf8"
)

var (
	ONEWIDTHS = map[rune]TokenType{
		'(': LEFT_PAREN,
		')': RIGHT_PAREN,
		'{': LEFT_BRACE,
		'}': RIGHT_BRACE,
		',': COMMA,
		'.': DOT,
		'-': MINUS,
		'+': PLUS,
		';': SEMICOLON,
		'*': STAR,
	}
)

type Scanner struct {
	src        string
	tokens     []Token
	start, cur int
	currune    rune
	curwidth   int
	line       int
}

func NewScanner(src string) *Scanner {
	return &Scanner{
		src:    src,
		tokens: make([]Token, 0),
		line:   1,
	}
}

// Tokens scans the input source and returns its tokens.
func (s *Scanner) Tokens() []Token {
	for !s.atEnd() {
		s.start = s.cur
		s.scanToken()
	}
	s.tokens = append(s.tokens, Token{EOF, "", nil, s.line})
	return s.tokens
}

func (s *Scanner) atEnd() bool {
	return s.cur >= len(s.src)
}

func (s *Scanner) scanToken() {
	c := s.advance()

	// simple lookups first
	if tok, ok := ONEWIDTHS[c]; ok {
		s.addToken1(tok)
		return
	}

	// switch on larger statments
	switch c {
	case '!':
		s.addToken1(s.match('=', BANG_EQUAL, BANG))
	case '=':
		s.addToken1(s.match('=', EQUAL_EQUAL, EQUAL))
	case '<':
		s.addToken1(s.match('=', LESS_EQUAL, LESS))
	case '>':
		s.addToken1(s.match('=', GREATER_EQUAL, GREATER))
	case '/':
		if s.peek() == '/' {
			// this is a consume -- consume it
			for s.peek() != '\n' && !s.atEnd() {
				s.advance()
			}
		} else {
			s.addToken1(SLASH)
		}
	case '"':
		s.eatString()
	case '\n':
		s.line += 1
		fallthrough
	case ' ', '\r', '\t':
		break
	default:
		complain(s.line, "Unexpected codepoint %q", c)
	}

}

func (s *Scanner) peek() rune {
	if s.atEnd() {
		return 0
	}

	// if width was reset, we have to do a fetch
	if s.curwidth == 0 {
		// Using UTF-8 package to get the full char - thanks Go!
		next, width := utf8.DecodeRuneInString(s.src[s.cur:])
		s.curwidth = width
		s.currune = next
		return next
	} else {
		// if a width is still stored, we still have the rune
		return s.currune
	}
}

func (s *Scanner) advance() rune {
	s.peek()
	s.cur += s.curwidth
	s.curwidth = 0
	return s.currune
}

func (s *Scanner) match(expect rune, match, nomatch TokenType) TokenType {
	if s.atEnd() || s.peek() != expect {
		return nomatch
	}

	// here we performed peek and found what we expected
	s.advance()
	return match
}

func (s *Scanner) eatString() {
	for s.peek() != '"' && !s.atEnd() {
		if s.peek() == '\n' {
			s.line += 1
		}
		s.advance()
	}

	// the document ended before the string..
	if s.atEnd() {
		complain(s.line, "unterminated string")
		return
	}

	s.advance() // corresponds to last "

	// cut out the quotes from the value to add
	// TODO Implement escape sequences
	val := s.src[s.start+1 : s.cur-1]
	s.addToken(STRING, val)
}

func (s *Scanner) addToken(tok TokenType, lit interface{}) {
	s.tokens = append(s.tokens, Token{
		tok,
		s.src[s.start:s.cur],
		lit,
		s.line})
}

func (s *Scanner) addToken1(tok TokenType) {
	s.addToken(tok, nil)
}
