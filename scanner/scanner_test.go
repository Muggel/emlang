package scanner

import (
	"testing"

	"github.com/muggel/emlang/examples"
	"github.com/muggel/emlang/token"
	"github.com/stretchr/testify/assert"
)

type tokenLitPair struct {
	token token.Token
	lit   string
}

func TestScanner__Next(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected []tokenLitPair
	}{
		{
			name:   "scans_empty_source",
			source: "",
			expected: []tokenLitPair{
				{token.EOF, ""},
			},
		},
		{
			name:   "scans_handles_illegal_tokens",
			source: "@",
			expected: []tokenLitPair{
				{token.ILLEGAL, "@"}, {token.EOF, ""},
			},
		},
		{
			name:   "scans_whitespace",
			source: " 	\n\r",
			expected: []tokenLitPair{
				{token.EOF, ""},
			},
		},
		{
			name:   "scans_binary_operators",
			source: "+-*/",
			expected: []tokenLitPair{
				{token.ADD, "+"}, {token.SUB, "-"}, {token.MUL, "*"}, {token.DIV, "/"}, {token.EOF, ""},
			},
		},
		{
			name:   "scans_brackets",
			source: "(){}",
			expected: []tokenLitPair{
				{token.LPAREN, "("}, {token.RPAREN, ")"}, {token.LBRACE, "{"}, {token.RBRACE, "}"}, {token.EOF, ""},
			},
		},
		{
			name:   "scans_semicolon_comma_and_assign",
			source: ",=;",
			expected: []tokenLitPair{
				{token.COMMA, ","}, {token.ASSIGN, "="}, {token.SEMICOLON, ";"}, {token.EOF, ""},
			},
		},
		{
			name:   "scans_integers",
			source: "123",
			expected: []tokenLitPair{
				{token.INT, "123"}, {token.EOF, ""},
			},
		},
		{
			name:   "scans_identifiers",
			source: "_abc",
			expected: []tokenLitPair{
				{token.IDENT, "_abc"}, {token.EOF, ""},
			},
		},
		{
			name:   "scans_keywords",
			source: "fn return",
			expected: []tokenLitPair{
				{token.FN, "fn"}, {token.RETURN, "return"}, {token.EOF, ""},
			},
		},
		{
			name:   "scans_function",
			source: examples.Function,
			expected: []tokenLitPair{
				{token.FN, "fn"},
				{token.IDENT, "main"},
				{token.LPAREN, "("},
				{token.RPAREN, ")"},
				{token.IDENT, "int"},
				{token.LBRACE, "{"},
				{token.RETURN, "return"},
				{token.INT, "1"},
				{token.SEMICOLON, ";"},
				{token.RBRACE, "}"},
				{token.EOF, ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewScanner(tt.source)
			res := scanAll(t, s)
			assert.Equal(t, tt.expected, res)
		})
	}
}

func scanAll(t *testing.T, s *Scanner) (res []tokenLitPair) {
	for {
		tok, lit := s.Next()
		res = append(res, tokenLitPair{tok, lit})
		if tok == token.EOF {
			break
		}
	}
	return res
}
