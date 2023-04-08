package parser

import (
	"testing"

	"github.com/muggel/emlang/ast"
	"github.com/muggel/emlang/examples"
	"github.com/muggel/emlang/scanner"
	"github.com/muggel/emlang/token"
	"github.com/stretchr/testify/assert"
)

func TestParser_parseIntLiteral(t *testing.T) {
	t.Run("parses_int_literal", func(t *testing.T) {
		s := scanner.NewScanner("123")
		p := NewParser(s)
		res := p.parseIntLiteral()
		expected := &ast.IntLiteral{Token: token.INT, Literal: "123", Value: 123}
		assert.Equal(t, expected, res)
	})

	t.Run("adds_error_to_parser_if_int_literal_is_invalid", func(t *testing.T) {
		s := scanner.NewScanner("abc")
		p := NewParser(s)
		p.parseIntLiteral()
		assert.Equal(t, 1, len(p.Errors))
	})
}

func TestParser_parseIdentifier(t *testing.T) {
	s := scanner.NewScanner("abc")
	p := NewParser(s)
	res := p.parseIdentifier()
	expected := &ast.Identifier{Token: token.IDENT, Literal: "abc", Value: "abc"}
	assert.Equal(t, expected, res)
}

func TestParser_parseCallExpression(t *testing.T) {
	s := scanner.NewScanner("func()")
	p := NewParser(s)
	res := p.parseCallExpression()
	expected := &ast.CallExpression{Token: token.IDENT, Function: &ast.Identifier{Token: token.IDENT, Literal: "func", Value: "func"}}
	assert.Equal(t, expected, res)
}

func TestParser_parseReturnStatement(t *testing.T) {
	t.Run("parses_return_statement", func(t *testing.T) {
		s := scanner.NewScanner("return 123;")
		p := NewParser(s)
		res := p.parseReturnStatement()
		expected := &ast.ReturnStatement{Token: token.RETURN, ReturnValue: &ast.IntLiteral{Token: token.INT, Literal: "123", Value: 123}}
		assert.Equal(t, expected, res)
	})

	t.Run("adds_error_to_parser_if_return_statement_is_missing_semicolon", func(t *testing.T) {
		s := scanner.NewScanner("return 123")
		p := NewParser(s)
		p.parseReturnStatement()
		assert.Equal(t, 1, len(p.Errors))
	})
}

func TestParser_parseAssignmentStatement(t *testing.T) {
	t.Run("parses_assignment_statement", func(t *testing.T) {
		s := scanner.NewScanner("foo = 123;")
		p := NewParser(s)
		res := p.parseAssignmentStatement()
		expected := &ast.AssignmentStatement{
			Token:      token.ASSIGN,
			Identifier: &ast.Identifier{Token: token.IDENT, Literal: "foo", Value: "foo"},
			Value:      &ast.IntLiteral{Token: token.INT, Literal: "123", Value: 123},
		}

		assert.Equal(t, expected, res)
	})

	t.Run("adds_error_to_parser_if_assignment_statement_is_missing_semicolon", func(t *testing.T) {
		s := scanner.NewScanner("foo = 123")
		p := NewParser(s)
		p.parseAssignmentStatement()
		assert.Equal(t, 1, len(p.Errors))
	})

	t.Run("adds_error_to_parser_if_assignment_statement_does_not_contain_assignment_token", func(t *testing.T) {
		s := scanner.NewScanner("foo - 123;")
		p := NewParser(s)
		p.parseAssignmentStatement()
		assert.Equal(t, 1, len(p.Errors))
	})
}

func TestParser_parseBlockStatement(t *testing.T) {
	t.Run("parses_block_statement", func(t *testing.T) {
		s := scanner.NewScanner("{\nfoo = 123;\nreturn foo;\n}")
		p := NewParser(s)
		res := p.parseBlockStatement()
		expected := &ast.BlockStatement{
			Token: token.LBRACE,
			Statements: []ast.Statement{
				&ast.AssignmentStatement{
					Token:      token.ASSIGN,
					Identifier: &ast.Identifier{Token: token.IDENT, Literal: "foo", Value: "foo"},
					Value:      &ast.IntLiteral{Token: token.INT, Literal: "123", Value: 123},
				},
				&ast.ReturnStatement{
					Token:       token.RETURN,
					ReturnValue: &ast.Identifier{Token: token.IDENT, Literal: "foo", Value: "foo"},
				},
			},
		}

		assert.Equal(t, expected, res)
	})

	t.Run("adds_error_to_parser_if_block_statement_is_missing_closing_brace", func(t *testing.T) {
		s := scanner.NewScanner("{return 123;")
		p := NewParser(s)
		p.parseBlockStatement()
		assert.Equal(t, 1, len(p.Errors))
	})
}

func TestParser_parseFunctionDeclaration(t *testing.T) {
	t.Run("parses_function_declarations", func(t *testing.T) {
		s := scanner.NewScanner("fn foo() int {\nreturn foo;\n}")
		p := NewParser(s)
		res := p.parseFunctionDeclaration()
		expected := &ast.FunctionDeclaration{
			Token:      token.FN,
			Literal:    "fn",
			Identifier: &ast.Identifier{Token: token.IDENT, Literal: "foo", Value: "foo"},
			ReturnType: &ast.Identifier{Token: token.IDENT, Literal: "int", Value: "int"},
			Body: &ast.BlockStatement{
				Token: token.LBRACE,
				Statements: []ast.Statement{
					&ast.ReturnStatement{
						Token:       token.RETURN,
						ReturnValue: &ast.Identifier{Token: token.IDENT, Literal: "foo", Value: "foo"},
					},
				},
			},
		}
		assert.Equal(t, expected, res)
	})

	t.Run("handles_void_as_return_type", func(t *testing.T) {
		s := scanner.NewScanner("fn foo() {\nabc = 123;\n}")
		p := NewParser(s)
		res := p.parseFunctionDeclaration()
		expected := &ast.FunctionDeclaration{
			Token:      token.FN,
			Literal:    "fn",
			Identifier: &ast.Identifier{Token: token.IDENT, Literal: "foo", Value: "foo"},
			ReturnType: &ast.Identifier{Token: token.IDENT, Literal: "void", Value: "void"},
			Body: &ast.BlockStatement{
				Token: token.LBRACE,
				Statements: []ast.Statement{
					&ast.AssignmentStatement{
						Token:      token.ASSIGN,
						Literal:    "=",
						Identifier: &ast.Identifier{Token: token.IDENT, Literal: "abc", Value: "abc"},
						Value:      &ast.IntLiteral{Token: token.INT, Literal: "123", Value: 123},
					},
				},
			},
		}
		assert.Equal(t, expected, res)
	})
}

func TestExamples(t *testing.T) {
	t.Run("main_an_helper.em", func(t *testing.T) {
		s := scanner.NewScanner(examples.MainAndHelper)
		p := NewParser(s)
		program := p.parseProgram()
		expected := &ast.Program{
			TopLevelDeclarations: []ast.TopLevelDeclaration{
				&ast.FunctionDeclaration{
					Token:      token.FN,
					Literal:    "fn",
					Identifier: &ast.Identifier{Token: token.IDENT, Literal: "helper", Value: "helper"},
					ReturnType: &ast.Identifier{Token: token.IDENT, Literal: "int", Value: "int"},
					Body: &ast.BlockStatement{
						Token:   token.LBRACE,
						Literal: "{",
						Statements: []ast.Statement{
							&ast.ReturnStatement{
								Token:       token.RETURN,
								Literal:     "return",
								ReturnValue: &ast.IntLiteral{Token: token.INT, Literal: "123", Value: 123},
							},
						},
					},
				},
				&ast.FunctionDeclaration{
					Token:      token.FN,
					Literal:    "fn",
					Identifier: &ast.Identifier{Token: token.IDENT, Literal: "main", Value: "main"},
					ReturnType: &ast.Identifier{Token: token.IDENT, Literal: "void", Value: "void"},
					Body: &ast.BlockStatement{
						Token:   token.LBRACE,
						Literal: "{",
						Statements: []ast.Statement{
							&ast.AssignmentStatement{
								Token:      token.ASSIGN,
								Literal:    "=",
								Identifier: &ast.Identifier{Token: token.IDENT, Literal: "foo", Value: "foo"},
								Value: &ast.CallExpression{
									Token:    token.IDENT,
									Literal:  "helper",
									Function: &ast.Identifier{Token: token.IDENT, Literal: "helper", Value: "helper"},
								},
							},
						},
					},
				},
			},
		}

		assert.Equal(t, expected, program)
	})
}
