package vm

import (
	"monkey/code"
	"monkey/object"
)

// A Frame holds execution-relevant information.
type Frame struct {
	fn *object.CompiledFunction // Compiled function referenced by this Frame
	ip int                      // Instruction pointer in this Frame, for this Frame
}

func NewFrame(fn *object.CompiledFunction) *Frame {
	return &Frame{fn: fn, ip: -1}
}

func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}
