package routes

import (
	"log"
	"net/http"
	"html/template"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"todo-app/internal/utils"
)

func RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Handle("/static/*",
        http.StripPrefix("/static/",
            http.FileServer(http.Dir("web/static")),
        ),
    )

	r.Get("/", indexHandler)
	return r
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	todos, err := utils.FetchAllTodosFromJSON("data/todos.json")
    if err != nil {
        log.Printf("Could not get all todos: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    tmpl, err := template.ParseFiles("web/templates/index.html")
    if err != nil {
        log.Printf("Could not parse template: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    // Execute the root template for the full index page render
    err = tmpl.Execute(w, todos)
    if err != nil {
        log.Printf("Could not execute template: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}
