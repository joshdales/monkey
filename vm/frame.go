package vm

import (
	"monkey/code"
	"monkey/object"
)

// A Frame holds execution-relevant information.
type Frame struct {
	fn          *object.CompiledFunction // Compiled function referenced by this Frame
	ip          int                      // Instruction pointer in this Frame, for this Frame
	basepointer int                      // The pointer value before the function was executed
}

func NewFrame(fn *object.CompiledFunction, basePointer int) *Frame {
	return &Frame{fn: fn, ip: -1, basepointer: basePointer}
}

func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}
