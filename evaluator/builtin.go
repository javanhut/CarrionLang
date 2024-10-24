// evaluator/builtin.go
package evaluator

import (
    "carrionlang/object"
    "fmt"
)

var builtins = map[string]*object.Builtin{
    // Other built-in functions can be added here
}

func init() {
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
            // Add more methods to 'munin' here if needed
        },
    }
    builtins["munin"] = &object.Builtin{
        Fn: func(args ...object.Object) object.Object {
            return munin
        },
    }
}

