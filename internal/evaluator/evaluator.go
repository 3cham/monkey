package evaluator

import (
	"monkey/internal/ast"
	"monkey/internal/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalStatements(node.Statements)
	case *ast.IfExpression:
		condition := Eval(node.Condition)
		if isTruthy(condition) {
			return Eval(node.Consequence)
		} else if node.Alternative != nil {
			return Eval(node.Alternative)
		} else {
			return NULL
		}

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolean(node.Value)
	}

	return NULL
}

func isTruthy(condition object.Object) bool {
	switch condition {
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

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return evalBooleanInfixExpression(operator, left, right)
	default:
		return NULL
	}
}

func evalBooleanInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftEval := left.(*object.Boolean).Value
	rightEval := right.(*object.Boolean).Value
	switch operator {
	case "==":
		return nativeBoolean(leftEval == rightEval)
	case "!=":
		return nativeBoolean(leftEval != rightEval)
	default:
		return NULL
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftEval := left.(*object.Integer).Value
	rightEval := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftEval + rightEval}
	case "-":
		return &object.Integer{Value: leftEval - rightEval}
	case "*":
		return &object.Integer{Value: leftEval * rightEval}
	case "/":
		return &object.Integer{Value: leftEval / rightEval}
	case "<":
		return nativeBoolean(leftEval < rightEval)
	case ">":
		return nativeBoolean(leftEval > rightEval)
	case "==":
		return nativeBoolean(leftEval == rightEval)
	case "!=":
		return nativeBoolean(leftEval != rightEval)
	default:
		return NULL
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return NULL
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return NULL
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func nativeBoolean(value bool) object.Object {
	if value {
		return TRUE
	}
	return FALSE
}

func evalStatements(statements []ast.Statement) object.Object {
	var result object.Object
	for _, statement := range statements {
		result = Eval(statement)
	}

	return result
}
