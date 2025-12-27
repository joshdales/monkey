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
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      any
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		program := setupProgram(t, tt.input, 1)
		stmt := program.Statements[0]
		assertLetStatement(t, tt.expectedIdentifier, stmt)
		val := stmt.(*ast.LetStatement).Value
		assertLiteralExpression(t, tt.expectedValue, val)
	}
}

func TestReturnStatement(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`
	program := setupProgram(t, input, 3)
	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		require.Truef(t, ok, "expected statement to be ReturnStatement, got %T", stmt)
		assert.Equal(t, "return", returnStmt.TokenLiteral(), "TokenLiteral not 'return'")
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	program := setupProgram(t, input, 1)
	stmt := assertExpressionStatement(t, program.Statements[0])
	assertLiteralExpression(t, "foobar", stmt.Expression)
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	program := setupProgram(t, input, 1)
	stmt := assertExpressionStatement(t, program.Statements[0])
	assertLiteralExpression(t, 5, stmt.Expression)
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`
	program := setupProgram(t, input, 1)
	stmt := assertExpressionStatement(t, program.Statements[0])
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	require.True(t, ok, "expected expression to be StringLiteral, got %t", stmt)
	assert.Equal(t, "hello world", literal.Value)
}

func TestBooleanExpression(t *testing.T) {
	input := "true;"
	program := setupProgram(t, input, 1)
	stmt := assertExpressionStatement(t, program.Statements[0])
	assertLiteralExpression(t, true, stmt.Expression)
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
		program := setupProgram(t, tt.input, 1)
		stmt := assertExpressionStatement(t, program.Statements[0])
		assertPrefixExpression(t, stmt.Expression, tt.operator, tt.value)
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
		program := setupProgram(t, tt.input, 1)
		stmt := assertExpressionStatement(t, program.Statements[0])
		assertInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue)
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
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},
	}

	for _, tt := range tests {
		program := setupProgram(t, tt.input, 0)
		actual := program.String()
		assert.Equal(t, tt.expected, actual)
	}
}

func TestIfExpression(t *testing.T) {
	input := "if (x < y) { x }"
	program := setupProgram(t, input, 1)
	stmt := assertExpressionStatement(t, program.Statements[0])
	exp, ok := stmt.Expression.(*ast.IfExpression)
	require.Truef(t, ok, "expected IfStatement, got %T", stmt.Expression)
	assertInfixExpression(t, exp.Condition, "x", "<", "y")
	assert.Len(t, exp.Consequence.Statements, 1)
	consequence := assertExpressionStatement(t, exp.Consequence.Statements[0])
	assertIdentifier(t, "x", consequence.Expression)
	require.Nil(t, exp.Alternative)
}

func TestIfElseExpression(t *testing.T) {
	input := "if (x < y) { x } else { y }"
	program := setupProgram(t, input, 1)
	stmt := assertExpressionStatement(t, program.Statements[0])
	exp, ok := stmt.Expression.(*ast.IfExpression)
	require.Truef(t, ok, "expected expression to be IfExpression, got %T", exp)
	assertInfixExpression(t, exp.Condition, "x", "<", "y")
	assert.Len(t, exp.Consequence.Statements, 1)
	consequence := assertExpressionStatement(t, exp.Consequence.Statements[0])
	assertIdentifier(t, "x", consequence.Expression)
	alternative := assertExpressionStatement(t, exp.Alternative.Statements[0])
	assertIdentifier(t, "y", alternative.Expression)
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := "fn(x, y) { x + y }"
	program := setupProgram(t, input, 1)
	stmt := assertExpressionStatement(t, program.Statements[0])
	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	require.Truef(t, ok, "expected expression to be FunctionLiteral, got %T", stmt.Expression)
	assert.Len(t, function.Parameters, 2, "wrong number of function parameters")
	assertLiteralExpression(t, "x", function.Parameters[0])
	assertLiteralExpression(t, "y", function.Parameters[1])
	assert.Len(t, function.Body.Statements, 1, "wrong number of function body statements")
	bodyStmt := assertExpressionStatement(t, function.Body.Statements[0])
	assertInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {}", expectedParams: []string{}},
		{input: "fn(x) {}", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {}", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		program := setupProgram(t, tt.input, 0)
		stmt := assertExpressionStatement(t, program.Statements[0])
		function := stmt.Expression.(*ast.FunctionLiteral)
		assert.Len(t, function.Parameters, len(tt.expectedParams), "length of parameters wrong")
		for idx, ident := range tt.expectedParams {
			assertLiteralExpression(t, ident, function.Parameters[idx])
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"
	program := setupProgram(t, input, 1)
	stmt := assertExpressionStatement(t, program.Statements[0])
	exp, ok := stmt.Expression.(*ast.CallExpression)
	require.Truef(t, ok, "expected expression to be CallExpression, got %T", stmt.Expression)
	assertIdentifier(t, "add", exp.Function)
	assert.Len(t, exp.Arguments, 3, "length of arguments wrong")
	assertLiteralExpression(t, 1, exp.Arguments[0])
	assertInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	assertInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestCallExpressionParameterParsing(t *testing.T) {
	tests := []struct {
		input         string
		expectedIdent string
		expectedArgs  []string
	}{
		{input: "add();", expectedIdent: "add", expectedArgs: []string{}},
		{input: "add(1);", expectedIdent: "add", expectedArgs: []string{"1"}},
		{input: "add(1, 2 * 3, 4 + 5);", expectedIdent: "add", expectedArgs: []string{"1", "(2 * 3)", "(4 + 5)"}},
	}

	for _, tt := range tests {
		program := setupProgram(t, tt.input, 0)
		stmt := assertExpressionStatement(t, program.Statements[0])

		exp, ok := stmt.Expression.(*ast.CallExpression)
		require.Truef(t, ok, "expected expression to be CallExpression, got %T", stmt.Expression)
		assert.Len(t, exp.Arguments, len(tt.expectedArgs), "length of arguments wrong")

		for idx, arg := range tt.expectedArgs {
			assert.Equal(t, arg, exp.Arguments[idx].String())
		}
	}
}

// Test Helpers

func setupProgram(t *testing.T, input string, stmtLen int) *ast.Program {
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

func assertExpressionStatement(t *testing.T, stmt ast.Statement) *ast.ExpressionStatement {
	t.Helper()
	exp, ok := stmt.(*ast.ExpressionStatement)
	require.Truef(t, ok, "expected statement to be ExpressionStatement, got %T", stmt)
	return exp
}

func assertLetStatement(t *testing.T, name string, stmt ast.Statement) {
	t.Helper()
	assert.Equal(t, "let", stmt.TokenLiteral(), "TokenLiteral not 'let'")
	letStmt, ok := stmt.(*ast.LetStatement)
	require.Truef(t, ok, "expected statement to be LetStatement, got %T", stmt)
	assert.Equal(t, name, letStmt.Name.Value)
	assert.Equal(t, name, letStmt.Name.TokenLiteral())
}

func assertLiteralExpression(t *testing.T, expected any, exp ast.Expression) {
	t.Helper()
	switch value := expected.(type) {
	case int:
		assertIntegerLiteral(t, int64(value), exp)
	case int64:
		assertIntegerLiteral(t, value, exp)
	case string:
		assertIdentifier(t, value, exp)
	case bool:
		assertBoolean(t, value, exp)
	default:
		t.Errorf("type of exp not handled. got=%T", exp)
	}
}

func assertIntegerLiteral(t *testing.T, value int64, exp ast.Expression) {
	t.Helper()
	il, ok := exp.(*ast.IntegerLiteral)
	require.Truef(t, ok, "expected expression to be IntegerLiteral, got %T", exp)
	assert.Equal(t, value, il.Value)
	assert.Equal(t, fmt.Sprintf("%d", value), il.Token.Literal)
}

func assertIdentifier(t *testing.T, value string, exp ast.Expression) {
	t.Helper()
	iden, ok := exp.(*ast.Identifier)
	require.Truef(t, ok, "expected expression to be Identifier, got %T", exp)
	assert.Equal(t, value, iden.Value)
	assert.Equal(t, value, iden.TokenLiteral())
}

func assertBoolean(t *testing.T, value bool, exp ast.Expression) {
	t.Helper()
	b, ok := exp.(*ast.Boolean)
	require.Truef(t, ok, "expected expression to be Boolean, got %T", exp)
	assert.Equal(t, value, b.Value)
	assert.Equal(t, fmt.Sprintf("%t", value), b.TokenLiteral())
}

func assertPrefixExpression(t *testing.T, exp ast.Expression, operator string, right any) {
	t.Helper()
	opExp, ok := exp.(*ast.PrefixExpression)
	require.Truef(t, ok, "expected expression to be PrefixExpression, got %T", exp)
	assert.Equal(t, operator, opExp.Operator)
	assertLiteralExpression(t, right, opExp.Right)
}

func assertInfixExpression(t *testing.T, exp ast.Expression, left any, operator string, right any) {
	t.Helper()
	opExp, ok := exp.(*ast.InfixExpression)
	require.Truef(t, ok, "expected expression to be InfixExpression, got %T", exp)
	assertLiteralExpression(t, left, opExp.Left)
	assert.Equal(t, operator, opExp.Operator)
	assertLiteralExpression(t, right, opExp.Right)
}
