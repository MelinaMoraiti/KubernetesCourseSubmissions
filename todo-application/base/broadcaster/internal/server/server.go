package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"broadcaster/internal/config"
	"broadcaster/internal/health"
	"broadcaster/internal/messaging"
)

// NewServer builds the HTTP server (health endpoints) and connects to NATS.
func NewServer(cfg *config.Config) (*http.Server, *messaging.NatsService, error) {
	readyz, err := health.NewHandler(health.Nats(cfg.NatsURL))
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	natsService, err := messaging.NewNatsService(ctx, cfg)
	if err != nil {
		return nil, nil, err
	}

	mux := http.NewServeMux()
	// Liveness: the process is up. Must not depend on NATS.
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	// Readiness: NATS is reachable.
	mux.Handle("/readyz", readyz)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return srv, natsService, nil
}
