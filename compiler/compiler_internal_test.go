package compiler

import (
	"monkey/code"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompilerScopes(t *testing.T) {
	compiler := New()
	assert.Equalf(t, 0, compiler.scopeIndex, "scopeIndex wrong. got=%d, want=%d", compiler.scopeIndex, 0)

	compiler.emit(code.OpMul)
	compiler.enterScope()
	assert.Equalf(t, 1, compiler.scopeIndex, "scopeIndex wrong. got=%d, want=%d", compiler.scopeIndex, 1)

	compiler.emit(code.OpSub)
	assert.Len(t, compiler.scopes[compiler.scopeIndex].instructions, 1, "instructions length wrong")

	last := compiler.scopes[compiler.scopeIndex].lastInstruction
	assert.Equal(t, code.OpSub, last.Opcode, "lastInstruction.OpCode wrong. got=%d, want=%d")

	compiler.leaveScope()
	assert.Equalf(t, 0, compiler.scopeIndex, "scopeIndex wrong. got=%d, want=%d", compiler.scopeIndex, 0)

	compiler.emit(code.OpAdd)
	assert.Len(t, compiler.scopes[compiler.scopeIndex].instructions, 2, "instructions length wrong")

	last = compiler.scopes[compiler.scopeIndex].lastInstruction
	assert.Equal(t, code.OpAdd, last.Opcode, "lastInstruction.OpCode wrong. got=%d, want=%d")

	previous := compiler.scopes[compiler.scopeIndex].previousInstruction
	assert.Equal(t, code.OpMul, previous.Opcode, "lastInstruction.OpCode wrong. got=%d, want=%d")
}
