package compiler_test

import (
	"monkey/code"
	"monkey/compiler"
	"monkey/object"
	"monkey/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type compilerTestCase struct {
	input                string
	expectedConstants    []any
	expectedInstructions []code.Instructions
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: "1 + 2", expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
			},
		},
	}

	runCompilerTests(t, tests)
}

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()

	for _, tt := range tests {

		program := testutil.SetupProgram(t, tt.input, 0)
		compiler := compiler.New()
		err := compiler.Compile(program)
		require.NoError(t, err)

		bytecode := compiler.Bytecode()
		err = testInstructions(t, tt.expectedInstructions, bytecode.Instructions)
		require.NoError(t, err)

		err = testConstants(t, tt.expectedConstants, bytecode.Constants)
		require.NoError(t, err)

	}
}

func testInstructions(t *testing.T, expected []code.Instructions, actual code.Instructions) error {
	t.Helper()
	concatted := concatInstructions(expected)

	require.Len(t, actual, len(concatted))

	for i, ins := range concatted {
		assert.EqualValues(t, ins, actual[i])
	}
	return nil
}

func concatInstructions(instructions []code.Instructions) code.Instructions {
	out := code.Instructions{}
	for _, ins := range instructions {
		out = append(out, ins...)
	}

	return out
}

func testConstants(t *testing.T, expected []any, actual []object.Object) error {
	t.Helper()
	require.Len(t, actual, len(expected))

	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			testutil.AssertIntegerObject(t, int64(constant), actual[i])
		}
	}
	return nil
}
