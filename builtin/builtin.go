// evaluator/builtins.go
package evaluator

import (
    "carrionlang/object"
    "fmt"
)

var builtins = map[string]*object.Builtin{
    "munin.print": {
        Fn: func(args ...object.Object) object.Object {
            for _, arg := range args {
                fmt.Print(arg.Inspect())
            }
            fmt.Println()
            return NULL
        },
    },
    // Add more built-in functions here
}

