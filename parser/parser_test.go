package parser_test

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/parser"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTests(t *testing.T, input string, stmtLen int) *ast.Program {
	t.Helper()
	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	require.NotNil(t, program, "program.ParseProgram() returned nil")
	checkParserErrors(t, p)
	require.Len(t, program.Statements, stmtLen, "program.Statements does not contain %d statements.", stmtLen)
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

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`

	program := setupTests(t, input, 3)
	tests := []struct{ expectedIdentifier string }{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		testLetStatement(t, tt.expectedIdentifier, stmt)
	}
}

func testLetStatement(t *testing.T, name string, s ast.Statement) {
	t.Helper()
	assert.Equal(t, "let", s.TokenLiteral(), "TokenLiteral not 'let'")
	letStmt, ok := s.(*ast.LetStatement)
	require.True(t, ok)
	assert.Equal(t, name, letStmt.Name.Value)
	assert.Equal(t, name, letStmt.Name.TokenLiteral())
}

func TestReturnStatement(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`
	program := setupTests(t, input, 3)
	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		require.True(t, ok)
		assert.Equal(t, "return", returnStmt.TokenLiteral(), "TokenLiteral not 'return'")
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	program := setupTests(t, input, 1)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)
	ident, ok := stmt.Expression.(*ast.Identifier)
	require.True(t, ok)
	assert.Equal(t, "foobar", ident.Value)
	assert.Equal(t, "foobar", ident.Token.Literal)
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	program := setupTests(t, input, 1)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	require.True(t, ok)
	assert.Equal(t, int64(5), literal.Value)
	assert.Equal(t, "5", literal.TokenLiteral())
}

func testIntegerLiteral(t *testing.T, value int64, exp ast.Expression) {
	t.Helper()
	il, ok := exp.(*ast.IntegerLiteral)
	require.True(t, ok)
	assert.Equal(t, value, il.Value)
	assert.Equal(t, fmt.Sprintf("%d", value), il.Token.Literal)
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range prefixTests {
		program := setupTests(t, tt.input, 1)
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		require.True(t, ok)
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		require.True(t, ok)
		assert.Equal(t, tt.operator, exp.Operator)
		testIntegerLiteral(t, tt.integerValue, exp.Right)
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		program := setupTests(t, tt.input, 1)
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		require.True(t, ok)
		exp, ok := stmt.Expression.(*ast.InfixExpression)
		require.True(t, ok)
		testIntegerLiteral(t, tt.leftValue, exp.Left)
		assert.Equal(t, tt.operator, exp.Operator)
		testIntegerLiteral(t, tt.rightValue, exp.Right)
	}
}
