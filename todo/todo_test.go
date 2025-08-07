package todo

import (
	"os"
	"testing"
)

func TestCreate(t *testing.T) {
	// Создаём временный файл
	tmpFile, err := os.CreateTemp("", "todos_test_*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name()) // очистка после теста
	tmpFile.Close()

	// Изначально записываем пустой массив, чтобы Load корректно прочитал
	err = os.WriteFile(tmpFile.Name(), []byte("[]"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Добавляем новую задачу
	err = Create(tmpFile.Name(), "test task")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Загружаем задачи обратно
	todos := Load(tmpFile.Name())
	if len(todos) != 1 {
		t.Fatalf("expected 1 todo, got %d", len(todos))
	}

	if todos[0].Content != "test task" {
		t.Errorf("expected content 'test task', got '%s'", todos[0].Content)
	}

	if todos[0].Done != false {
		t.Errorf("expected Done false, got true")
	}

	if todos[0].ID != 1 {
		t.Errorf("expected ID 1, got %d", todos[0].ID)
	}
}
