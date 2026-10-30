package main

import (
	"log"

	"todo-backend/internal/repository"
	"todo-backend/internal/server"

	"todo-backend/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("invalid config: %v", err)
	}

	redisRepository,_ := repository.NewRedisTodoRepository(cfg.RedisAddr())

	err = redisRepository.InitializeFromJSON("data/todos.json")
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	srv, err  := server.NewServer()
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
	log.Printf("Server listening on %s", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}