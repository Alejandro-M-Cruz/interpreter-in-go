# interpreter-in-go

Interpreter for a simple scripting language, called Mandrill, based on the Monkey programming language from the book <a href="https://interpreterbook.com" target="_blank">_Writing an Interpreter in Go_</a> by Thorsten Ball. 

## Language overview 

The language supports integers, booleans (`true`, `false`) and strings. Semicolons are optional. 

It includes assignment (`let`) and return statements, while everything else is considered an expression, including if/else.

```
>> let x = 2
>> let y = 3
>> let min = if (x < y) { x } else { y } 
>> min
2
>> "Hello" + ", " + "world!";
Hello, world!
```

Lastly, the language has first-class functions and implicit return, and it fully supports closures.

```
>> let newAdder = fn(x) { fn(n) { x + n } }
>> let addThree = newAdder(3)
>> addThree(4)
7
```



## How to run

To run the interpreter, follow these steps:

1. Ensure that Go is installed on your system. You can download and install it from the <a href="https://go.dev/dl" target="_blank">official Go website</a>.

1. Clone this repository to your local machine.

1. Open your terminal and navigate to the project directory.

1. Execute the following command to start the Read-Eval-Print Loop (REPL):

```
go run .
```
The REPL evaluates your input and displays the result of expressions.
