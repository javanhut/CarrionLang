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

// Eval evaluates the AST nodes
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

    case *ast.SpellDeclaration:
        return evalSpellDeclaration(node, env)

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

    case *ast.ReturnStatement:
        val := Eval(node.ReturnValue, env)
        if isError(val) {
            return val
        }
        return &object.ReturnValue{Value: val}

    case *ast.MemberExpression:
        return evalMemberExpression(node, env)

    // Add more cases as needed...
    }
    return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
    var result object.Object

    for _, statement := range program.Statements {
        result = Eval(statement, env)

        switch result := result.(type) {
        case *object.ReturnValue:
            return result.Value
        case *object.Error:
            return result
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
        if rightVal == 0 {
            return newError("division by zero")
        }
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

func evalSpellDeclaration(sd *ast.SpellDeclaration, env *object.Environment) object.Object {
    spell := &object.Function{
        Parameters: sd.Parameters,
        Body:       sd.Body,
        Env:        env,
    }

    env.Set(sd.Name.Value, spell)
    return spell
}

func evalMemberExpression(me *ast.MemberExpression, env *object.Environment) object.Object {
    object := Eval(me.Object, env)
    if isError(object) {
        return object
    }

    return evalPropertyAccess(object, me.Property.Value)
}

func evalPropertyAccess(obj object.Object, property string) object.Object {
    switch obj := obj.(type) {
    case *object.BuiltinObject:
        if prop, ok := obj.Properties[property]; ok {
            return prop
        }
        return newError("property '%s' not found on built-in object", property)
    case *object.Instance:
        // Handle instance properties if needed
        return newError("property access not implemented for instances")
    case *object.Class:
        // Handle class properties or methods if needed
        return newError("property access not implemented for classes")
    default:
        return newError("cannot access property '%s' of type %s", property, obj.Type())
    }
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
    case *object.Function:
        extendedEnv := extendFunctionEnv(fn, args)
        evaluated := Eval(fn.Body, extendedEnv)
        return unwrapReturnValue(evaluated)
    case *object.Builtin:
        return fn.Fn(args...)
    default:
        return newError("not a function: %s", fn.Type())
    }
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
    env := object.NewEnclosedEnvironment(fn.Env)
    for paramIdx, param := range fn.Parameters {
        if paramIdx < len(args) {
            env.Set(param.Value, args[paramIdx])
        } else {
            env.Set(param.Value, NULL) // Default to NULL if not enough args
        }
    }
    return env
}

func unwrapReturnValue(obj object.Object) object.Object {
    if returnValue, ok := obj.(*object.ReturnValue); ok {
        return returnValue.Value
    }
    return obj
}

func DefineBuiltins(env *object.Environment) {
    // Define 'munin' object with 'print' method
    munin := &object.BuiltinObject{
        Properties: map[string]object.Object{
            "print": &object.Builtin{
                Fn: func(args ...object.Object) object.Object {
                    for _, arg := range args {
                        fmt.Print(arg.Inspect())
                    }
                    fmt.Println()
                    return NULL
                },
            },
            // Add more methods if needed
        },
    }
    env.Set("munin", munin)
}

func isError(obj object.Object) bool {
    return obj != nil && obj.Type() == object.ERROR_OBJ
}

func newError(format string, a ...interface{}) *object.Error {
    return &object.Error{Message: fmt.Sprintf(format, a...)}
}

