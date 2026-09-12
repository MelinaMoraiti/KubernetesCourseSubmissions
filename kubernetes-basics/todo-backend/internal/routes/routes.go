package routes

import (
	"log"
	"net/http"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"todo-backend/internal/utils"
	"todo-backend/internal/models"
	"github.com/go-chi/cors"
)

func RegisterRoutes() http.Handler {

	r := chi.NewRouter()
    r.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"http://localhost:8000"},
        AllowedMethods:   []string{"GET", "POST"},
        AllowedHeaders:   []string{"Accept", "Content-Type", "X-CSRF-Token"},
        AllowCredentials: true,
    }))
	r.Use(middleware.Logger)

	r.Get("/todos", readTodosHandler)
	r.Post("/todos", createTodosHandler)
	return r
}

func readTodosHandler(w http.ResponseWriter, r *http.Request) {
	todos, err := utils.FetchAllTodosFromJSON("data/todos.json")
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
func createTodosHandler(w http.ResponseWriter, r *http.Request) {
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
    _, err2 := utils.AppendNewTodoInJSON("data/todos.json", todo)
    if err2 != nil {
        log.Printf("Could not create todo: %v", err2)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("Todo created successfully"))
}