package health

import (
	"time"
	hg "github.com/hellofresh/health-go/v5"
	healthRedis "github.com/hellofresh/health-go/v5/checks/redis"
)

const defaultTimeout = 5 * time.Second

func Redis(dsn string) hg.Config {
	return hg.Config{
		Name:    "redis-check",
		Timeout: defaultTimeout,
		Check:   healthRedis.New(healthRedis.Config{DSN: dsn}),
	}
}