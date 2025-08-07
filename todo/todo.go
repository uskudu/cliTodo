package todo

import (
	"encoding/json"
	"fmt"
	"os"
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
	data, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	var todos []Todo
	if len(data) > 0 {
		if err := json.Unmarshal(data, &todos); err != nil {
			return err
		}
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
	data, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	var todos []Todo
	if len(data) > 0 {
		if err := json.Unmarshal(data, &todos); err != nil {
			return err
		}
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

func PrintHelp() {
	fmt.Println("This is a todo CLI application.")
	fmt.Println("Todos are stored in json file.")
	fmt.Println("Available commands: add, list, done, del.")
	fmt.Println("add: adds a new todo to your file.")
	fmt.Println("example: go run main.go add buy tomatoes and potatoes")
	fmt.Println("list: shows your todos in json format.")
	fmt.Println("example: go run main.go list")
	fmt.Println("done: marks todo as done (specify id of todo).")
	fmt.Println("example: go run main.go done 14")
	fmt.Println("del: deletes todo (specify id of todo).")
	fmt.Println("example: go run main.go del 88")

}

// func Search(fileName string, search string) "search by word for included in content"
// func EditContent(fileName string, todoID int, content string) "search for included in content"
// func SortFile(fileName string, sortBy string) sort file by id/done/createdAt
// func Stats(fileName string) returns avg duration, avg created in the morning, avg created/done in one day
