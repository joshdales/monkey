package compiler

import (
	"monkey/code"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompilerScopes(t *testing.T) {
	compiler := New()
	assert.Equalf(t, 0, compiler.scopeIndex, "scopeIndex wrong. got=%d, want=%d", compiler.scopeIndex, 0)

	globalSymbolTable := compiler.symbolTable

	compiler.emit(code.OpMul)

	compiler.enterScope()
	assert.Equalf(t, 1, compiler.scopeIndex, "scopeIndex wrong. got=%d, want=%d", compiler.scopeIndex, 1)

	compiler.emit(code.OpSub)
	assert.Len(t, compiler.scopes[compiler.scopeIndex].instructions, 1, "instructions length wrong")

	last := compiler.scopes[compiler.scopeIndex].lastInstruction
	assert.Equal(t, code.OpSub, last.Opcode, "lastInstruction.OpCode wrong. got=%d, want=%d")

	assert.Equal(t, compiler.symbolTable.Outer, globalSymbolTable, "compiler did not enclose symbolTable")

	compiler.leaveScope()
	assert.Equalf(t, 0, compiler.scopeIndex, "scopeIndex wrong. got=%d, want=%d", compiler.scopeIndex, 0)

	assert.Equal(t, compiler.symbolTable, globalSymbolTable, "compiler did not restore global symbol table")
	assert.Nil(t, compiler.symbolTable.Outer, "compiler modified glodal symbol table incorrectly")

	compiler.emit(code.OpAdd)
	assert.Len(t, compiler.scopes[compiler.scopeIndex].instructions, 2, "instructions length wrong")

	last = compiler.scopes[compiler.scopeIndex].lastInstruction
	assert.Equalf(t, code.OpAdd, last.Opcode, "lastInstruction.OpCode wrong. got=%d, want=%d", last.Opcode, code.OpAdd)

	previous := compiler.scopes[compiler.scopeIndex].previousInstruction
	assert.Equalf(t, code.OpMul, previous.Opcode, "previousInstruction.OpCode wrong. got=%d, want=%d", previous.Opcode, code.OpMul)
}
