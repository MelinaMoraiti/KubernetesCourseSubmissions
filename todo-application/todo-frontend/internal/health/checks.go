package health

import (
	"context"
	"errors"
	"sync/atomic"
	"time"

	hg "github.com/hellofresh/health-go/v5"
	healthHttp "github.com/hellofresh/health-go/v5/checks/http"
)
type Switch struct {
	broken atomic.Bool
}
const defaultTimeout = 5 * time.Second
// Backend verifies the frontend can reach the todo backend and that the
// backend itself reports ready (its /readyz also covers Redis).
func Backend(baseURL string) hg.Config {
	return hg.Config{
		Name:    "todo-backend-check",
		Timeout: defaultTimeout ,
		Check: healthHttp.New(healthHttp.Config{
			URL:            baseURL,
			RequestTimeout: 2 * time.Second,
		}),
	}
}
// Break makes every handler that includes Check() fail until the process restarts.
func (s *Switch) Break() { s.broken.Store(true) }

func (s *Switch) Broken() bool { return s.broken.Load() }

func (s *Switch) Check() hg.Config {
	return hg.Config{
		Name:    "manual-break",
		Timeout: 2 * time.Second,
		Check: func(ctx context.Context) error {
			if s.broken.Load() {
				return errors.New("app was manually broken")
			}
			return nil
		},
	}
}