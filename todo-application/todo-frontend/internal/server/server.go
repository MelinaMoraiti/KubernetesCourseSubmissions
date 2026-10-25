package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"todo-app/internal/config"
	"todo-app/internal/routes"
)

type Server struct {
	port int
}

func NewServer() *http.Server {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("invalid config: %v", err)
	}

	newServer := &Server{
		port: cfg.Port,
	}
	handler, err := routes.RegisterRoutes(cfg.TodoBackendURL, cfg.StoredImagesPath)
	if err != nil {
		log.Fatalf("failed to register routes: %v", err)
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", newServer.port),
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}