package routes

import (
	"log"
	"net/http"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"todo-backend/internal/repository"
	"todo-backend/internal/models"
)

func RegisterRoutes(todoRepository repository.TodoRepository) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/todos", func(w http.ResponseWriter, r *http.Request) {
		readTodosHandler(w, r, todoRepository)
	})

	r.Post("/todos", func(w http.ResponseWriter, r *http.Request) {
		createTodosHandler(w, r, todoRepository)
	})

	return r
}

func readTodosHandler(w http.ResponseWriter, r *http.Request, todoRepository repository.TodoRepository,) {
	todos, err := todoRepository.FetchAllTodos()
	if err != nil {
		log.Printf("Could not get all todos: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Ensure an empty result is [] rather than null.
	if todos == nil {
		todos = make([]models.Todo, 0)
	}

	if err := json.NewEncoder(w).Encode(todos); err != nil {
		log.Printf("Could not encode todos: %v", err)
	}
}
func createTodosHandler(w http.ResponseWriter, r *http.Request,todoRepository repository.TodoRepository,) {
    err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
    // Get form value
    task := r.FormValue("todo")

    if task == "" {
        http.Error(w, "Task cannot be empty", http.StatusBadRequest)
        return
    }

    // Create new todo
    todo := models.Todo{
        Task: task,
        Done: false,
    }

    // Append todo to JSON file
    _, err2 := todoRepository.CreateTodo(todo)
    if err2 != nil {
        log.Printf("Could not create todo: %v", err2)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("Todo created successfully"))
}