package evaluator

import (
	"bytes"
	"encoding/json"
	"errors"
	"example.com/writing-an-interpreter/object"
	"fmt"
	"net/http"
	"os"
	"rsc.io/quote/v4"
	"strings"
)

var builtins = map[string]*object.Builtin{
	"print":  {Fn: builtinPrint},
	"len":    {Fn: builtinLen},
	"append": {Fn: builtinAppend},
	"first":  {Fn: builtinFirst},
	"last":   {Fn: builtinLast},
	"skip":   {Fn: builtinSkip},
	"quote":  {Fn: builtinQuote},
}

func builtinPrint(args ...object.Object) object.Object {
	var arguments []string

	for _, a := range args {
		arguments = append(arguments, a.Inspect())
	}

	fmt.Println(strings.Join(arguments, " "))
	return NULL
}

func builtinLen(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newArgumentNumberError(1, len(args), false)
	}

	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}
	default:
		return newInvalidArgumentError("len", arg)
	}
}

func builtinFirst(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newArgumentNumberError(1, len(args), false)
	}

	switch arg := args[0].(type) {
	case *object.Array:
		if len(arg.Elements) == 0 {
			return NULL
		}
		return arg.Elements[0]
	case *object.String:
		if len(arg.Value) == 0 {
			return NULL
		}
		return &object.String{Value: string(arg.Value[0])}
	default:
		return newInvalidArgumentError("first", arg)
	}
}

func builtinSkip(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newArgumentNumberError(2, len(args), false)
	}

	skip, ok := args[1].(*object.Integer)
	if !ok {
		return newInvalidArgumentError("skip", args[1])
	}
	s := skip.Value

	switch arg := args[0].(type) {
	case *object.Array:
		length := int64(len(arg.Elements))
		if s > length {
			return &object.Array{}
		}
		newElements := make([]object.Object, length-s)
		copy(newElements, arg.Elements[s:])
		return &object.Array{Elements: newElements}
	case *object.String:
		if s > int64(len(arg.Value)) {
			return &object.String{}
		}
		return &object.String{Value: arg.Value[s:]}
	default:
		return newInvalidArgumentError("first", arg)
	}
}

func builtinLast(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newArgumentNumberError(1, len(args), false)
	}

	switch arg := args[0].(type) {
	case *object.Array:
		if len(arg.Elements) == 0 {
			return NULL
		}
		return arg.Elements[len(arg.Elements)-1]
	case *object.String:
		if len(arg.Value) == 0 {
			return NULL
		}
		return &object.String{Value: string(arg.Value[len(arg.Value)-1])}
	default:
		return newInvalidArgumentError("last", arg)
	}
}

func builtinQuote(_ ...object.Object) object.Object {
	result := &object.String{}
	q, err := getRandomQuote()

	if err != nil {
		result.Value = quote.Opt()
	} else {
		result.Value = q
	}

	return result
}

func getRandomQuote() (string, error) {
	response, err := http.Get(os.Getenv("RANDOM_QUOTE_ENDPOINT"))

	if err != nil {
		return "", err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	var quotes Quotes

	if err := json.NewDecoder(response.Body).Decode(&quotes); err != nil {
		return "", err
	}

	if len(quotes) == 0 || len(quotes[0].Quote) == 0 {
		return "", errors.New("could not retrieve a quote")
	}

	return quotes[0].Quote, nil
}

type Quotes []struct {
	Quote  string `json:"q"`
	Author string `json:"a"`
	Html   string `json:"h"`
}

func builtinAppend(args ...object.Object) object.Object {
	if len(args) < 2 {
		return newArgumentNumberError(2, len(args), true)
	}

	switch arg := args[0].(type) {
	case *object.Array:
		return &object.Array{Elements: append(arg.Elements, args[1:]...)}
	default:
		return newInvalidArgumentError("append", arg)
	}
}

func newInvalidArgumentError(functionName string, arg object.Object) *object.Error {
	return newError("invalid argument for the `%s` function, got %s", functionName, arg.Type())
}

func newArgumentNumberError(expected int, given int, canBeMore bool) *object.Error {
	var out bytes.Buffer

	if canBeMore {
		out.WriteString("expected at least")
	} else {
		out.WriteString("expected")
	}

	out.WriteString(" %d ")

	if expected == 1 {
		out.WriteString("argument")
	} else {
		out.WriteString("arguments")
	}

	out.WriteString(", received %d")
	return newError(out.String(), expected, given)
}
