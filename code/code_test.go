package code_test

import (
	"monkey/code"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       code.Opcode
		operands []int
		expected []byte
	}{
		{code.OpConstant, []int{65534}, []byte{byte(code.OpConstant), 255, 254}},
		{code.OpAdd, []int{}, []byte{byte(code.OpAdd)}},
	}

	for _, tt := range tests {
		instruction := code.Make(tt.op, tt.operands...)
		assert.Len(t, instruction, len(tt.expected), "instruction has wrong length")
		assert.True(t, assert.ObjectsAreEqualValues(tt.expected, instruction))
	}
}

func TestInstructionsString(t *testing.T) {
	instructions := []code.Instructions{
		code.Make(code.OpAdd),
		code.Make(code.OpConstant, 2),
		code.Make(code.OpConstant, 65535),
	}

	expected := `0000 OpAdd
0001 OpConstant 2
0004 OpConstant 65535
`

	concatted := code.Instructions{}
	for _, ins := range instructions {
		concatted = append(concatted, ins...)
	}

	assert.Equalf(t,
		expected,
		concatted.String(),
		"instructions wrongly formatted.\nwant %q\ngot %q", expected, concatted.String(),
	)
}

func TestReadOperands(t *testing.T) {
	tests := []struct {
		op        code.Opcode
		operands  []int
		bytesRead int
	}{
		{code.OpConstant, []int{65535}, 2},
	}

	for _, tt := range tests {
		instruction := code.Make(tt.op, tt.operands...)
		def, err := code.Lookup(byte(tt.op))
		require.NoErrorf(t, err, "definition not found: %q", err)
		operandsRead, n := code.ReadOperands(def, instruction[1:])
		assert.Equal(t, tt.bytesRead, n)
		for i, want := range tt.operands {
			assert.Equal(t, want, operandsRead[i])
		}
	}
}
