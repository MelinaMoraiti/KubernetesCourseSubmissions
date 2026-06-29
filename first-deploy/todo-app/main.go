package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	port := os.Getenv("PORT")


	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		fmt.Fprint(w, `<!DOCTYPE html>
        <html>
        <head>
          <title>To-Do App</title>
        </head>
        <body>
          <h2>To-Do List</h2>
          <p>Server running on port `+port+`</p>
        </body>
        </html>`)
	})

	http.ListenAndServe("0.0.0.0:"+port, r)
}