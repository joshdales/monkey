package vm_test

import (
	"monkey/testutil"
	"testing"

	"github.com/stretchr/testify/require"
)

type vmTestCase struct {
	input    string
	expected any
}

func runVmTest(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		comp := testutil.Compile(t, tt.input)
		vm := New(comp.Bytecode())
		err := vm.Run()
		require.NoError(t, err)
		stackElm := vm.StackTop

		testutil.AssertObject(t, stackElm, tt.expected)
	}
}
