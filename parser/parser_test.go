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

	program := setupProgram(t, input, 3)
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
	stmt := testExpressionStatement(t, program.Statements[0])
	testLiteralExpression(t, "foobar", stmt.Expression)
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	program := setupProgram(t, input, 1)
	stmt := testExpressionStatement(t, program.Statements[0])
	testLiteralExpression(t, 5, stmt.Expression)
}

func TestBooleanExpression(t *testing.T) {
	input := "true;"
	program := setupProgram(t, input, 1)
	stmt := testExpressionStatement(t, program.Statements[0])
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
		program := setupProgram(t, tt.input, 1)
		stmt := testExpressionStatement(t, program.Statements[0])
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
		program := setupProgram(t, tt.input, 1)
		stmt := testExpressionStatement(t, program.Statements[0])
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
		program := setupProgram(t, tt.input, 0)
		actual := program.String()
		assert.Equal(t, tt.expected, actual)
	}
}

func TestIfExpression(t *testing.T) {
	input := "if (x < y) { x }"
	program := setupProgram(t, input, 1)
	stmt := testExpressionStatement(t, program.Statements[0])
	exp, ok := stmt.Expression.(*ast.IfExpression)
	require.Truef(t, ok, "expected IfStatement, got %T", stmt.Expression)
	testInfixExpression(t, exp.Condition, "x", "<", "y")
	assert.Len(t, exp.Consequence.Statements, 1)
	consequence := testExpressionStatement(t, exp.Consequence.Statements[0])
	testIdentifier(t, "x", consequence.Expression)
	require.Nil(t, exp.Alternative)
}

func TestIfElseExpression(t *testing.T) {
	input := "if (x < y) { x } else { y }"
	program := setupProgram(t, input, 1)
	stmt := testExpressionStatement(t, program.Statements[0])
	exp, ok := stmt.Expression.(*ast.IfExpression)
	require.Truef(t, ok, "expected expression to be IfExpression, got %T", exp)
	testInfixExpression(t, exp.Condition, "x", "<", "y")
	assert.Len(t, exp.Consequence.Statements, 1)
	consequence := testExpressionStatement(t, exp.Consequence.Statements[0])
	testIdentifier(t, "x", consequence.Expression)
	alternative := testExpressionStatement(t, exp.Alternative.Statements[0])
	testIdentifier(t, "y", alternative.Expression)
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := "fn(x, y) { x + y }"
	program := setupProgram(t, input, 1)
	stmt := testExpressionStatement(t, program.Statements[0])
	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	require.Truef(t, ok, "expected expression to be FunctionLiteral, got %T", stmt.Expression)
	assert.Len(t, function.Parameters, 2, "wrong number of function parameters")
	testLiteralExpression(t, "x", function.Parameters[0])
	testLiteralExpression(t, "y", function.Parameters[1])
	assert.Len(t, function.Body.Statements, 1, "wrong number of function body statements")
	bodyStmt := testExpressionStatement(t, function.Body.Statements[0])
	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestMainFunctionParameterParsing(t *testing.T) {
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
		stmt := testExpressionStatement(t, program.Statements[0])
		function := stmt.Expression.(*ast.FunctionLiteral)
		assert.Len(t, function.Parameters, len(tt.expectedParams), "length of parameters wrong")
		for idx, ident := range tt.expectedParams {
			testLiteralExpression(t, ident, function.Parameters[idx])
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"
	program := setupProgram(t, input, 1)
	stmt := testExpressionStatement(t, program.Statements[0])
	exp, ok := stmt.Expression.(*ast.CallExpression)
	require.Truef(t, ok, "expected expression to be CallExpression, got %T", stmt.Expression)
	testIdentifier(t, "add", exp.Function)
	assert.Len(t, exp.Arguments, 3, "length of arguments wrong")
	testLiteralExpression(t, 1, exp.Arguments[0])
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
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

func testExpressionStatement(t *testing.T, stmt ast.Statement) *ast.ExpressionStatement {
	t.Helper()
	exp, ok := stmt.(*ast.ExpressionStatement)
	require.Truef(t, ok, "expected statement to be ExpressionStatement, got %T", stmt)
	return exp
}

func testLetStatement(t *testing.T, name string, stmt ast.Statement) {
	t.Helper()
	assert.Equal(t, "let", stmt.TokenLiteral(), "TokenLiteral not 'let'")
	letStmt, ok := stmt.(*ast.LetStatement)
	require.Truef(t, ok, "expected statement to be LetStatement, got %T", stmt)
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
	require.Truef(t, ok, "expected expression to be IntegerLiteral, got %T", exp)
	assert.Equal(t, value, il.Value)
	assert.Equal(t, fmt.Sprintf("%d", value), il.Token.Literal)
}

func testIdentifier(t *testing.T, value string, exp ast.Expression) {
	t.Helper()
	iden, ok := exp.(*ast.Identifier)
	require.Truef(t, ok, "expected expression to be Identifier, got %T", exp)
	assert.Equal(t, value, iden.Value)
	assert.Equal(t, value, iden.TokenLiteral())
}

func testBoolean(t *testing.T, value bool, exp ast.Expression) {
	t.Helper()
	b, ok := exp.(*ast.Boolean)
	require.Truef(t, ok, "expected expression to be Boolean, got %T", exp)
	assert.Equal(t, value, b.Value)
	assert.Equal(t, fmt.Sprintf("%t", value), b.TokenLiteral())
}

func testPrefixExpression(t *testing.T, exp ast.Expression, operator string, right any) {
	t.Helper()
	opExp, ok := exp.(*ast.PrefixExpression)
	require.Truef(t, ok, "expected expression to be PrefixExpression, got %T", exp)
	assert.Equal(t, operator, opExp.Operator)
	testLiteralExpression(t, right, opExp.Right)
}

func testInfixExpression(t *testing.T, exp ast.Expression, left any, operator string, right any) {
	t.Helper()
	opExp, ok := exp.(*ast.InfixExpression)
	require.Truef(t, ok, "expected expression to be InfixExpression, got %T", exp)
	testLiteralExpression(t, left, opExp.Left)
	assert.Equal(t, operator, opExp.Operator)
	testLiteralExpression(t, right, opExp.Right)
}
