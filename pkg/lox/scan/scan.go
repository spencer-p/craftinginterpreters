package scan

import (
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/spencer-p/craftinginterpreters/pkg/lox/errtrack"
	. "github.com/spencer-p/craftinginterpreters/pkg/lox/tok"
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

	RESERVED = map[string]TokenType{
		"and":    AND,
		"class":  CLASS,
		"else":   ELSE,
		"false":  FALSE,
		"for":    FOR,
		"fn":     FN, // I prefer fn over Lox's fun.
		"fun":    FN, // However, we support both.
		"if":     IF,
		"nil":    NIL,
		"or":     OR,
		"print":  PRINT,
		"return": RETURN,
		"super":  SUPER,
		"this":   THIS,
		"true":   TRUE,
		"var":    VAR,
		"while":  WHILE,
	}
)

type Scanner struct {
	src        string
	tokens     []Token
	start, cur int
	line       int

	lookahead [2]struct {
		char  rune
		width int
	}
	lookaheadi int
}

func New(src string) *Scanner {
	return &Scanner{
		src:        src,
		tokens:     make([]Token, 0),
		line:       1,
		lookaheadi: -1,
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
	r := s.advance()

	// simple lookups first
	if tok, ok := ONEWIDTHS[r]; ok {
		s.addToken1(tok)
		return
	}

	// switch on larger statments
	switch r {
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
			// this is a comment -- consume it
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
		if unicode.IsDigit(r) {
			s.eatNumber()
		} else if isAlphaNum(r) {
			s.eatIdent()
		} else {
			errtrack.Complain(s.line, "unexpected rune %q", r)
		}
	}

}

func (s *Scanner) peek() rune {
	if s.atEnd() {
		return 0
	}

	if s.lookaheadi >= 0 {
		return s.lookahead[0].char
	}

	// Using UTF-8 package to get the full char - thanks Go!
	next, width := utf8.DecodeRuneInString(s.src[s.cur:])
	s.lookahead[0].char = next
	s.lookahead[0].width = width
	s.lookaheadi = 0
	return next
}

func (s *Scanner) peekNext() rune {
	if s.atEnd() {
		return 0 // TODO will EOF be detected properly here?
	}

	if s.lookaheadi >= 1 {
		return s.lookahead[1].char
	}

	// Get offset we should be reading from
	// by making sure a peek already happened
	if s.peek() == 0 {
		return 0 // nothing to peek at
	}

	// we now know the peeked character's width. read the peek next
	offset := s.lookahead[0].width
	next, width := utf8.DecodeRuneInString(s.src[s.cur+offset:])
	s.lookahead[1].char = next
	s.lookahead[1].width = width
	s.lookaheadi = 1
	return next
}

func (s *Scanner) advance() rune {
	s.peek()
	s.cur += s.lookahead[0].width
	char := s.lookahead[0].char

	if s.lookaheadi == 1 {
		// there's a peekNext we can move
		s.lookahead[0] = s.lookahead[1]
	}
	s.lookaheadi -= 1

	return char
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
		errtrack.Complain(s.line, "unterminated string")
		return
	}

	s.advance() // corresponds to last "

	// cut out the quotes from the value to add
	// TODO Implement escape sequences
	val := s.src[s.start+1 : s.cur-1]
	s.addToken(STRING, val)
}

func (s *Scanner) eatNumber() {
	for unicode.IsDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && unicode.IsDigit(s.peekNext()) {
		s.advance() // the dot character

		for unicode.IsDigit(s.peek()) {
			s.advance()
		}
	}

	val, err := strconv.ParseFloat(s.src[s.start:s.cur], 64)
	if err != nil {
		errtrack.Complain(s.line, "number does not parse: %v", err)
		return
	}
	s.addToken(NUMBER, val)
}

func (s *Scanner) eatIdent() {
	for isAlphaNum(s.peek()) {
		s.advance()
	}

	val := s.src[s.start:s.cur]
	if tok, ok := RESERVED[val]; ok {
		s.addToken1(tok)
	} else {
		s.addToken1(IDENT)
	}
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

// true if r is alphanumeric (in L, M, N, So (includes emoji) or '_')
func isAlphaNum(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsNumber(r) || r == '_' || unicode.IsMark(r) || unicode.Is(unicode.So, r)
}
