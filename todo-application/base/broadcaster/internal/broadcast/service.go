// Package broadcast turns todo events into chat messages.
package broadcast

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"broadcaster/internal/messaging"
	"broadcaster/internal/models"
	"broadcaster/internal/notifier"
)

const sendTimeout = 10 * time.Second

type Service struct {
	notifier notifier.Notifier
	status   messaging.StatusPublisher
	botName  string
}

func NewService(n notifier.Notifier, status messaging.StatusPublisher, botName string) *Service {
	return &Service{notifier: n, status: status, botName: botName}
}

// HandleTodoEvent matches messaging.TodoEventHandler.
func (s *Service) HandleTodoEvent(ctx context.Context, ev models.TodoEvent) error {
	msg := models.StatusMessage{User: s.botName, Message: Describe(ev)}

	// Detach from the app context so a message that is already being sent
	// can finish while the service shuts down.
	sendCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), sendTimeout)
	defer cancel()

	if err := s.notifier.Send(sendCtx, msg); err != nil {
		return fmt.Errorf("forward event %s: %w", ev.ID, err)
	}
	slog.Info("event forwarded", "event_id", ev.ID, "type", ev.Type)

	// Best effort: record what was sent on todos.status.
	if err := s.status.PublishStatusMessage(sendCtx, msg); err != nil {
		slog.Warn("could not publish status message", "event_id", ev.ID, "err", err)
	}
	return nil
}

func Describe(ev models.TodoEvent) string {
	switch ev.Type {
	case models.EventTypeCreated:
		return fmt.Sprintf("A todo was created: %q", ev.Todo.Task)
	case models.EventTypeCompleted:
		return fmt.Sprintf("A todo was completed: %q", ev.Todo.Task)
	default:
		return fmt.Sprintf("A todo was updated: %q", ev.Todo.Task)
	}
}
