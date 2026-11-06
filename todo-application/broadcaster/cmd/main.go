package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"broadcaster/internal/broadcast"
	"broadcaster/internal/config"
	"broadcaster/internal/notifier"
	"broadcaster/internal/server"
)

func main() {
	if err := run(); err != nil {
		slog.Error("broadcaster failed", "err", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv, natsService, err := server.NewServer(cfg)
	if err != nil {
		return err
	}
	defer natsService.Close()

	webhook, err := notifier.New(cfg.WebhookFormat, cfg.WebhookURL)
	if err != nil {
		return err
	}
	svc := broadcast.NewService(webhook, natsService, cfg.BotName)

	if err := natsService.SubscribeTodoEvents(ctx, svc.HandleTodoEvent); err != nil {
		return err
	}

	serveErr := make(chan error, 1)
	go func() {
		slog.Info("broadcaster started", "addr", srv.Addr, "format", cfg.WebhookFormat)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serveErr <- err
		}
	}()

	select {
	case <-ctx.Done():
		slog.Info("shutting down")
	case err := <-serveErr:
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
