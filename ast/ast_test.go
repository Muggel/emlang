package ast

import (
	"testing"

	"github.com/muggel/emlang/token"
	"github.com/stretchr/testify/assert"
)

func TestIntLiteral_String(t *testing.T) {
	intLit := &IntLiteral{Token: token.INT, Literal: "42", Value: 42}
	assert.Equal(t, "42", intLit.String())
}

func TestIdentifier_String(t *testing.T) {
	ident := &Identifier{Token: token.IDENT, Literal: "foo", Value: "foo"}
	assert.Equal(t, "foo", ident.String())
}

func TestReturnStatement_String(t *testing.T) {
	ident := &Identifier{Token: token.IDENT, Literal: "foo", Value: "foo"}
	retStmt := &ReturnStatement{Token: token.RETURN, Literal: "return", ReturnValue: ident}
	assert.Equal(t, "return foo;\n", retStmt.String())
}

func TestAssignmentStatement_String(t *testing.T) {
	ident := &Identifier{Token: token.IDENT, Literal: "foo", Value: "foo"}
	intLit := &IntLiteral{Token: token.INT, Literal: "42", Value: 42}
	assignStmt := &AssignmentStatement{
		Token: token.ASSIGN, Literal: "=", Identifier: ident, Value: intLit,
	}
	assert.Equal(t, "foo = 42;\n", assignStmt.String())
}

func TestBlockStatement_String(t *testing.T) {
	ident := &Identifier{Token: token.IDENT, Literal: "foo", Value: "foo"}
	intLit := &IntLiteral{Token: token.INT, Literal: "42", Value: 42}
	assignStmt := &AssignmentStatement{Token: token.ASSIGN, Literal: "=", Identifier: ident, Value: intLit}
	retStmt := &ReturnStatement{Token: token.RETURN, Literal: "return", ReturnValue: ident}
	blockStmt := &BlockStatement{Token: token.LBRACE, Literal: "{", Statements: []Statement{assignStmt, retStmt}}
	assert.Equal(t, "{\nfoo = 42;\nreturn foo;\n}\n", blockStmt.String())
}

func TestFunctionDeclaration_String(t *testing.T) {
	main := &Identifier{Token: token.IDENT, Literal: "main", Value: "main"}
	resValue := &Identifier{Token: token.IDENT, Literal: "int", Value: "int"}
	intLit := &IntLiteral{Token: token.INT, Literal: "42", Value: 42}
	retStmt := &ReturnStatement{Token: token.RETURN, Literal: "return", ReturnValue: intLit}
	blockStmt := &BlockStatement{Token: token.LBRACE, Literal: "{", Statements: []Statement{retStmt}}
	funcDecl := &FunctionDeclaration{
		Token: token.FN, Literal: "fn", Identifier: main, ReturnType: resValue, Body: blockStmt,
	}

	assert.Equal(t, "fn main() int {\nreturn 42;\n}\n", funcDecl.String())
}
