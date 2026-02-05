package compiler_test

import (
	"monkey/code"
	"monkey/testutil"
	"testing"
)

type compilerTestCase struct {
	input                string
	expectedConstants    []any
	expectedInstructions []code.Instructions
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "1 + 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
			},
		},
	}

	runCompilerTests(t, tests)
}

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()

	for _, tt := range tests {
		comp := testutil.Compile(t, tt.input)
		bytecode := comp.Bytecode()

		testutil.AssertInstructions(t, bytecode.Instructions, tt.expectedInstructions)
		testutil.AssertConstants(t, bytecode.Constants, tt.expectedConstants)
	}
}
