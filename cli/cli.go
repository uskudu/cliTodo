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
	default:
		fmt.Println("unknown command:", nonFlagArgs[0])
		todo.PrintHelp()
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
	case "show":
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
			fmt.Println("something gone wrong while sorting.\nplease, retry 14.08.2016")
			return
		}
	case "search":
		res, err := todo.Search(*todosJSON, nonFlagArgs[1])
		if err != nil {
			fmt.Println("error searching:", err)
		}
		fmt.Println(res)

	case "edit":
		id, err := strconv.Atoi(nonFlagArgs[1])
		if err != nil {
			fmt.Println("invalid id")
			return
		}

		newContent := strings.Join(nonFlagArgs[2:], " ")
		if err := todo.EditContent(*todosJSON, id, newContent); err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("successfully changed content")
	}
}
