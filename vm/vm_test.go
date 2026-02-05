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
		{"1 - 2", -1},
		{"1 * 2", 2},
		{"4 / 2", 2},
		{"50 / 2 * 2 - 5", 45},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"5 * (2 * 10)", 100},
	}

	runVmTest(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"true", true},
		{"false", false},
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
