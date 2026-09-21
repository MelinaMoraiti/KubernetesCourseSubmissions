package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"todo-backend/internal/repository"
	"todo-backend/internal/routes"
)

type Server struct {
	port int
}

func NewServer() *http.Server {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic("PORT is not set or invalid")
	}

	newServer := &Server{
		port: port,
	}

	todoRepository,_ := repository.NewRedisTodoRepository()

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", newServer.port),

		Handler: routes.RegisterRoutes(todoRepository),

		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}