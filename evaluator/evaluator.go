package evaluator

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(env *object.Environment, node ast.Node) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(env, node)
	case *ast.BlockStatement:
		return evalBlockStatements(env, node)
	case *ast.IfExpression:
		return evalIfStatement(env, node)
	case *ast.LetStatement:
		val := Eval(env, node.Value)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return nil
	case *ast.ExpressionStatement:
		return Eval(env, node.Expression)
	case *ast.ReturnStatement:
		val := Eval(env, node.ReturnValue)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	// Expressions
	case *ast.PrefixExpression:
		right := Eval(env, node.Right)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(env, node.Left)
		if isError(left) {
			return left
		}
		right := Eval(env, node.Right)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.CallExpression:
		return evalCallExpression(env, node)
	case *ast.IndexExpression:
		left := Eval(env, node.Left)
		if isError(left) {
			return left
		}
		index := Eval(env, node.Index)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	// Literals
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.ArrayLiteral:
		elements := evalExpressions(env, node.Elements)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.HashLiteral:
		return evalHashLiteral(env, node)
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.Identifier:
		return evalIdentifier(env, node)
	default:
		return nil
	}
}

func evalIdentifier(env *object.Environment, node *ast.Identifier) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: %s", node.Value)
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func newError(format string, a ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}

	return false
}

func evalHashLiteral(env *object.Environment, node *ast.HashLiteral) object.Object {
	pairs := make(map[object.HashKey]object.HashPair, len(node.Pairs))

	for keyNode, valueNode := range node.Pairs {
		key := Eval(env, keyNode)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as a hash key: %s", key.Type())
		}

		value := Eval(env, valueNode)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}
