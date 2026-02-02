package parser_test

import (
	"monkey/ast"
	"monkey/testutil"
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
		program := testutil.SetupProgram(t, tt.input, 1)
		stmt := program.Statements[0]
		testutil.AssertLetStatement(t, stmt, tt.expectedIdentifier)
		val := stmt.(*ast.LetStatement).Value
		testutil.AssertLiteralExpression(t, val, tt.expectedValue)
	}
}

func TestReturnStatement(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`
	program := testutil.SetupProgram(t, input, 3)
	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		require.Truef(t, ok, "expected statement to be ReturnStatement, got %T", stmt)
		assert.Equal(t, "return", returnStmt.TokenLiteral(), "TokenLiteral not 'return'")
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	program := testutil.SetupProgram(t, input, 1)
	stmt := testutil.AssertExpressionStatement(t, program.Statements[0])
	testutil.AssertLiteralExpression(t, stmt.Expression, "foobar")
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	program := testutil.SetupProgram(t, input, 1)
	stmt := testutil.AssertExpressionStatement(t, program.Statements[0])
	testutil.AssertLiteralExpression(t, stmt.Expression, 5)
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`
	program := testutil.SetupProgram(t, input, 1)
	stmt := testutil.AssertExpressionStatement(t, program.Statements[0])
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	require.True(t, ok, "expected expression to be StringLiteral, got %t", stmt)
	assert.Equal(t, "hello world", literal.Value)
}

func TestBooleanExpression(t *testing.T) {
	input := "true;"
	program := testutil.SetupProgram(t, input, 1)
	stmt := testutil.AssertExpressionStatement(t, program.Statements[0])
	testutil.AssertLiteralExpression(t, stmt.Expression, true)
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
		program := testutil.SetupProgram(t, tt.input, 1)
		stmt := testutil.AssertExpressionStatement(t, program.Statements[0])
		testutil.AssertPrefixExpression(t, stmt.Expression, tt.operator, tt.value)
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
		program := testutil.SetupProgram(t, tt.input, 1)
		stmt := testutil.AssertExpressionStatement(t, program.Statements[0])
		testutil.AssertInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue)
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
		{"a * [1, 2, 3, 4][b *  c] * d", "((a * ([1, 2, 3, 4][(b * c)])) * d)"},
		{"add(a * b[2], b[1], 2 * [1, 2][1])", "add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))"},
	}

	for _, tt := range tests {
		program := testutil.SetupProgram(t, tt.input, 0)
		actual := program.String()
		assert.Equal(t, tt.expected, actual)
	}
}

func TestIfExpression(t *testing.T) {
	input := "if (x < y) { x }"
	program := testutil.SetupProgram(t, input, 1)
	stmt := testutil.AssertExpressionStatement(t, program.Statements[0])
	exp, ok := stmt.Expression.(*ast.IfExpression)
	require.Truef(t, ok, "expected IfStatement, got %T", stmt.Expression)
	testutil.AssertInfixExpression(t, exp.Condition, "x", "<", "y")
	assert.Len(t, exp.Consequence.Statements, 1)
	consequence := testutil.AssertExpressionStatement(t, exp.Consequence.Statements[0])
	testutil.AssertIdentifier(t, consequence.Expression, "x")
	require.Nil(t, exp.Alternative)
}

func TestIfElseExpression(t *testing.T) {
	input := "if (x < y) { x } else { y }"
	program := testutil.SetupProgram(t, input, 1)
	stmt := testutil.AssertExpressionStatement(t, program.Statements[0])
	exp, ok := stmt.Expression.(*ast.IfExpression)
	require.Truef(t, ok, "expected expression to be IfExpression, got %T", exp)
	testutil.AssertInfixExpression(t, exp.Condition, "x", "<", "y")
	assert.Len(t, exp.Consequence.Statements, 1)
	consequence := testutil.AssertExpressionStatement(t, exp.Consequence.Statements[0])
	testutil.AssertIdentifier(t, consequence.Expression, "x")
	alternative := testutil.AssertExpressionStatement(t, exp.Alternative.Statements[0])
	testutil.AssertIdentifier(t, alternative.Expression, "y")
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := "fn(x, y) { x + y }"
	program := testutil.SetupProgram(t, input, 1)
	stmt := testutil.AssertExpressionStatement(t, program.Statements[0])
	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	require.Truef(t, ok, "expected expression to be FunctionLiteral, got %T", stmt.Expression)
	assert.Len(t, function.Parameters, 2, "wrong number of function parameters")
	testutil.AssertLiteralExpression(t, function.Parameters[0], "x")
	testutil.AssertLiteralExpression(t, function.Parameters[1], "y")
	assert.Len(t, function.Body.Statements, 1, "wrong number of function body statements")
	bodyStmt := testutil.AssertExpressionStatement(t, function.Body.Statements[0])
	testutil.AssertInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
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
		program := testutil.SetupProgram(t, tt.input, 0)
		stmt := testutil.AssertExpressionStatement(t, program.Statements[0])
		function := stmt.Expression.(*ast.FunctionLiteral)
		assert.Len(t, function.Parameters, len(tt.expectedParams), "length of parameters wrong")
		for idx, ident := range tt.expectedParams {
			testutil.AssertLiteralExpression(t, function.Parameters[idx], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"
	program := testutil.SetupProgram(t, input, 1)
	stmt := testutil.AssertExpressionStatement(t, program.Statements[0])
	exp, ok := stmt.Expression.(*ast.CallExpression)
	require.Truef(t, ok, "expected expression to be CallExpression, got %T", stmt.Expression)
	testutil.AssertIdentifier(t, exp.Function, "add")
	assert.Len(t, exp.Arguments, 3, "length of arguments wrong")
	testutil.AssertLiteralExpression(t, exp.Arguments[0], 1)
	testutil.AssertInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testutil.AssertInfixExpression(t, exp.Arguments[2], 4, "+", 5)
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
		program := testutil.SetupProgram(t, tt.input, 0)
		stmt := testutil.AssertExpressionStatement(t, program.Statements[0])

		exp, ok := stmt.Expression.(*ast.CallExpression)
		require.Truef(t, ok, "expected expression to be CallExpression, got %T", stmt.Expression)
		assert.Len(t, exp.Arguments, len(tt.expectedArgs), "length of arguments wrong")

		for idx, arg := range tt.expectedArgs {
			assert.Equal(t, arg, exp.Arguments[idx].String())
		}
	}
}

func TestArrayParsing(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	program := testutil.SetupProgram(t, input, 1)
	stmt := testutil.AssertExpressionStatement(t, program.Statements[0])
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	require.Truef(t, ok, "expected expression to be ArrayLiteral, got %T", stmt.Expression)
	testutil.AssertIntegerLiteral(t, array.Elements[0], 1)
	testutil.AssertInfixExpression(t, array.Elements[1], 2, "*", 2)
	testutil.AssertInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestHashLiteralParsing(t *testing.T) {
	t.Run("with string keys", func(t *testing.T) {
		input := `{"one": 1, "two": 2, "three": 3}`
		program := testutil.SetupProgram(t, input, 0)
		stmt := testutil.AssertExpressionStatement(t, program.Statements[0])
		hash, ok := stmt.Expression.(*ast.HashLiteral)
		require.Truef(t, ok, "expected expression to be HashLiteral, got %T", stmt.Expression)
		require.Len(t, hash.Pairs, 3)

		expected := map[string]int64{
			"one":   1,
			"two":   2,
			"three": 3,
		}

		for key, value := range hash.Pairs {
			literal, ok := key.(*ast.StringLiteral)
			require.Truef(t, ok, "key is not a string literal, got %T", key)
			expectedValue := expected[literal.String()]
			testutil.AssertIntegerLiteral(t, value, expectedValue)
		}
	})

	t.Run("with integer keys", func(t *testing.T) {
		input := `{1: "one", 2: "two", 3: "three"}`
		program := testutil.SetupProgram(t, input, 0)
		stmt := testutil.AssertExpressionStatement(t, program.Statements[0])
		hash, ok := stmt.Expression.(*ast.HashLiteral)
		require.Truef(t, ok, "expected expression to be HashLiteral, got %T", stmt.Expression)
		require.Len(t, hash.Pairs, 3)

		expected := map[int64]string{
			1: "one",
			2: "two",
			3: "three",
		}

		for key, value := range hash.Pairs {
			literal, ok := key.(*ast.IntegerLiteral)
			require.Truef(t, ok, "key is not a integer literal, got %T", key)
			expectedValue := expected[literal.Value]
			assert.Equal(t, expectedValue, value.String())
		}
	})

	t.Run("with boolean keys", func(t *testing.T) {
		input := `{true: 1, false: 0}`
		program := testutil.SetupProgram(t, input, 0)
		stmt := testutil.AssertExpressionStatement(t, program.Statements[0])
		hash, ok := stmt.Expression.(*ast.HashLiteral)
		require.Truef(t, ok, "expected expression to be HashLiteral, got %T", stmt.Expression)
		require.Len(t, hash.Pairs, 2)

		expected := map[bool]int64{
			true:  1,
			false: 0,
		}

		for key, value := range hash.Pairs {
			literal, ok := key.(*ast.Boolean)
			require.Truef(t, ok, "key is not a boolean, got %T", key)
			expectedValue := expected[literal.Value]
			testutil.AssertIntegerLiteral(t, value, expectedValue)
		}
	})

	t.Run("with expression values", func(t *testing.T) {
		input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`
		program := testutil.SetupProgram(t, input, 0)
		stmt := testutil.AssertExpressionStatement(t, program.Statements[0])
		hash, ok := stmt.Expression.(*ast.HashLiteral)
		require.Truef(t, ok, "expected expression to be HashLiteral, got %T", stmt.Expression)
		require.Len(t, hash.Pairs, 3)

		expected := map[string]func(ast.Expression){
			"one": func(e ast.Expression) {
				testutil.AssertInfixExpression(t, e, 0, "+", 1)
			},
			"two": func(e ast.Expression) {
				testutil.AssertInfixExpression(t, e, 10, "-", 8)
			},
			"three": func(e ast.Expression) {
				testutil.AssertInfixExpression(t, e, 15, "/", 5)
			},
		}

		for key, value := range hash.Pairs {
			literal, ok := key.(*ast.StringLiteral)
			require.Truef(t, ok, "key is not a string literal, got %T", key)
			testFunc := expected[literal.String()]
			testFunc(value)
		}
	})

	t.Run("with an empty hash", func(t *testing.T) {
		input := `{}`
		program := testutil.SetupProgram(t, input, 0)
		stmt := testutil.AssertExpressionStatement(t, program.Statements[0])
		hash, ok := stmt.Expression.(*ast.HashLiteral)
		require.Truef(t, ok, "expected expression to be HashLiteral, got %T", stmt.Expression)
		assert.Len(t, hash.Pairs, 0)
	})
}

func TestIndexExpressionParsing(t *testing.T) {
	input := "myArray[1 + 1]"
	program := testutil.SetupProgram(t, input, 0)
	stmt := testutil.AssertExpressionStatement(t, program.Statements[0])
	idxExp, ok := stmt.Expression.(*ast.IndexExpression)
	require.Truef(t, ok, "expected expression to be IndexExpression, got %T", stmt.Expression)
	testutil.AssertIdentifier(t, idxExp.Left, "myArray")
	testutil.AssertInfixExpression(t, idxExp.Index, 1, "+", 1)
}

func TestMacroLiteralParsing(t *testing.T) {
	input := `macro(x, y) { x + y }`
	program := testutil.SetupProgram(t, input, 1)
	stmt := testutil.AssertExpressionStatement(t, program.Statements[0])
	macro, ok := stmt.Expression.(*ast.MacroLiteral)
	require.Truef(t, ok, "statement is not ast.MacroLiteral, got %T", stmt.Expression)
	assert.Len(t, macro.Parameters, 2)
	testutil.AssertLiteralExpression(t, macro.Parameters[0], "x")
	testutil.AssertLiteralExpression(t, macro.Parameters[1], "y")
	assert.Len(t, macro.Body.Statements, 1)
	bodyStmt := testutil.AssertExpressionStatement(t, macro.Body.Statements[0])
	testutil.AssertInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}
