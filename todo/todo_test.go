package todo

import (
	"encoding/json"
	//"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var testData = `[
        {
            "id": 1,
            "content": "test",
            "done": false,
            "createdAt": "14.08.2025 02:30:42",
            "doneAt": "14.08.2025 02:31:04",
            "duration": "22.193747377s"
        }
    ]`

var testDataBig = `[
 {
  "id": 1,
  "content": "test_1",
  "done": true,
  "createdAt": "14.08.2025 02:30:42",
  "doneAt": "10.08.2025 02:31:04",
  "duration": "16.193747377s"
 },
 {
  "id": 2,
  "content": "test_2_even",
  "done": false,
  "createdAt": "14.08.2025 02:30:42",
  "doneAt": "11.08.2025 02:31:04",
  "duration": "64.193747377s"
 },
 {
  "id": 3,
  "content": "test_3",
  "done": true,
  "createdAt": "14.08.2025 02:30:42",
  "doneAt": "12.08.2025 02:31:04",
  "duration": "32.193747377s"
 },
 {
  "id": 4,
  "content": "test_4_even",
  "done": false,
  "createdAt": "14.08.2025 02:30:42",
  "doneAt": "13.08.2025 02:31:04",
  "duration": "8.193747377s"
 }
]`

func TestCreate(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "todos_test_*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	err = os.WriteFile(tmpFile.Name(), []byte("[]"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = Create(tmpFile.Name(), "test todo")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	todos, err := LoadTodos(tmpFile.Name())
	if len(todos) != 1 {
		t.Fatalf("expected 1 todo, got %d", len(todos))
	}

	if todos[0].Content != "test todo" {
		t.Errorf("expected content 'test todo', got '%s'", todos[0].Content)
	}

	if todos[0].Done != false {
		t.Errorf("expected Done false, got true")
	}

	if todos[0].ID != 1 {
		t.Errorf("expected ID 1, got %d", todos[0].ID)
	}
}

func TestShowJSON(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "todos_test_*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	err = os.WriteFile(tmpFile.Name(), []byte(testData), 0644)
	if err != nil {
		t.Fatal(err)
	}

	result := ShowJSON(tmpFile.Name())
	if result != testData {
		t.Errorf("Expected:\n%s\n\nGot:\n%s", testData, result)
	}
}

func TestMarkDone(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "todos_test_*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	err = os.WriteFile(tmpFile.Name(), []byte("[]"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(tmpFile.Name(), []byte(testData), 0644)
	if err != nil {
		t.Fatal(err)
	}

	todos, err := LoadTodos(tmpFile.Name())

	err = MarkDone(tmpFile.Name(), 1)
	if err != nil {
		t.Errorf("Error while marking as done: %s", err)
	}
	todos, err = LoadTodos(tmpFile.Name())

	if todos[0].Done != true {
		t.Errorf("Expected true, got false.\ntrue was not set after MarkDone called.")
	}
}

func TestDelete(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "todos_test_*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	err = os.WriteFile(tmpFile.Name(), []byte("[]"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(tmpFile.Name(), []byte(testData), 0644)
	if err != nil {
		t.Fatal(err)
	}

	todos, err := LoadTodos(tmpFile.Name())
	if len(todos) != 1 {
		t.Errorf("Expected len == 1, got %d", len(todos))
	}
	err = Delete(tmpFile.Name(), 1)
	if err != nil {
		t.Errorf("Error while deleting: %s", err)
	}

	todos, err = LoadTodos(tmpFile.Name())
	if len(todos) != 0 {
		t.Errorf("Expected len == 0, got %d", len(todos))
	}
}

func TestSortTodos(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "todos_test_*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	err = os.WriteFile(tmpFile.Name(), []byte("[]"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(tmpFile.Name(), []byte(testDataBig), 0644)
	if err != nil {
		t.Fatal(err)
	}

	todos, err := LoadTodos(tmpFile.Name())

	sortTodos(todos, "done")
	if todos[0].Done != false || todos[3].Done != true {
		t.Errorf("Failed while sorting by done")
	}

	sortTodos(todos, "duration")
	if !(todos[0].Duration < todos[1].Duration &&
		todos[1].Duration < todos[2].Duration &&
		todos[2].Duration < todos[3].Duration) {
		t.Errorf("Failed while sorting by duration")
	}
	sortTodos(todos, "id")
	if todos[0].ID != 1 || todos[3].ID != 4 {
		t.Errorf("Failed while sorting by id")
	}
}

func TestSearch(t *testing.T) {
	var mock_exp = `[
	 {
		"id": 2,
		"content": "test_2_even",
		"done": false,
		"createdAt": "14.08.2025 02:30:42",
		"doneAt": "11.08.2025 02:31:04",
		"duration": "1m4.193747377s"
	 },
	 {
		"id": 4,
		"content": "test_4_even",
		"done": false,
		"createdAt": "14.08.2025 02:30:42",
		"doneAt": "13.08.2025 02:31:04",
		"duration": "8.193747377s"
	 }
	]`

	tmpFile, err := os.CreateTemp("", "todos_test_*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	err = os.WriteFile(tmpFile.Name(), []byte("[]"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(tmpFile.Name(), []byte(testDataBig), 0644)
	if err != nil {
		t.Fatal(err)
	}

	result, err := Search(tmpFile.Name(), "even")
	if err != nil {
		t.Errorf("Error while searching: %s", err)
	}

	normalize := func(s string) string {
		var v interface{}
		if err := json.Unmarshal([]byte(s), &v); err != nil {
			return s
		}
		b, _ := json.MarshalIndent(v, "", "    ")
		return string(b)
	}

	if normalize(result) != normalize(mock_exp) {
		t.Errorf("Search failed. Expected:\n%s\n\nGot:\n%s", normalize(mock_exp), normalize(result))
	}
}

func TestEditContent(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "todos_test_*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	err = os.WriteFile(tmpFile.Name(), []byte("[]"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(tmpFile.Name(), []byte(testData), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = EditContent(tmpFile.Name(), 1, "new_content_test")
	if err != nil {
		t.Fatal(err)
	}

	todos, err := LoadTodos(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	if todos[0].Content != "new_content_test" {
		t.Errorf("Edit failed: content didnt change")
	}
}
