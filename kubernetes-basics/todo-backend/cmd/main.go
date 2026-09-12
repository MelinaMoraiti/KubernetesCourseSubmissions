package main

import (
	"log"

	"todo-backend/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	srv := server.NewServer()
	log.Printf("Server listening on %s", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}