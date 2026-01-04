package ast_test

import (
	"monkey/ast"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModify(t *testing.T) {
	one := func() ast.Expression { return &ast.IntegerLiteral{Value: 1} }
	two := func() ast.Expression { return &ast.IntegerLiteral{Value: 2} }

	turnOneIntoTwo := func(node ast.Node) ast.Node {
		integer, ok := node.(*ast.IntegerLiteral)
		if !ok {
			return node
		}

		if integer.Value != 1 {
			return node
		}

		integer.Value = 2
		return integer
	}

	testCases := []struct {
		desc     string
		input    ast.Node
		expected ast.Node
	}{
		{
			"Program",
			&ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{Expression: one()},
				},
			},
			&ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{Expression: two()},
				},
			},
		},
		{
			"InfixExpression1",
			&ast.InfixExpression{Left: one(), Operator: "+", Right: two()},
			&ast.InfixExpression{Left: two(), Operator: "+", Right: two()},
		},
		{
			"InfixExpression2",
			&ast.InfixExpression{Left: two(), Operator: "+", Right: two()},
			&ast.InfixExpression{Left: two(), Operator: "+", Right: two()},
		},
		{
			"PrefixExpression1",
			&ast.PrefixExpression{Operator: "-", Right: one()},
			&ast.PrefixExpression{Operator: "-", Right: two()},
		},
		{
			"PrefixExpression2",
			&ast.IndexExpression{Left: one(), Index: one()},
			&ast.IndexExpression{Left: two(), Index: two()},
		},
		{
			"IfExpression",
			&ast.IfExpression{
				Condition: one(),
				Consequence: &ast.BlockStatement{
					Statements: []ast.Statement{&ast.ExpressionStatement{Expression: one()}},
				},
				Alternative: &ast.BlockStatement{
					Statements: []ast.Statement{&ast.ExpressionStatement{Expression: one()}},
				},
			},
			&ast.IfExpression{
				Condition: one(),
				Consequence: &ast.BlockStatement{
					Statements: []ast.Statement{&ast.ExpressionStatement{Expression: two()}},
				},
				Alternative: &ast.BlockStatement{
					Statements: []ast.Statement{&ast.ExpressionStatement{Expression: two()}},
				},
			},
		},
		{
			"ReturnStatement",
			&ast.ReturnStatement{ReturnValue: one()},
			&ast.ReturnStatement{ReturnValue: two()},
		},
		{
			"LetStatement",
			&ast.LetStatement{Value: one()},
			&ast.LetStatement{Value: two()},
		},
		{
			"FunctionLiteral",
			&ast.FunctionLiteral{
				Parameters: []*ast.Identifier{},
				Body: &ast.BlockStatement{
					Statements: []ast.Statement{
						&ast.ExpressionStatement{Expression: one()},
					},
				},
			},
			&ast.FunctionLiteral{
				Parameters: []*ast.Identifier{},
				Body: &ast.BlockStatement{
					Statements: []ast.Statement{
						&ast.ExpressionStatement{Expression: two()},
					},
				},
			},
		},
		{
			"ArrayLiteral",
			&ast.ArrayLiteral{Elements: []ast.Expression{one(), two()}},
			&ast.ArrayLiteral{Elements: []ast.Expression{two(), two()}},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			modified := ast.Modify(tC.input, turnOneIntoTwo)
			assert.ObjectsAreEqualValues(tC.expected, modified)
		})
	}
}
