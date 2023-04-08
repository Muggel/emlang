package parser

import (
	"errors"
	"strconv"

	"github.com/muggel/emlang/ast"
	"github.com/muggel/emlang/scanner"
	"github.com/muggel/emlang/token"
)

type Parser struct {
	s *scanner.Scanner

	currentToken   token.Token
	currentLiteral string
	peekToken      token.Token
	peekLiteral    string
	hasPeeked      bool

	Errors []error
}

func NewParser(s *scanner.Scanner) *Parser {
	parser := &Parser{s: s}

	// fill both current and peek
	parser.readNext()
	parser.readNext()

	return parser
}

func (p *Parser) Parse() (ast.Node, error) {
	return nil, nil
}

func (p *Parser) parseProgram() *ast.Program {
	program := &ast.Program{}

	for p.currentToken != token.EOF {
		decl := p.parseTopLevelDeclaration()
		if decl != nil {
			program.TopLevelDeclarations = append(program.TopLevelDeclarations, decl)
		}
		// p.readNext()
	}

	return program
}

func (p *Parser) parseTopLevelDeclaration() ast.TopLevelDeclaration {
	if p.currentToken == token.FN {
		return p.parseFunctionDeclaration()
	}
	return nil
}

func (p *Parser) parseFunctionDeclaration() *ast.FunctionDeclaration {
	stmt := &ast.FunctionDeclaration{
		Token:   p.currentToken,
		Literal: p.currentLiteral,
	}
	p.readNext()

	stmt.Identifier = p.parseIdentifier()
	p.readNext()

	if p.currentToken != token.LPAREN || p.peekToken != token.RPAREN {
		p.Errors = append(p.Errors, errors.New("expected empty parameter list"))
	}
	p.readNext()
	p.readNext()

	if p.currentToken != token.LBRACE {
		stmt.ReturnType = p.parseIdentifier()
		p.readNext()
	} else {
		stmt.ReturnType = &ast.Identifier{Token: token.IDENT, Literal: "void", Value: "void"}
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currentToken, Literal: p.currentLiteral}
	p.readNext()

	for p.currentToken != token.RBRACE && p.currentToken != token.EOF {
		stmt := p.parseStatement()
		block.Statements = append(block.Statements, stmt)
	}
	if p.currentToken == token.EOF {
		p.Errors = append(p.Errors, errors.New("unexpected end of file"))
	}

	p.readNext()

	return block
}

func (p *Parser) parseStatement() ast.Statement {
	if p.currentToken == token.RETURN {
		return p.parseReturnStatement()
	}
	if p.currentToken == token.IDENT && p.peekToken == token.ASSIGN {
		return p.parseAssignmentStatement()
	}
	return nil
}

func (p *Parser) parseAssignmentStatement() *ast.AssignmentStatement {
	if p.peekToken != token.ASSIGN || p.currentToken != token.IDENT {
		p.Errors = append(p.Errors, errors.New("expected single identifier before assignment operator"))
	}
	stmt := &ast.AssignmentStatement{Token: p.peekToken, Literal: p.peekLiteral}
	stmt.Identifier = p.parseIdentifier()
	// read the assignment token
	p.readNext()
	// read the value
	p.readNext()

	value := p.parseExpression()
	stmt.Value = value
	p.readNext()

	p.consumeSemicolon()

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currentToken, Literal: p.currentLiteral}
	p.readNext()

	stmt.ReturnValue = p.parseExpression()
	p.readNext()

	p.consumeSemicolon()

	return stmt
}

func (p *Parser) parseExpression() ast.Expression {
	switch p.currentToken {
	case token.IDENT:
		if p.peekToken == token.LPAREN {
			return p.parseCallExpression()
		} else {
			return p.parseIdentifier()
		}
	case token.INT:
		return p.parseIntLiteral()
	default:
		p.Errors = append(p.Errors, errors.New("unknown expression"))
		return nil
	}
}

func (p *Parser) parseIntLiteral() *ast.IntLiteral {
	intValue, err := strconv.ParseInt(p.currentLiteral, 10, 64)
	if err != nil {
		p.Errors = append(p.Errors, errors.New("could not parse int literal"))
	}
	return &ast.IntLiteral{Token: p.currentToken, Literal: p.currentLiteral, Value: intValue}
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	return &ast.Identifier{Token: p.currentToken, Literal: p.currentLiteral, Value: p.currentLiteral}
}

func (p *Parser) parseCallExpression() *ast.CallExpression {
	call := &ast.CallExpression{Token: p.currentToken, Literal: p.currentLiteral}
	call.Function = p.parseIdentifier()
	p.readNext()

	if p.currentToken != token.LPAREN {
		p.Errors = append(p.Errors, errors.New("expected left parenthesis"))
	}
	p.readNext()
	if p.currentToken != token.RPAREN {
		p.Errors = append(p.Errors, errors.New("expected closing parenthesis"))
	}

	return call
}

func (p *Parser) readNext() {
	p.currentToken, p.currentLiteral = p.peekToken, p.peekLiteral
	p.peekToken, p.peekLiteral = p.s.Next()
}

func (p *Parser) consumeSemicolon() {
	if p.currentToken != token.SEMICOLON {
		p.Errors = append(p.Errors, errors.New("expected semicolon"))
	}
	p.readNext()
}
