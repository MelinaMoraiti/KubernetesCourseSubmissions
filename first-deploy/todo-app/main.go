package main

import (
	"net/http"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	_ = godotenv.Load()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" 
	}
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("Server started in port %s",port)))
	})
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:"+port), r)
}