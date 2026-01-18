package testutil

import (
	"monkey/evaluator"
	"monkey/object"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
