package testutil

import (
	"fmt"
	"monkey/ast"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func AssertLiteralExpression(t *testing.T, expected any, exp ast.Expression) {
	t.Helper()
	switch value := expected.(type) {
	case int:
		AssertIntegerLiteral(t, int64(value), exp)
	case int64:
		AssertIntegerLiteral(t, value, exp)
	case string:
		AssertIdentifier(t, value, exp)
	case bool:
		AssertBoolean(t, value, exp)
	default:
		t.Errorf("type of exp not handled. got=%T", exp)
	}
}

func AssertIntegerLiteral(t *testing.T, value int64, exp ast.Expression) {
	t.Helper()
	il, ok := exp.(*ast.IntegerLiteral)
	require.Truef(t, ok, "expected expression to be IntegerLiteral, got %T", exp)
	assert.Equal(t, value, il.Value)
	assert.Equal(t, fmt.Sprintf("%d", value), il.Token.Literal)
}

func AssertIdentifier(t *testing.T, value string, exp ast.Expression) {
	t.Helper()
	iden, ok := exp.(*ast.Identifier)
	require.Truef(t, ok, "expected expression to be Identifier, got %T", exp)
	assert.Equal(t, value, iden.Value)
	assert.Equal(t, value, iden.TokenLiteral())
}

func AssertBoolean(t *testing.T, value bool, exp ast.Expression) {
	t.Helper()
	b, ok := exp.(*ast.Boolean)
	require.Truef(t, ok, "expected expression to be Boolean, got %T", exp)
	assert.Equal(t, value, b.Value)
	assert.Equal(t, fmt.Sprintf("%t", value), b.TokenLiteral())
}

func AssertPrefixExpression(t *testing.T, exp ast.Expression, operator string, right any) {
	t.Helper()
	opExp, ok := exp.(*ast.PrefixExpression)
	require.Truef(t, ok, "expected expression to be PrefixExpression, got %T", exp)
	assert.Equal(t, operator, opExp.Operator)
	AssertLiteralExpression(t, right, opExp.Right)
}

func AssertInfixExpression(t *testing.T, exp ast.Expression, left any, operator string, right any) {
	t.Helper()
	opExp, ok := exp.(*ast.InfixExpression)
	require.Truef(t, ok, "expected expression to be InfixExpression, got %T", exp)
	AssertLiteralExpression(t, left, opExp.Left)
	assert.Equal(t, operator, opExp.Operator)
	AssertLiteralExpression(t, right, opExp.Right)
}
