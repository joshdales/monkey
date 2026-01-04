package evaluator_test

import (
	"monkey/evaluator"
	"monkey/object"
	"monkey/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefineMacros(t *testing.T) {
	input := `
	let number = 1
	let function = fn(x, y) {x + y};
	let mymacro = macro(x, y) {x + y};
	`
	env := object.NewEnvironment()
	program := testutil.SetupProgram(t, input, 2)

	evaluator.DefineMacros(env, program)

	_, ok := env.Get("number")
	assert.Falsef(t, ok, "number should not be defined")
	_, ok = env.Get("function")
	assert.Falsef(t, ok, "function should not be defined")
	obj, ok := env.Get("mymacro")
	assert.Truef(t, ok, "mymacro not in environment")
	macro, ok := obj.(*object.Macro)
	require.Truef(t, ok, "object is not a Macro, got %T (%+v)", obj, obj)
	assert.Len(t, macro.Parameters, 2)
	assert.Equal(t, "x", macro.Parameters[0].String())
	assert.Equal(t, "y", macro.Parameters[1].String())
	assert.Equal(t, "(x + y)", macro.Body.String())
}
