package testutil

import (
	"fmt"
	"monkey/ast"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func AssertLiteralExpression(t *testing.T, actual ast.Expression, expected any) {
	t.Helper()
	switch value := expected.(type) {
	case int:
		AssertIntegerLiteral(t, actual, int64(value))
	case int64:
		AssertIntegerLiteral(t, actual, value)
	case string:
		AssertIdentifier(t, actual, value)
	case bool:
		AssertBoolean(t, actual, value)
	default:
		t.Errorf("type of exp not handled. got=%T", actual)
	}
}

func AssertIntegerLiteral(t *testing.T, actual ast.Expression, expected int64) {
	t.Helper()
	il, ok := actual.(*ast.IntegerLiteral)
	require.Truef(t, ok, "expected expression to be IntegerLiteral, got %T", actual)
	assert.Equal(t, expected, il.Value)
	assert.Equal(t, fmt.Sprintf("%d", expected), il.Token.Literal)
}

func AssertIdentifier(t *testing.T, actual ast.Expression, expected string) {
	t.Helper()
	iden, ok := actual.(*ast.Identifier)
	require.Truef(t, ok, "expected expression to be Identifier, got %T", actual)
	assert.Equal(t, expected, iden.Value)
	assert.Equal(t, expected, iden.TokenLiteral())
}

func AssertBoolean(t *testing.T, actual ast.Expression, expected bool) {
	t.Helper()
	b, ok := actual.(*ast.Boolean)
	require.Truef(t, ok, "expected expression to be Boolean, got %T", actual)
	assert.Equal(t, expected, b.Value)
	assert.Equal(t, fmt.Sprintf("%t", expected), b.TokenLiteral())
}

func AssertPrefixExpression(t *testing.T, actual ast.Expression, expectedOperator string, expectedRight any) {
	t.Helper()
	opExp, ok := actual.(*ast.PrefixExpression)
	require.Truef(t, ok, "expected expression to be PrefixExpression, got %T", actual)
	assert.Equal(t, expectedOperator, opExp.Operator)
	AssertLiteralExpression(t, opExp.Right, expectedRight)
}

func AssertInfixExpression(t *testing.T, actual ast.Expression, expectedLeft any, expectedOperator string, expectedRight any) {
	t.Helper()
	opExp, ok := actual.(*ast.InfixExpression)
	require.Truef(t, ok, "expected expression to be InfixExpression, got %T", actual)
	AssertLiteralExpression(t, opExp.Left, expectedLeft)
	assert.Equal(t, expectedOperator, opExp.Operator)
	AssertLiteralExpression(t, opExp.Right, expectedRight)
}
