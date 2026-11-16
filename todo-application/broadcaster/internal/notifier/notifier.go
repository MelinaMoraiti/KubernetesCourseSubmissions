// Package notifier sends status messages to an external chat service.
package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"broadcaster/internal/models"
)

const discordMaxContent = 2000

type Notifier interface {
	Send(ctx context.Context, msg models.StatusMessage) error
}

type webhook struct {
	url     string
	client  *http.Client
	payload func(models.StatusMessage) any
}

// New returns a Notifier for the given format: "generic" or "discord".
func New(format, url string) (Notifier, error) {
	if url == "" {
		return nil, fmt.Errorf("webhook URL is empty")
	}

	var payload func(models.StatusMessage) any
	switch strings.ToLower(format) {
	case "", "generic":
		// {"user": "...", "message": "..."}
		payload = func(m models.StatusMessage) any { return m }
	case "discord":
		payload = func(m models.StatusMessage) any {
			content := m.Message
			if r := []rune(content); len(r) > discordMaxContent {
				content = string(r[:discordMaxContent])
			}
			return map[string]any{
				"username": m.User,
				"content":  content,
				"allowed_mentions": map[string]any{"parse": []string{}},
			}
		}
	default:
		return nil, fmt.Errorf("unknown WEBHOOK_FORMAT %q (use \"generic\" or \"discord\")", format)
	}

	return &webhook{
		url:     url,
		client:  &http.Client{Timeout: 5 * time.Second},
		payload: payload,
	}, nil
}

// Send posts the message once. It deliberately does not retry: a retry after
// a timeout could deliver the same message twice, and duplicates are worse
// than a missed message here.
func (w *webhook) Send(ctx context.Context, msg models.StatusMessage) error {
	body, err := json.Marshal(w.payload(msg))
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		snippet, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("webhook returned %d: %s", resp.StatusCode, strings.TrimSpace(string(snippet)))
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	return nil
}
