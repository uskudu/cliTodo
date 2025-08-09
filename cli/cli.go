package cli

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"cliTodo/todo"
)

func Run(args []string, defaultFile string) {
	todosJSON := flag.String("file", defaultFile, "json file of todos")

	err := flag.CommandLine.Parse(args)
	if err != nil {
		return
	}

	nonFlagArgs := flag.CommandLine.Args()
	if len(nonFlagArgs) == 0 {
		return
	}

	switch nonFlagArgs[0] {
	case "add":
		err := todo.Create(*todosJSON, strings.Join(nonFlagArgs[1:], " "))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("new todo added")

	case "list":
		content := todo.ShowJSON(*todosJSON)
		fmt.Println(content)

	case "done":
		if len(nonFlagArgs) < 2 {
			fmt.Println("missing todo ID")
			os.Exit(1)
		}
		id, err := strconv.Atoi(nonFlagArgs[1])
		if err != nil {
			fmt.Println("invalid id")
			os.Exit(1)
		}
		err = todo.MarkDone(*todosJSON, id)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("todo %v is now marked done\n", id)

	case "del":
		if len(nonFlagArgs) < 2 {
			fmt.Println("missing todo ID")
			os.Exit(1)
		}
		id, err := strconv.Atoi(nonFlagArgs[1])
		if err != nil {
			fmt.Println("invalid id")
			os.Exit(1)
		}
		err = todo.Delete(*todosJSON, id)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("deleted todo %v\n", id)

	case "help":
		todo.PrintHelp()
	case "sortby":
		err := todo.SortFile(*todosJSON, nonFlagArgs[1])
		if err != nil {
			return
		}
	default:
		fmt.Println("unknown command:", nonFlagArgs[0])
		todo.PrintHelp()
	}
}
