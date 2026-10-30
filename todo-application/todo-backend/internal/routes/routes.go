package routes

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"todo-backend/internal/messaging"
	"todo-backend/internal/models"
	"todo-backend/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RegisterRoutes(todoRepository repository.TodoRepository, readyz http.Handler, publisher messaging.EventPublisher) http.Handler {
	r := chi.NewRouter()
	r.Handle("/readyz", readyz)
	r.Group(func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		r.Get("/todos", func(w http.ResponseWriter, r *http.Request) {
			readTodosHandler(w, r, todoRepository)
		})
		r.Post("/todos", func(w http.ResponseWriter, r *http.Request) {
			createTodosHandler(w, r, todoRepository, publisher)
		})
		r.Put("/todos/{id}", func(w http.ResponseWriter, r *http.Request) {
    		completeTodoHandler(w, r, todoRepository, publisher)
		})
	})
	return r
}
func completeTodoHandler(w http.ResponseWriter, r *http.Request, todoRepository repository.TodoRepository, publisher messaging.EventPublisher) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	todo, err := todoRepository.CompleteTodo(id)
	if err != nil {
		if errors.Is(err, repository.ErrTodoNotFound) {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		log.Printf("Could not complete todo %d: %v", id, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
  
    event := models.NewTodoEvent(models.EventTypeCompleted, todo)

    if err := publisher.PublishTodoEvent(r.Context(), event); err != nil {
        log.Printf("Could not publish TodoCompleted event: %v", err)
        http.Error(w, "Failed to publish todo event", http.StatusInternalServerError)
        return
    }

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(todo); err != nil {
		log.Printf("Could not encode todo %d: %v", id, err)
	}
}
func readTodosHandler(w http.ResponseWriter, r *http.Request, todoRepository repository.TodoRepository) {
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

func createTodosHandler(w http.ResponseWriter, r *http.Request, todoRepository repository.TodoRepository, publisher messaging.EventPublisher) {
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
	createdTodo, err2 := todoRepository.CreateTodo(todo)
	if err2 != nil {
		log.Printf("Could not create todo: %v", err2)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	
    event := models.NewTodoEvent(models.EventTypeCreated, createdTodo)

    if err := publisher.PublishTodoEvent(r.Context(), event); err != nil {
        log.Printf("Could not publish TodoCreated event: %v", err)
        http.Error(w, "Failed to publish todo event", http.StatusInternalServerError)
        return
    }

	log.Printf("Todo: %s created successfully", createdTodo.Task)

	w.WriteHeader(http.StatusCreated)
	//w.Write([]byte("Todo created successfully"))
}
