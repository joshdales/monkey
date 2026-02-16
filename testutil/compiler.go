package testutil

import (
	"monkey/code"
	"monkey/object"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func AssertInstructions(t *testing.T, actual code.Instructions, expected []code.Instructions) {
	t.Helper()
	concatted := concatInstructions(expected)

	require.Len(t, actual, len(concatted), "wrong number of instructions")

	for i, ins := range concatted {
		assert.EqualValues(t, ins, actual[i], "wrong instruction")
	}
}

func concatInstructions(instructions []code.Instructions) code.Instructions {
	out := code.Instructions{}
	for _, ins := range instructions {
		out = append(out, ins...)
	}

	return out
}

func AssertConstants(t *testing.T, actual []object.Object, expected []any) {
	t.Helper()
	require.Lenf(t, actual, len(expected), "wrong number of constants")

	for i, constant := range expected {
		AssertObject(t, actual[i], constant)
	}
}
