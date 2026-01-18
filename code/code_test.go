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
