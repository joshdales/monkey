package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

func evalProgram(env *object.Environment, program *ast.Program) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(env, statement)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatements(env *object.Environment, block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(env, statement)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalIfStatement(env *object.Environment, ie *ast.IfExpression) object.Object {
	condition := Eval(env, ie.Condition)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(env, ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(env, ie.Alternative)
	} else {
		return NULL
	}
}
