package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"todo-backend/internal/config"
	"todo-backend/internal/health"
	"todo-backend/internal/messaging"

	hg "github.com/hellofresh/health-go/v5"

	"todo-backend/internal/repository"
	"todo-backend/internal/routes"
)

type Server struct {
	port int
}
func NewServer() (*http.Server, *messaging.NatsService, error) {
    cfg, err := config.Load()
    if err != nil {
        return nil, nil, err
    }

    checks := []hg.Config{
        health.Redis(cfg.RedisDSN()),
        health.Nats(cfg.NatsURL),
    }
    readyz, err := health.NewHandler(checks...)
    if err != nil {
        return nil, nil, err
    }

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel() 
    natsService, err := messaging.NewNatsService(ctx, cfg)
    if err != nil {
        return nil, nil, err
    }

    todoRepository, err := repository.NewRedisTodoRepository(cfg.RedisAddr())
    if err != nil {
        natsService.Close()
        return nil, nil, err
    }

    srv := &http.Server{
        Addr:         fmt.Sprintf(":%d", cfg.Port),
        Handler:      routes.RegisterRoutes(todoRepository, readyz, natsService),
        IdleTimeout:  time.Minute,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 30 * time.Second,
    }
    return srv, natsService, nil
}