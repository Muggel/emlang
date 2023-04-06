package scanner

import (
	"github.com/muggel/emlang/token"
)

const _eof = 0

type Scanner struct {
	source string

	position     int
	readPosition int
	ch           byte
}

func NewScanner(source string) *Scanner {
	s := &Scanner{source: source}
	s.readChar()
	return s
}

func (s *Scanner) Next() (token.Token, string /* literal */) {
	var tok token.Token

	s.skipWhitespace()

	literal := string(s.ch)
	switch s.ch {
	case '+':
		tok = token.ADD
	case '-':
		tok = token.SUB
	case '*':
		tok = token.MUL
	case '/':
		tok = token.DIV

	case '=':
		tok = token.ASSIGN
	case ',':
		tok = token.COMMA

	case '(':
		tok = token.LPAREN
	case ')':
		tok = token.RPAREN
	case '{':
		tok = token.LBRACE
	case '}':
		tok = token.RBRACE

	case 0:
		tok = token.EOF
		literal = ""
	default:
		if isNumber(s.ch) {
			literal = s.readNumber()
			tok = token.INT
			break
		} else if isLetter(s.ch) {
			literal = s.readIdentifier()
			tok = token.Lookup(literal)
			break
		} else {
			tok = token.ILLEGAL
		}
	}

	s.readChar()
	return tok, literal
}

func (s *Scanner) readIdentifier() string {
	start := s.position
	for isLetter(s.peekChar()) {
		s.readChar()
	}

	return s.source[start:s.readPosition]
}

func (s *Scanner) readNumber() string {
	start := s.position
	for isNumber(s.peekChar()) {
		s.readChar()
	}

	return s.source[start:s.readPosition]
}

func (s *Scanner) readChar() {
	s.ch = s.peekChar()
	s.position = s.readPosition
	s.readPosition++
}

func (s *Scanner) peekChar() byte {
	if s.readPosition >= len(s.source) {
		return _eof
	}
	return s.source[s.readPosition]
}

func (s *Scanner) skipWhitespace() {
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\n' || s.ch == '\r' {
		s.readChar()
	}
}

func isNumber(ch byte) bool { return '0' <= ch && ch <= '9' }
func isLetter(ch byte) bool { return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' }
