package main

import "os"
import "cliTodo/cli"
import "cliTodo/todo"

func main() {
	cli.Run(os.Args[1:], todo.Todos_file)
}
