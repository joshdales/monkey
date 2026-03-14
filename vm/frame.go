package vm

import (
	"monkey/code"
	"monkey/object"
)

// A Frame holds execution-relevant information.
type Frame struct {
	cl          *object.Closure // Closure for compiled function referenced by this Frame
	ip          int             // Instruction pointer in this Frame, for this Frame
	basepointer int             // The pointer value before the function was executed
}

func NewFrame(cl *object.Closure, basePointer int) *Frame {
	return &Frame{cl: cl, ip: -1, basepointer: basePointer}
}

func (f *Frame) Instructions() code.Instructions {
	return f.cl.Fn.Instructions
}
