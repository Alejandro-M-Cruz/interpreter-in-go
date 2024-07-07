package repl

import (
	"bufio"
	"example.com/writing-an-interpreter/evaluator"
	"example.com/writing-an-interpreter/lexer"
	"example.com/writing-an-interpreter/object"
	"example.com/writing-an-interpreter/parser"
	"io"
	"strings"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		_, err := out.Write([]byte(PROMPT))
		if err != nil {
			panic(err)
		}

		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.NewLexer(line)
		p := parser.NewParser(l)
		program := p.ParseProgram()
		result := evaluator.Eval(program, env)

		if len(p.Errors()) > 0 {
			err := printParseErrors(out, p.Errors())

			if err != nil {
				panic(err)
			}

			continue
		}

		if result == nil {
			continue
		}

		_, err = io.WriteString(out, result.Inspect()+"\n")

		if err != nil {
			panic(err)
		}
	}
}

func printParseErrors(out io.Writer, errors []string) error {
	messages := []string{
		"Woops! We ran into some monkey business here!\n",
		" parser errors:\n\t",
		strings.Join(errors, "\n\t"),
		"\n",
	}
	_, err := io.WriteString(out, strings.Join(messages, ""))
	return err
}
