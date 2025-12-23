package evaluator_test

import (
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		assertIntegerObject(t, tt.expected, evaluated)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		assertBooleanObject(t, tt.expected, evaluated)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!true", true},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		assertBooleanObject(t, tt.expected, evaluated)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			assertIntegerObject(t, int64(integer), evaluated)
		} else {
			assertNullObject(t, evaluated)
		}
	}
}

func TestReturnStatement(t *testing.T) {
	testCases := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 5 * 2; 9;", 10},
		{"9; return 2 * 5; 9", 10},
		{`
			if (10 > 1) {
				if (10 > 1) {
					return 10;
				}
			}

			return 1;
		`, 10,
		},
	}

	for _, tt := range testCases {
		evaluated := testEval(t, tt.input)
		assertIntegerObject(t, tt.expected, evaluated)
	}
}

func testEval(t *testing.T, input string) object.Object {
	t.Helper()

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return evaluator.Eval(program)
}

func assertIntegerObject(t *testing.T, expected int64, obj object.Object) {
	t.Helper()

	result, ok := obj.(*object.Integer)
	require.Truef(t, ok, "object is not an Integer, got %T (%+v)", obj, obj)
	assert.Equal(t, expected, result.Value)
}

func assertBooleanObject(t *testing.T, expected bool, obj object.Object) {
	t.Helper()

	result, ok := obj.(*object.Boolean)
	require.Truef(t, ok, "object is not an Boolean, got %T (%+v)", obj, obj)
	assert.Equal(t, expected, result.Value)
}

func assertNullObject(t *testing.T, obj object.Object) {
	t.Helper()

	assert.Equalf(t, evaluator.NULL, obj, "object is not NULL, got %T (%+v)", obj, obj)
}
