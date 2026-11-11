package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port int

	NatsURL        string
	NatsStreamName string
	NatsSubjects   []string

	// External service the events are forwarded to.
	WebhookURL    string
	WebhookFormat string // "generic" or "discord"
	BotName       string
}

func Load() (*Config, error) {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return nil, fmt.Errorf("PORT is not set or invalid: %w", err)
	}

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		return nil, fmt.Errorf("NATS_URL is not set")
	}

	natsSubjects := splitCSV(envOrDefault("NATS_SUBJECTS", "todos.>"))
	if len(natsSubjects) == 0 {
		return nil, fmt.Errorf("NATS_SUBJECTS is empty")
	}

	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		return nil, fmt.Errorf("WEBHOOK_URL is not set")
	}

	return &Config{
		Port:           port,
		NatsURL:        natsURL,
		NatsStreamName: envOrDefault("NATS_STREAM_NAME", "TODOS"),
		NatsSubjects:   natsSubjects,
		WebhookURL:     webhookURL,
		WebhookFormat:  strings.ToLower(envOrDefault("WEBHOOK_FORMAT", "generic")),
		BotName:        envOrDefault("BOT_NAME", "bot"),
	}, nil
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func splitCSV(s string) []string {
	var out []string
	for _, part := range strings.Split(s, ",") {
		if p := strings.TrimSpace(part); p != "" {
			out = append(out, p)
		}
	}
	return out
}
