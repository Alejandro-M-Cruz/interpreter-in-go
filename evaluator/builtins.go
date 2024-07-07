package evaluator

import (
	"example.com/writing-an-interpreter/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("expected 1 argument, %d were given", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("invalid argument for the `len` function, got %s", args[0].Type())
			}
		},
	},
}
