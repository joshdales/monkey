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
		{code.OpGetLocal, []int{255}, []byte{byte(code.OpGetLocal), 255}},
		{code.OpClosure, []int{65534, 255}, []byte{byte(code.OpClosure), 255, 254, 255}},
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
		code.Make(code.OpGetLocal, 1),
		code.Make(code.OpConstant, 2),
		code.Make(code.OpConstant, 65535),
		code.Make(code.OpClosure, 65535, 255),
	}

	expected := `0000 OpAdd
0001 OpGetLocal 1
0003 OpConstant 2
0006 OpConstant 65535
0009 OpClosure 65535 255
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
		name      string
		op        code.Opcode
		operands  []int
		bytesRead int
	}{
		{"Contstant", code.OpConstant, []int{65535}, 2},
		{"GetLocal", code.OpGetLocal, []int{255}, 1},
		{"Closure", code.OpClosure, []int{65535, 255}, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instruction := code.Make(tt.op, tt.operands...)
			def, err := code.Lookup(byte(tt.op))
			require.NoErrorf(t, err, "definition not found: %q", err)
			operandsRead, n := code.ReadOperands(def, instruction[1:])
			assert.Equal(t, tt.bytesRead, n)
			for i, want := range tt.operands {
				assert.Equal(t, want, operandsRead[i])
			}
		})
	}
}
