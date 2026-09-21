package repository

import "todo-backend/internal/models"

type TodoRepository interface {
	FetchAllTodos() ([]models.Todo, error)
	FetchTodo(id int) (models.Todo, error)
	CreateTodo(todo models.Todo) (models.Todo, error)
}