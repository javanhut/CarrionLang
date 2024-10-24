// object/object.go
package object

import (
    "fmt"
    "hash/fnv"
    "strings"
)

type ObjectType string

const (
    INTEGER_OBJ ObjectType = "INTEGER"
    STRING_OBJ  ObjectType = "STRING"
    FUNCTION_OBJ ObjectType = "FUNCTION"
    BUILTIN_OBJ ObjectType = "BUILTIN"
    CLASS_OBJ   ObjectType = "CLASS"
    INSTANCE_OBJ ObjectType = "INSTANCE"
    NULL_OBJ    ObjectType = "NULL"
    RETURN_VALUE_OBJ ObjectType = "RETURN_VALUE"
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

type Integer struct {
    Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type String struct {
    Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

type Function struct {
    Parameters []*ast.Identifier
    Body       *ast.BlockStatement
    Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
    var out strings.Builder
    params := []string{}
    for _, p := range f.Parameters {
        params = append(params, p.String())
    }
    out.WriteString("fn(")
    out.WriteString(strings.Join(params, ", "))
    out.WriteString(") {\n")
    out.WriteString(f.Body.String())
    out.WriteString("\n}")
    return out.String()
}

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
    Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

type Class struct {
    Name string
    Env  *Environment
}

func (c *Class) Type() ObjectType { return CLASS_OBJ }
func (c *Class) Inspect() string  { return "<class " + c.Name + ">" }

type Instance struct {
    Class *Class
    Env   *Environment
}

func (i *Instance) Type() ObjectType { return INSTANCE_OBJ }
func (i *Instance) Inspect() string  { return "<instance of " + i.Class.Name + ">" }

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

type ReturnValue struct {
    Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

