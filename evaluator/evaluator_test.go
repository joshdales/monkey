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
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		assertBooleanObject(t, tt.expected, evaluated)
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
