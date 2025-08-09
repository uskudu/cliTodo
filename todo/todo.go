package todo

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

const Todos_file = "my_todos.json"

type Todo struct {
	ID        int            `json:"id"`
	Content   string         `json:"content"`
	Done      bool           `json:"done"`
	CreatedAt MyTime         `json:"createdAt"`
	DoneAt    MyTime         `json:"doneAt"`
	Duration  DurationString `json:"duration"`
}

func PrintHelp() {
	fmt.Println("This is a todo CLI application.")
	fmt.Println("Todos are stored in json file.")
	fmt.Println("Available commands: add, list, done, del.")
	fmt.Println("")
	fmt.Println("add: adds a new todo to your file.")
	fmt.Println("example: go run main.go --file=BossTodos.json add buy tomatoes and potatoes")
	fmt.Println("")
	fmt.Println("list: shows your todos in json format.")
	fmt.Println("example: go run main.go list")
	fmt.Println("")
	fmt.Println("done: marks todo as done (specify id of todo).")
	fmt.Println("example: go run main.go done 14")
	fmt.Println("")
	fmt.Println("del: deletes todo (specify id of todo).")
	fmt.Println("example: go run main.go del 88")
	fmt.Println("")
	fmt.Println("sortby: sorts todos file by key provided")
	fmt.Println("example: go run main.go sortby done")
}

func JSONtoTodo(fileName string) ([]Todo, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var todos []Todo
	if len(data) > 0 {
		if err := json.Unmarshal(data, &todos); err != nil {
			return nil, err
		}
	}
	return todos, nil
}

func Load(fileName string) []Todo {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return []Todo{}
	}
	var todos []Todo
	if len(data) == 0 {
		return todos
	}
	err = json.Unmarshal(data, &todos)
	if err != nil {
		return []Todo{}
	}
	return todos
}

func Create(fileName string, content string) error {
	if len(content) == 0 {
		return nil
	}
	if len(content) > 100 {
		return fmt.Errorf("chill man")
	}

	todos := Load(fileName)
	for _, t := range todos {
		if t.Content == content {
			return fmt.Errorf("no duplications allowed")
		}
	}
	newID := 1
	for _, t := range todos {
		if t.ID >= newID {
			newID = t.ID + 1
		}
	}
	now := MyTime(time.Now())

	newTodo := Todo{
		ID:        newID,
		Content:   content,
		Done:      false,
		CreatedAt: now,
	}
	todos = append(todos, newTodo)

	data, err := json.MarshalIndent(todos, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, data, 0644)
}

func ShowJSON(fileName string) string {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return "{}"
	}
	return string(data)
}

func MarkDone(fileName string, todoID int) error {
	todos, err := JSONtoTodo(fileName)
	if err != nil {
		return err
	}

	found := false
	for i := range todos {
		if todos[i].ID == todoID {
			if todos[i].Done == true {
				return fmt.Errorf("todo #%v is already done", todoID)
			}
			now := MyTime(time.Now())
			todos[i].Done = true
			todos[i].DoneAt = now
			todos[i].Duration = DurationString(time.Time(now).Sub(time.Time(todos[i].CreatedAt)))
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("todo #%v not found", todoID)
	}

	newData, err := json.MarshalIndent(todos, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, newData, 0644)
}

func Delete(fileName string, todoID int) error {
	todos, err := JSONtoTodo(fileName)
	if err != nil {
		return err
	}

	found := false
	for i := range todos {
		if todos[i].ID == todoID {
			todos = append(todos[:i], todos[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("todo with id %v not found", todoID)
	}

	newData, err := json.MarshalIndent(todos, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, newData, 0644)
}

func sortTodos(todos []Todo, sortBy string) []Todo {
	sort.Slice(todos, func(i, j int) bool {
		switch sortBy {
		case "id":
			return todos[i].ID < todos[j].ID
		case "done":
			return !todos[i].Done && todos[j].Done
		case "duration":
			return todos[i].Duration < todos[j].Duration
		default:
			return false
		}
	})
	return todos
}

func SortFile(fileName string, sortBy string) error {
	todos, err := JSONtoTodo(fileName)
	if err != nil {
		return err
	}
	todos = sortTodos(todos, sortBy)
	newData, err := json.MarshalIndent(todos, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, newData, 0644)
}

// "search by word for included in content"
func Search(fileName string, search string) {
	todos, err := JSONtoTodo(fileName)
	if err != nil {
		return
	}

	result := []Todo{}
	for _, todo := range todos {
		if strings.Contains(todo.Content, search) {
			result = append(result, todo)
		}
	}
	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println("JSON coding error:", err)
		return
	}

	fmt.Println(string(jsonResult))
}

func EditContent(fileName string, todoID int, newContent string) error {
	todos, err := JSONtoTodo(fileName)
	if err != nil {
		return err
	}

	done := false
	for i := range todos {
		if todos[i].ID == todoID {
			todos[i].Content = newContent
			done = true
			break
		}
	}

	if !done {
		return fmt.Errorf("todo #%v is not found", todoID)
	}

	newData, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, newData, 0644)
}

// func Stats(fileName string) returns avg duration, avg created in the morning, avg created/done in one day
