package evaluator

import (
	"encoding/json"
	"errors"
	"example.com/writing-an-interpreter/object"
	"net/http"
	"os"
	"rsc.io/quote/v4"
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
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("invalid argument for the `len` function, got %s", arg.Type())
			}
		},
	},
	"quote": {
		Fn: func(args ...object.Object) object.Object {
			result := &object.String{}
			q, err := getRandomQuote()

			if err != nil {
				result.Value = quote.Opt()
			} else {
				result.Value = q
			}

			return result
		},
	},
	"append": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 2 {
				return newError("expected at least 2 arguments, %d were given", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Array{Elements: append(arg.Elements, args[1:]...)}
			default:
				return newError("invalid argument for the `append` function, got %s", arg.Type())
			}
		},
	},
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
