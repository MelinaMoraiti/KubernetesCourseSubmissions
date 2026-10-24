package main

import (
	"log"

	"todo-backend/internal/server"
	"todo-backend/internal/repository"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	redisRepository,_ := repository.NewRedisTodoRepository()

	err := redisRepository.InitializeFromJSON("data/todos.json")
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	srv := server.NewServer()
	log.Printf("Server listening on %s", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}