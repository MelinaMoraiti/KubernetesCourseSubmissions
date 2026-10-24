package server

import (
	"fmt"
	"net/http"
	"time"

	hg "github.com/hellofresh/health-go/v5"
	"todo-backend/internal/config"
	"todo-backend/internal/health"

	"todo-backend/internal/repository"
	"todo-backend/internal/routes"
)

type Server struct {
	port int
}

func NewServer() (*http.Server, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	newServer := &Server{
		port: cfg.Port,
	}
	checks := []hg.Config{
		health.Redis(cfg.RedisDSN()),
		// later we can add more checks here ...

	}
	readyz, err := health.NewReadinessHandler(checks...)
	if err != nil {
		return nil, err
	}

	todoRepository,_ := repository.NewRedisTodoRepository(cfg.RedisAddr())

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", newServer.port),

		Handler: routes.RegisterRoutes(todoRepository, readyz),

		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server, nil
}