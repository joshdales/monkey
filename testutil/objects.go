package testutil

import (
	"monkey/evaluator"
	"monkey/object"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
