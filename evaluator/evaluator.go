// evaluator/evaluator.go
package evaluator

import (
    "carrionlang/ast"
    "carrionlang/object"
    "fmt"
)

var (
    NULL = &object.Null{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
    switch node := node.(type) {

    case *ast.Program:
        return evalProgram(node, env)

    case *ast.ExpressionStatement:
        return Eval(node.Expression, env)

    case *ast.IntegerLiteral:
        return &object.Integer{Value: node.Value}

    case *ast.StringLiteral:
        return &object.String{Value: node.Value}

    case *ast.Identifier:
        return evalIdentifier(node, env)

    case *ast.VariableDeclaration:
        val := Eval(node.Value, env)
        if isError(val) {
            return val
        }
        env.Set(node.Name.Value, val)
        return val

    case *ast.InfixExpression:
        left := Eval(node.Left, env)
        if isError(left) {
            return left
        }
        right := Eval(node.Right, env)
        if isError(right) {
            return right
        }
        return evalInfixExpression(node.Operator, left, right)

    case *ast.SpellbookDeclaration:
        return evalSpellbookDeclaration(node, env)

    case *ast.CallExpression:
        function := Eval(node.Function, env)
        if isError(function) {
            return function
        }
        args := evalExpressions(node.Arguments, env)
        if len(args) == 1 && isError(args[0]) {
            return args[0]
        }
        return applyFunction(function, args)

    default:
        return nil
    }
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
    var result object.Object

    for _, statement := range program.Statements {
        result = Eval(statement, env)

        if returnValue, ok := result.(*object.ReturnValue); ok {
            return returnValue.Value
        }
    }
    return result
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
    if val, ok := env.Get(node.Value); ok {
        return val
    }

    // Check built-in functions
    if builtin, ok := builtins[node.Value]; ok {
        return builtin
    }

    return newError("identifier not found: " + node.Value)
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
    switch {
    case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
        return evalIntegerInfixExpression(operator, left, right)
    case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
        return evalStringInfixExpression(operator, left, right)
    default:
        return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
    }
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
    leftVal := left.(*object.Integer).Value
    rightVal := right.(*object.Integer).Value

    switch operator {
    case "+":
        return &object.Integer{Value: leftVal + rightVal}
    case "-":
        return &object.Integer{Value: leftVal - rightVal}
    case "*":
        return &object.Integer{Value: leftVal * rightVal}
    case "/":
        return &object.Integer{Value: leftVal / rightVal}
    default:
        return newError("unknown operator: %s", operator)
    }
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
    leftVal := left.(*object.String).Value
    rightVal := right.(*object.String).Value

    switch operator {
    case "+":
        return &object.String{Value: leftVal + rightVal}
    default:
        return newError("unknown operator: %s", operator)
    }
}

func evalSpellbookDeclaration(sd *ast.SpellbookDeclaration, env *object.Environment) object.Object {
    className := sd.Name.Value
    classEnv := object.NewEnclosedEnvironment(env)

    for _, stmt := range sd.Body {
        Eval(stmt, classEnv)
    }

    classObject := &object.Class{
        Name: className,
        Env:  classEnv,
    }

    env.Set(className, classObject)
    return classObject
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
    var result []object.Object

    for _, e := range exps {
        evaluated := Eval(e, env)
        if isError(evaluated) {
            return []object.Object{evaluated}
        }
        result = append(result, evaluated)
    }

    return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
    switch fn := fn.(type) {
    case *object.Builtin:
        return fn.Fn(args...)
    default:
        return newError("not a function: %s", fn.Type())
    }
}

func DefineBuiltins(env *object.Environment) {
    for name, builtin := range builtins {
        env.Set(name, builtin)
    }
}

func isError(obj object.Object) bool {
    return obj != nil && obj.Type() == object.ERROR_OBJ
}

func newError(format string, a ...interface{}) *object.Error {
    return &object.Error{Message: fmt.Sprintf(format, a...)}
}

