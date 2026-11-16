package utils
import (
	"os"
	"encoding/json"
	"fmt"
	"errors"
	"todo-backend/internal/models"
)

func FetchAllTodosFromJSON(filePath string) ([]models.Todo, error) {
    file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open json file: %w", err)
	}
	defer file.Close()

	var todos []models.Todo
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&todos); err != nil {
		return nil, fmt.Errorf("failed to decode json data: %w", err)
	}

	return todos, nil
}
func AppendNewTodoInJSON(filename string, newTodo models.Todo) ([]models.Todo, error) {
    var todos []models.Todo

    data, err := os.ReadFile(filename)
    if err != nil {
        if !errors.Is(err, os.ErrNotExist) {
            return nil, err
        }

        todos = []models.Todo{}
    } else if len(data) > 0 {
        if err := json.Unmarshal(data, &todos); err != nil {
            return nil, err
        }
    }

    // Generate ID
    newTodo.ID = 1

    if len(todos) > 0 {
        newTodo.ID = todos[len(todos)-1].ID + 1
    }

    todos = append(todos, newTodo)

    data, err = json.MarshalIndent(todos, "", "  ")
    if err != nil {
        return nil, err
    }

    if err := os.WriteFile(filename, data, 0644); err != nil {
        return nil, err
    }

    return todos, nil
}