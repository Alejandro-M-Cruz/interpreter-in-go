package main

import (
	"example.com/writing-an-interpreter/repl"
	"fmt"
	"os"
	"os/user"
)

func main() {
	u, err := user.Current()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Hi, %s! This is the Monkey Programming Language.\n", u.Username)
	fmt.Println("Feel free to type in commands...")

	repl.Start(os.Stdin, os.Stdout)
}
