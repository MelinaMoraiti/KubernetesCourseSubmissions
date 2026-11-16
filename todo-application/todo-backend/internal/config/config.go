package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port      int
	RedisHost string
	RedisPort string

	NatsURL string
	NatsStreamName  string
	NatsSubjects    []string
}

// RedisAddr returns host:port, for go-redis redis.Options{Addr: ...}.
func (c *Config) RedisAddr() string {
	return net.JoinHostPort(c.RedisHost, c.RedisPort)
}

// RedisDSN returns redis://host:port, for the health check.
func (c *Config) RedisDSN() string {
	return "redis://" + c.RedisAddr()
}

func Load() (*Config, error) {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return nil, fmt.Errorf("PORT is not set or invalid: %w", err)
	}

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		return nil, fmt.Errorf("REDIS_HOST is not set")
	}

	redisPort := os.Getenv("REDIS_PORT")
	if _, err := strconv.Atoi(redisPort); err != nil {
		return nil, fmt.Errorf("REDIS_PORT is not set or invalid: %w", err)
	}

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		return nil, fmt.Errorf("NATS_URL is not set")
	}

	natsSubjects := splitCSV(envOrDefault("NATS_SUBJECTS", "todos.>"))
	if len(natsSubjects) == 0 {
		return nil, fmt.Errorf("NATS_SUBJECTS is empty")
	}

	return &Config{
		Port:      port,
		RedisHost: redisHost,
		RedisPort: redisPort,
		NatsURL: natsURL,
		NatsStreamName: envOrDefault("NATS_STREAM_NAME", "TODOS"),
		NatsSubjects:   natsSubjects,
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
