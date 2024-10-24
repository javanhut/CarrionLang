// object/object.go
package object

import (
    "fmt"
    "carrionlang/ast"
)

type ObjectType string

const (
    INTEGER_OBJ          ObjectType = "INTEGER"
    STRING_OBJ           ObjectType = "STRING"
    FUNCTION_OBJ         ObjectType = "FUNCTION"
    BUILTIN_OBJ          ObjectType = "BUILTIN"
    BUILTIN_OBJECT_OBJ   ObjectType = "BUILTIN_OBJECT"
    CLASS_OBJ            ObjectType = "CLASS"
    INSTANCE_OBJ         ObjectType = "INSTANCE"
    NULL_OBJ             ObjectType = "NULL"
    RETURN_VALUE_OBJ     ObjectType = "RETURN_VALUE"
    ERROR_OBJ            ObjectType = "ERROR"
)

type Object interface {
    Type() ObjectType
    Inspect() string
}

type Environment struct {
    store map[string]Object
    outer *Environment
}

func NewEnvironment() *Environment {
    return &Environment{store: make(map[string]Object)}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
    env := NewEnvironment()
    env.outer = outer
    return env
}

func (e *Environment) Get(name string) (Object, bool) {
    obj, ok := e.store[name]
    if !ok && e.outer != nil {
        obj, ok = e.outer.Get(name)
    }
    return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
    e.store[name] = val
    return val
}

// Integer Object
type Integer struct {
    Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

// String Object
type String struct {
    Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

// Function Object
type Function struct {
    Parameters []*ast.Identifier
    Body       *ast.BlockStatement
    Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
    var out string
    params := []string{}
    for _, p := range f.Parameters {
        params = append(params, p.String())
    }
    out += "fn(" + fmt.Sprintf("%s", params) + ") {\n"
    out += f.Body.String()
    out += "\n}"
    return out
}

// Builtin Function Object
type BuiltinFunction func(args ...Object) Object

type Builtin struct {
    Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

// BuiltinObject represents objects with built-in methods (e.g., 'munin')
type BuiltinObject struct {
    Properties map[string]Object
}

func (bo *BuiltinObject) Type() ObjectType { return BUILTIN_OBJECT_OBJ }
func (bo *BuiltinObject) Inspect() string  { return "builtin object" }

// Class Object
type Class struct {
    Name string
    Env  *Environment
}

func (c *Class) Type() ObjectType { return CLASS_OBJ }
func (c *Class) Inspect() string  { return "<class " + c.Name + ">" }

// Instance Object
type Instance struct {
    Class *Class
    Env   *Environment
}

func (i *Instance) Type() ObjectType { return INSTANCE_OBJ }
func (i *Instance) Inspect() string  { return "<instance of " + i.Class.Name + ">" }

// Null Object
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

// Return Value Object
type ReturnValue struct {
    Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

// Error Object
type Error struct {
    Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

