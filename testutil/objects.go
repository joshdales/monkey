package testutil

import (
	"monkey/evaluator"
	"monkey/object"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func AssertObject(t *testing.T, actual object.Object, expected any) {
	t.Helper()

	switch value := expected.(type) {
	case int:
		AssertIntegerObject(t, actual, int64(value))
	case []int:
		AssertIntegerArray(t, actual, value)
	case int64:
		AssertIntegerObject(t, actual, value)
	case bool:
		AssertBooleanObject(t, actual, value)
	case string:
		AssertStringObject(t, actual, value)
	case nil, *object.Null:
		AssertNullObject(t, actual)
	}
}

func AssertIntegerObject(t *testing.T, actual object.Object, expected int64) {
	t.Helper()

	result, ok := actual.(*object.Integer)
	require.Truef(t, ok, "object is not an Integer, got %T (%+v)", actual, actual)
	assert.Equal(t, expected, result.Value)
}

func AssertBooleanObject(t *testing.T, actual object.Object, expected bool) {
	t.Helper()

	result, ok := actual.(*object.Boolean)
	require.Truef(t, ok, "object is not an Boolean, got %T (%+v)", actual, actual)
	assert.Equal(t, expected, result.Value)
}

func AssertNullObject(t *testing.T, obj object.Object) {
	t.Helper()

	assert.Equalf(t, evaluator.NULL, obj, "object is not NULL, got %T (%+v)", obj, obj)
}

func AssertStringObject(t *testing.T, actual object.Object, expected string) {
	t.Helper()

	result, ok := actual.(*object.String)
	require.Truef(t, ok, "object is not a String, got %T (%+v)", actual)
	assert.Equal(t, expected, result.Value)
}

func AssertIntegerArray(t *testing.T, actual object.Object, expected []int) {
	t.Helper()

	array, ok := actual.(*object.Array)
	require.Truef(t, ok, "object us not an Array, got %T, (%+v)", actual, actual)
	assert.Len(t, array.Elements, len(expected))
	for i, expectedElem := range expected {
		AssertIntegerObject(t, array.Elements[i], int64(expectedElem))
	}

}
