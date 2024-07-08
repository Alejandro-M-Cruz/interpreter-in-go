package main

import (
	"example.com/writing-an-interpreter/repl"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"os/user"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	u, err := user.Current()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Hi, %s! This is the Mandrill programming language.\n", u.Username)
	fmt.Println("Feel free to type in commands...")

	repl.Start(os.Stdin, os.Stdout)
}
