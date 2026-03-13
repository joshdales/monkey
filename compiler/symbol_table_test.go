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
		"c": {Name: "c", Scope: compiler.LocalScope, Index: 0},
		"d": {Name: "d", Scope: compiler.LocalScope, Index: 1},
		"e": {Name: "e", Scope: compiler.LocalScope, Index: 0},
		"f": {Name: "f", Scope: compiler.LocalScope, Index: 1},
	}

	global := compiler.NewSymbolTable()

	a := global.Define("a")
	assert.Equal(t, expected["a"], a)
	b := global.Define("b")
	assert.Equal(t, expected["b"], b)

	firstLocal := compiler.NewEnclosedSymbolTable(global)
	c := firstLocal.Define("c")
	assert.Equal(t, expected["c"], c)
	d := firstLocal.Define("d")
	assert.Equal(t, expected["d"], d)

	secondLocal := compiler.NewEnclosedSymbolTable(firstLocal)
	e := secondLocal.Define("e")
	assert.Equal(t, expected["e"], e)
	f := secondLocal.Define("f")
	assert.Equal(t, expected["f"], f)
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

func TestResoveLocal(t *testing.T) {
	global := compiler.NewSymbolTable()
	global.Define("a")
	global.Define("b")

	local := compiler.NewEnclosedSymbolTable(global)
	local.Define("c")
	local.Define("d")

	expected := []compiler.Symbol{
		{Name: "a", Scope: compiler.GlobalScope, Index: 0},
		{Name: "b", Scope: compiler.GlobalScope, Index: 1},
		{Name: "c", Scope: compiler.LocalScope, Index: 0},
		{Name: "d", Scope: compiler.LocalScope, Index: 1},
	}

	for _, sym := range expected {
		t.Run(sym.Name, func(t *testing.T) {
			result, ok := local.Resolve(sym.Name)
			require.Truef(t, ok, "name %s not resolvable", sym.Name)
			assert.Equalf(t, sym, result, "expected %s to resolve to %+v, got %+v", sym.Name, sym, result)
		})
	}
}

func TestResoveNestedLocal(t *testing.T) {
	global := compiler.NewSymbolTable()
	global.Define("a")
	global.Define("b")

	firstLocal := compiler.NewEnclosedSymbolTable(global)
	firstLocal.Define("c")
	firstLocal.Define("d")

	secondLocal := compiler.NewEnclosedSymbolTable(firstLocal)
	secondLocal.Define("e")
	secondLocal.Define("f")

	tests := []struct {
		table           *compiler.SymbolTable
		expectedSymbols []compiler.Symbol
	}{
		{firstLocal, []compiler.Symbol{
			{Name: "a", Scope: compiler.GlobalScope, Index: 0},
			{Name: "b", Scope: compiler.GlobalScope, Index: 1},
			{Name: "c", Scope: compiler.LocalScope, Index: 0},
			{Name: "d", Scope: compiler.LocalScope, Index: 1},
		}},
		{secondLocal, []compiler.Symbol{
			{Name: "a", Scope: compiler.GlobalScope, Index: 0},
			{Name: "b", Scope: compiler.GlobalScope, Index: 1},
			{Name: "e", Scope: compiler.LocalScope, Index: 0},
			{Name: "f", Scope: compiler.LocalScope, Index: 1},
		}},
	}

	for _, tt := range tests {
		for _, sym := range tt.expectedSymbols {
			t.Run(sym.Name, func(t *testing.T) {
				result, ok := tt.table.Resolve(sym.Name)
				require.Truef(t, ok, "name %s not resolvable", sym.Name)
				assert.Equalf(t, sym, result, "expected %s to resolve to %+v, got %+v", sym.Name, sym, result)
			})
		}
	}
}

func TestDefineResolveBuiltins(t *testing.T) {
	global := compiler.NewSymbolTable()
	firstLocal := compiler.NewEnclosedSymbolTable(global)
	secondLocal := compiler.NewEnclosedSymbolTable(firstLocal)

	expected := []compiler.Symbol{
		{Name: "a", Scope: compiler.BuiltinScope, Index: 0},
		{Name: "c", Scope: compiler.BuiltinScope, Index: 1},
		{Name: "e", Scope: compiler.BuiltinScope, Index: 2},
		{Name: "f", Scope: compiler.BuiltinScope, Index: 3},
	}

	for i, v := range expected {
		global.DefineBuiltin(i, v.Name)
	}

	for _, table := range []*compiler.SymbolTable{global, firstLocal, secondLocal} {
		for _, sym := range expected {
			t.Run(sym.Name, func(t *testing.T) {
				result, ok := table.Resolve(sym.Name)
				require.Truef(t, ok, "name %s not resolvable", sym.Name)
				assert.Equalf(t, sym, result, "expected %s to resolve to %+v, got %+v", sym.Name, sym, result)
			})
		}
	}
}
