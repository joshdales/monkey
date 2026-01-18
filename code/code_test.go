package code_test

import (
	"monkey/code"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       code.Opcode
		operands []int
		expected []byte
	}{
		{code.OpConstant, []int{65534}, []byte{byte(code.OpConstant), 255, 254}},
	}

	for _, tt := range tests {
		instruction := code.Make(tt.op, tt.operands...)
		assert.Len(t, instruction, len(tt.expected), "instruction has wrong length")
		assert.True(t, assert.ObjectsAreEqualValues(tt.expected, instruction))
	}
}

func TestInstructionsString(t *testing.T) {
	instructions := []code.Instructions{
		code.Make(code.OpConstant, 1),
		code.Make(code.OpConstant, 2),
		code.Make(code.OpConstant, 65535),
	}

	expected := `0000 Opconstant 1
0003 OpConstant 2
0006 OpConstant 65535`

	concatted := code.Instructions{}
	for _, ins := range instructions {
		concatted = append(concatted, ins...)
	}

	assert.Equal(t, expected, concatted.String(), "instructions wrongly formatted")
}
