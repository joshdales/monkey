package compiler_test

import (
	"monkey/compiler"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefine(t *testing.T) {
	expected := map[string]compiler.Symbol{
		"a": {Name: "a", Scope: compiler.GlobalScope, Index: 0},
		"b": {Name: "b", Scope: compiler.GlobalScope, Index: 1},
	}

	global := compiler.NewSymbolTable()

	a := global.Define("a")
	assert.Equal(t, expected["a"], a)
	b := global.Define("b")
	assert.Equal(t, expected["b"], b)
}

func TestResolveGlobal(t *testing.T) {
	global := compiler.NewSymbolTable()
	global.Define("a")
	global.Define("b")

	expected := []compiler.Symbol{
		{Name: "a", Scope: compiler.GlobalScope, Index: 0},
		{Name: "b", Scope: compiler.GlobalScope, Index: 1},
	}

	for _, sym := range expected {
		t.Run(sym.Name, func(t *testing.T) {
			result, ok := global.Resolve(sym.Name)
			require.Truef(t, ok, "name %s not resolvable", sym.Name)
			assert.Equalf(t, sym, result, "expected %s to resolve to %+v, got %+v", sym.Name, sym, result)
		})
	}
}
