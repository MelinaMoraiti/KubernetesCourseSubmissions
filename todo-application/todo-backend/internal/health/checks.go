package health

import (
	"time"
	hg "github.com/hellofresh/health-go/v5"
	healthRedis "github.com/hellofresh/health-go/v5/checks/redis"
	healthNats "github.com/hellofresh/health-go/v5/checks/nats"
)

const defaultTimeout = 5 * time.Second

func Redis(dsn string) hg.Config {
	return hg.Config{
		Name:    "redis-check",
		Timeout: defaultTimeout,
		Check:   healthRedis.New(healthRedis.Config{DSN: dsn}),
	}
}

func Nats(dsn string) hg.Config {
	return hg.Config{
		Name: "nats-check",
		Timeout: defaultTimeout,
		Check:   healthNats.New(healthNats.Config{DSN: dsn}),
	}
}