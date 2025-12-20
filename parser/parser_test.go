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
	testLiteralExpression(t, "foobar", stmt.Expression)
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	program := setupTests(t, input, 1)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)
	testLiteralExpression(t, 5, stmt.Expression)
}

func TestBooleanExpression(t *testing.T) {
	input := "true;"
	program := setupTests(t, input, 1)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)
	testLiteralExpression(t, true, stmt.Expression)
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    any
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range prefixTests {
		program := setupTests(t, tt.input, 1)
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		require.True(t, ok)
		testPrefixExpression(t, stmt.Expression, tt.operator, tt.value)
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  any
		operator   string
		rightValue any
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != true", true, "!=", true},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		program := setupTests(t, tt.input, 1)
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		require.True(t, ok)
		testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e -f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 > 4 != 3 < 4", "((5 > 4) != (3 < 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"!(true == true)", "(!(true == true))"},
	}

	for _, tt := range tests {
		program := setupTests(t, tt.input, 0)
		actual := program.String()
		assert.Equal(t, tt.expected, actual)
	}
}

// Test Helpers

func setupTests(t *testing.T, input string, stmtLen int) *ast.Program {
	t.Helper()
	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	require.NotNil(t, program, "program.ParseProgram() returned nil")
	checkParserErrors(t, p)
	if stmtLen > 0 {
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

func testLetStatement(t *testing.T, name string, s ast.Statement) {
	t.Helper()
	assert.Equal(t, "let", s.TokenLiteral(), "TokenLiteral not 'let'")
	letStmt, ok := s.(*ast.LetStatement)
	require.True(t, ok)
	assert.Equal(t, name, letStmt.Name.Value)
	assert.Equal(t, name, letStmt.Name.TokenLiteral())
}

func testLiteralExpression(t *testing.T, expected interface{}, exp ast.Expression) {
	t.Helper()
	switch value := expected.(type) {
	case int:
		testIntegerLiteral(t, int64(value), exp)
	case int64:
		testIntegerLiteral(t, value, exp)
	case string:
		testIdentifier(t, value, exp)
	case bool:
		testBoolean(t, value, exp)
	default:
		t.Errorf("type of exp not handled. got=%T", exp)
	}
}

func testIntegerLiteral(t *testing.T, value int64, exp ast.Expression) {
	t.Helper()
	il, ok := exp.(*ast.IntegerLiteral)
	require.True(t, ok)
	assert.Equal(t, value, il.Value)
	assert.Equal(t, fmt.Sprintf("%d", value), il.Token.Literal)
}

func testIdentifier(t *testing.T, value string, exp ast.Expression) {
	t.Helper()
	iden, ok := exp.(*ast.Identifier)
	require.True(t, ok)
	assert.Equal(t, value, iden.Value)
	assert.Equal(t, value, iden.TokenLiteral())
}

func testBoolean(t *testing.T, value bool, exp ast.Expression) {
	t.Helper()
	b, ok := exp.(*ast.Boolean)
	require.True(t, ok)
	assert.Equal(t, value, b.Value)
	assert.Equal(t, fmt.Sprintf("%t", value), b.TokenLiteral())
}

func testPrefixExpression(t *testing.T, exp ast.Expression, operator string, right any) {
	t.Helper()
	opExp, ok := exp.(*ast.PrefixExpression)
	require.True(t, ok)
	assert.Equal(t, operator, opExp.Operator)
	testLiteralExpression(t, right, opExp.Right)
}

func testInfixExpression(t *testing.T, exp ast.Expression, left any, operator string, right any) {
	t.Helper()
	opExp, ok := exp.(*ast.InfixExpression)
	require.True(t, ok)
	testLiteralExpression(t, left, opExp.Left)
	assert.Equal(t, operator, opExp.Operator)
	testLiteralExpression(t, right, opExp.Right)
}
