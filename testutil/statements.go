package testutil

import (
	"monkey/ast"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func AssertExpressionStatement(t *testing.T, stmt ast.Statement) *ast.ExpressionStatement {
	t.Helper()
	exp, ok := stmt.(*ast.ExpressionStatement)
	require.Truef(t, ok, "expected statement to be ExpressionStatement, got %T", stmt)
	return exp
}

func AssertLetStatement(t *testing.T, stmt ast.Statement, expectedName string) {
	t.Helper()
	assert.Equal(t, "let", stmt.TokenLiteral(), "TokenLiteral not 'let'")
	letStmt, ok := stmt.(*ast.LetStatement)
	require.Truef(t, ok, "expected statement to be LetStatement, got %T", stmt)
	assert.Equal(t, expectedName, letStmt.Name.Value)
	assert.Equal(t, expectedName, letStmt.Name.TokenLiteral())
}
