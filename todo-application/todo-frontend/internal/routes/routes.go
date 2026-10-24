package routes

import (
    "encoding/json"
    "html/template"
    "io"
    "log"
    "net/http"
    "net/url"
    "os"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"

    "todo-app/internal/models"
)

func RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Handle("/images/*",
        http.StripPrefix("/images/",
            http.FileServer(http.Dir(os.Getenv("STORED_IMAGES_PATH"))),
        ),
    )

	r.Get("/", indexHandler)
	r.Post("/todos", createTodoHandler)
	return r
}
var todoBackendURL = os.Getenv("TODO_BACKEND_URL")

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if todoBackendURL == "" {
		log.Println("TODO_BACKEND_URL is not set")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resp, err := http.Get(todoBackendURL + "/todos")
	if err != nil {
		log.Printf("Could not get all todos: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Todo API returned status: %s", resp.Status)
		http.Error(w, "Could not get todos", http.StatusInternalServerError)
		return
	}

	var todos []models.Todo
	if err := json.NewDecoder(resp.Body).Decode(&todos); err != nil {
		log.Printf("Could not decode todos: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		log.Printf("Could not parse template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, todos); err != nil {
		log.Printf("Could not execute template: %v", err)
	}
}
func createTodoHandler(w http.ResponseWriter, r *http.Request) {
	todo := r.FormValue("todo")

	form := url.Values{}
	form.Set("todo", todo)

	resp, err := http.PostForm(
		todoBackendURL+"/todos",
		form,
	)
	if err != nil {
		http.Error(w, "Could not contact todo backend", http.StatusBadGateway)
		return
	}

	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}