package testutil

import (
	"monkey/ast"
	"monkey/compiler"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"

	"github.com/stretchr/testify/require"
)

func SetupProgram(t *testing.T, input string, stmtLen int) *ast.Program {
	t.Helper()
	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	require.NotNil(t, program, "program.ParseProgram() returned nil")
	if stmtLen > 0 {
		checkParserErrors(t, p)
		require.Len(t, program.Statements, stmtLen, "program.Statements does not contain %d statements.", stmtLen)
	}
	return program
}

func Compile(t *testing.T, input string) *compiler.Compiler {
	t.Helper()
	program := SetupProgram(t, input, 0)
	comp := compiler.New()
	err := comp.Compile(program)
	require.NoError(t, err)
	return comp
}

func checkParserErrors(t *testing.T, p *parser.Parser) {
	t.Helper()
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestEval(t *testing.T, input string) object.Object {
	t.Helper()

	env := object.NewEnvironment()
	program := SetupProgram(t, input, 0)

	return evaluator.Eval(env, program)
}
