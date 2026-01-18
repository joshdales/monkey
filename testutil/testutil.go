package testutil

import (
	"fmt"
	"monkey/ast"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func SetupProgram(t *testing.T, input string, stmtLen int) *ast.Program {
	t.Helper()
	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	require.NotNil(t, program, "program.ParseProgram() returned nil")
	if stmtLen > 0 {
		checkParserErrors(t, p)
		require.Len(t, program.Statements, stmtLen, "program.Statements does not contain %d statements.", stmtLen)
	}
	return program
}

func checkParserErrors(t *testing.T, p *parser.Parser) {
	t.Helper()
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestEval(t *testing.T, input string) object.Object {
	t.Helper()

	env := object.NewEnvironment()
	program := SetupProgram(t, input, 0)

	return evaluator.Eval(env, program)
}

func AssertExpressionStatement(t *testing.T, stmt ast.Statement) *ast.ExpressionStatement {
	t.Helper()
	exp, ok := stmt.(*ast.ExpressionStatement)
	require.Truef(t, ok, "expected statement to be ExpressionStatement, got %T", stmt)
	return exp
}

func AssertLetStatement(t *testing.T, name string, stmt ast.Statement) {
	t.Helper()
	assert.Equal(t, "let", stmt.TokenLiteral(), "TokenLiteral not 'let'")
	letStmt, ok := stmt.(*ast.LetStatement)
	require.Truef(t, ok, "expected statement to be LetStatement, got %T", stmt)
	assert.Equal(t, name, letStmt.Name.Value)
	assert.Equal(t, name, letStmt.Name.TokenLiteral())
}

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

func AssertIntegerObject(t *testing.T, expected int64, obj object.Object) {
	t.Helper()

	result, ok := obj.(*object.Integer)
	require.Truef(t, ok, "object is not an Integer, got %T (%+v)", obj, obj)
	assert.Equal(t, expected, result.Value)
}

func AssertBooleanObject(t *testing.T, expected bool, obj object.Object) {
	t.Helper()

	result, ok := obj.(*object.Boolean)
	require.Truef(t, ok, "object is not an Boolean, got %T (%+v)", obj, obj)
	assert.Equal(t, expected, result.Value)
}

func AssertNullObject(t *testing.T, obj object.Object) {
	t.Helper()

	assert.Equalf(t, evaluator.NULL, obj, "object is not NULL, got %T (%+v)", obj, obj)
}
