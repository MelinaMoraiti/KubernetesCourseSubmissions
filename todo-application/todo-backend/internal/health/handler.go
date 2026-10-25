package health

import (
	"fmt"
	"net/http"

	hg "github.com/hellofresh/health-go/v5"
)

// NewReadinessHandler returns an http.Handler that runs all given checks.
func NewHandler(checks ...hg.Config) (http.Handler, error) {
	h, err := hg.New(hg.WithSystemInfo())
	if err != nil {
		return nil, fmt.Errorf("create health checker: %w", err)
	}

	for _, c := range checks {
		if err := h.Register(c); err != nil {
			return nil, fmt.Errorf("register %q health check: %w", c.Name, err)
		}
	}

	return h.Handler(), nil
}