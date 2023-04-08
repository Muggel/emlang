package ast

import (
	"strings"

	"github.com/muggel/emlang/token"
)

// INTERFACES
// this block describes the interfaces that are used in the AST
type Node interface {
	TokenLiteral() string
	String() string
}

type TopLevelDeclaration interface {
	Node
	topLevelDeclaration()
}

type Statement interface {
	Node
	statement()
}

type Expression interface {
	Node
	expression()
}

// TYPES
// this block describes the types that are used in the AST
type Program struct {
	TopLevelDeclarations []TopLevelDeclaration
}

func (p *Program) TokenLiteral() string {
	if len(p.TopLevelDeclarations) > 0 {
		return p.TopLevelDeclarations[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var res strings.Builder

	for _, decl := range p.TopLevelDeclarations {
		res.WriteString(decl.String())
	}

	return res.String()
}

type FunctionDeclaration struct {
	Token   token.Token
	Literal string

	Identifier *Identifier
	// TODO Add parameters to functions
	// Parameters []*Parameter
	ReturnType *Identifier
	Body       *BlockStatement
}

func (fd *FunctionDeclaration) topLevelDeclaration() {}
func (fd *FunctionDeclaration) TokenLiteral() string { return fd.Literal }
func (fd *FunctionDeclaration) String() string {
	var res strings.Builder

	res.WriteString(fd.Literal + " ")
	res.WriteString(fd.Identifier.String())
	res.WriteString("() ")
	res.WriteString(fd.ReturnType.String())
	res.WriteString(" ")
	res.WriteString(fd.Body.String())

	return res.String()
}

type BlockStatement struct {
	Token   token.Token
	Literal string

	Statements []Statement
}

func (bs *BlockStatement) statement()           {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Literal }
func (bs *BlockStatement) String() string {
	var res strings.Builder

	res.WriteString("{\n")

	for _, stmt := range bs.Statements {
		res.WriteString(stmt.String())
	}

	res.WriteString("}\n")

	return res.String()
}

type AssignmentStatement struct {
	Token   token.Token
	Literal string

	Identifier *Identifier
	Value      Expression
}

func (as *AssignmentStatement) statement()           {}
func (as *AssignmentStatement) TokenLiteral() string { return as.Literal }
func (as *AssignmentStatement) String() string {
	var res strings.Builder

	res.WriteString(as.Identifier.String())
	res.WriteString(" = ")
	res.WriteString(as.Value.String())
	res.WriteString(";\n")

	return res.String()
}

type ReturnStatement struct {
	Token       token.Token
	Literal     string
	ReturnValue Expression
}

func (rs *ReturnStatement) statement()           {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Literal }
func (rs *ReturnStatement) String() string {
	var res strings.Builder

	res.WriteString(rs.Literal + " ")
	res.WriteString(rs.ReturnValue.String())
	res.WriteString(";\n")

	return res.String()
}

type Identifier struct {
	Token   token.Token
	Literal string
	Value   string
}

func (i *Identifier) expression()          {}
func (i *Identifier) TokenLiteral() string { return i.Literal }
func (i *Identifier) String() string       { return i.Value }

type IntLiteral struct {
	Token   token.Token
	Literal string
	Value   int64
}

func (il *IntLiteral) expression()          {}
func (il *IntLiteral) TokenLiteral() string { return il.Literal }
func (il *IntLiteral) String() string       { return il.Literal }

type CallExpression struct {
	Token    token.Token
	Literal  string
	Function *Identifier
}

func (ce *CallExpression) expression()          {}
func (ce *CallExpression) TokenLiteral() string { return ce.Literal }
func (ce *CallExpression) String() string {
	var res strings.Builder

	res.WriteString(ce.Function.String())
	res.WriteString("()")

	return res.String()
}
