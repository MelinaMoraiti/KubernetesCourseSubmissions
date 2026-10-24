package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

type Config struct {
	Port      int
	RedisHost string
	RedisPort string
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

	return &Config{
		Port:      port,
		RedisHost: redisHost,
		RedisPort: redisPort,
	}, nil
}