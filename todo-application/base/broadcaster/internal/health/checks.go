package health

import (
	hg "github.com/hellofresh/health-go/v5"
	healthNats "github.com/hellofresh/health-go/v5/checks/nats"
	"time"
)

const defaultTimeout = 5 * time.Second

func Nats(dsn string) hg.Config {
	return hg.Config{
		Name:    "nats-check",
		Timeout: defaultTimeout,
		Check:   healthNats.New(healthNats.Config{DSN: dsn}),
	}
}
