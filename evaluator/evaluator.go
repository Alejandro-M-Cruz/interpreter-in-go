package evaluator

import (
	"example.com/writing-an-interpreter/ast"
	"example.com/writing-an-interpreter/object"
	"fmt"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch n := node.(type) {
	case *ast.Program:
		return evalProgram(n.Statements, env)
	case *ast.BlockStatement:
		return evalBlockStatement(n.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(n.Expression, env)
	case *ast.LetStatement:
		value := Eval(n.Value, env)
		if isError(value) {
			return value
		}
		env.Set(n.Name.Value, value)
		return nil
	case *ast.Identifier:
		return evalIdentifier(n.Value, env)
	case *ast.CallExpression:
		function := Eval(n.Function, env)
		if isError(function) {
			return function
		}
		arguments := evalExpressions(n.Arguments, env)
		if len(arguments) == 1 && isError(arguments[0]) {
			return arguments[0]
		}
		return evalCallExpression(function, arguments)
	case *ast.FunctionLiteral:
		return &object.Function{
			Parameters:  n.Parameters,
			Body:        n.Body,
			Environment: env,
		}
	case *ast.ReturnStatement:
		value := Eval(n.ReturnValue, env)
		if isError(value) {
			return value
		}
		return &object.ReturnValue{Value: value}
	case *ast.IfExpression:
		return evalIfExpression(n, env)
	case *ast.InfixExpression:
		left := Eval(n.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(n.Operator, left, right)
	case *ast.PrefixExpression:
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(n.Operator, right)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: n.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(n.Value)
	case *ast.StringLiteral:
		return &object.String{Value: n.Value}
	default:
		return nil
	}
}

func evalProgram(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement, env)

		switch r := result.(type) {
		case *object.ReturnValue:
			return r.Value
		case *object.Error:
			return r
		}
	}

	return result
}

func evalBlockStatement(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement, env)

		if result == nil {
			continue
		}

		rt := result.Type()
		if rt == object.RETURN_VALUE || rt == object.ERROR {
			return result
		}
	}

	return result
}

func evalExpressions(expressions []ast.Expression, env *object.Environment) []object.Object {
	var results []object.Object

	for _, exp := range expressions {
		evaluated := Eval(exp, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		results = append(results, evaluated)
	}

	return results
}

func evalCallExpression(f object.Object, args []object.Object) object.Object {
	function, ok := f.(*object.Function)

	if !ok {
		return newError("not a function: %s", f.Type())
	}

	extendedEnv := object.NewEnclosedEnvironment(function.Environment)

	for i, arg := range args {
		extendedEnv.Set(function.Parameters[i].Value, arg)
	}

	return unwrapReturnValue(Eval(function.Body, extendedEnv))
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	}

	if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	}

	return NULL
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER && right.Type() == object.INTEGER:
		return evalIntegerInfixExpression(operator, left.(*object.Integer), right.(*object.Integer))
	case left.Type() == object.STRING && right.Type() == object.STRING:
		return evalStringInfixExpression(operator, left.(*object.String), right.(*object.String))
	case operator == "==":
		return nativeBoolToBooleanObject(left == right) // pointer comparison -> true and false are always the same
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, l *object.Integer, r *object.Integer) object.Object {
	switch operator {
	case "+":
		return &object.Integer{Value: l.Value + r.Value}
	case "-":
		return &object.Integer{Value: l.Value - r.Value}
	case "*":
		return &object.Integer{Value: l.Value * r.Value}
	case "/":
		return &object.Integer{Value: l.Value / r.Value}
	case ">":
		return nativeBoolToBooleanObject(l.Value > r.Value)
	case "<":
		return nativeBoolToBooleanObject(l.Value < r.Value)
	case "==":
		return nativeBoolToBooleanObject(l.Value == r.Value)
	case "!=":
		return nativeBoolToBooleanObject(l.Value != r.Value)
	default:
		return newError("unknown operator: %s %s %s", l.Type(), operator, r.Type())
	}
}

func evalStringInfixExpression(operator string, l *object.String, r *object.String) object.Object {
	switch operator {
	case "+":
		return &object.String{Value: l.Value + r.Value}
	case "==":
		return nativeBoolToBooleanObject(l.Value == r.Value)
	case "!=":
		return nativeBoolToBooleanObject(l.Value != r.Value)
	default:
		return newError("unknown operator: %s %s %s", l.Type(), operator, r.Type())
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorPrefixExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	return nativeBoolToBooleanObject(!isTruthy(right))
}

func evalMinusOperatorPrefixExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER {
		return newError("unknown operator: -%s", right.Type())
	}
	return &object.Integer{Value: -right.(*object.Integer).Value}
}

func nativeBoolToBooleanObject(value bool) *object.Boolean {
	if value {
		return TRUE
	}
	return FALSE
}

func isTruthy(value object.Object) bool {
	switch value {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	}

	switch value.Type() {
	case object.INTEGER:
		return value.(*object.Integer).Value != 0
	case object.STRING:
		return value.(*object.String).Value != ""
	default:
		return true
	}
}

func newError(format string, a ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	return obj != nil && obj.Type() == object.ERROR
}

func evalIdentifier(name string, env *object.Environment) object.Object {
	if value, ok := env.Get(name); ok {
		return value
	}
	return newError("identifier not found: %s", name)
}
