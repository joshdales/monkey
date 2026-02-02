package compiler_test

import (
	"monkey/code"
	"monkey/compiler"
	"monkey/testutil"
	"testing"

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

		testutil.AssertInstructions(t, bytecode.Instructions, tt.expectedInstructions)
		testutil.AssertConstants(t, bytecode.Constants, tt.expectedConstants)

	}
}
