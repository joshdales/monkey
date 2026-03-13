package evaluator_test

import (
	"monkey/object"
	"monkey/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuiltinFunctions(t *testing.T) {
	t.Run("len", func(t *testing.T) {
		tests := []struct {
			input    string
			expected any
		}{
			{`len("")`, 0},
			{`len("four")`, 4},
			{`len("hello world")`, 11},
			{`len(1)`, "argument to `len` not supported, got=INTEGER"},
			{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
		}

		for _, tt := range tests {
			evaluated := testutil.TestEval(t, tt.input)
			switch expected := tt.expected.(type) {
			case int:
				testutil.AssertIntegerObject(t, evaluated, int64(expected))
			case string:
				errObj, ok := evaluated.(*object.Error)
				require.Truef(t, ok, "object is not an Error, got %T (%+v)", evaluated, evaluated)
				assert.Equal(t, expected, errObj.Message)
			}
		}
	})
}
