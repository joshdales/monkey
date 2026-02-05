package vm_test

import (
	"monkey/testutil"
	"monkey/vm"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
	}

	runVmTest(t, tests)
}

type vmTestCase struct {
	input    string
	expected any
}

func runVmTest(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		comp := testutil.Compile(t, tt.input)
		vm := vm.New(comp.Bytecode())
		err := vm.Run()
		require.NoError(t, err)
		stackElm := vm.LastPoppedStackElem()

		testutil.AssertObject(t, stackElm, tt.expected)
	}
}
