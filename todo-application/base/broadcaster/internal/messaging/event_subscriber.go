package messaging

import (
	"context"

	"broadcaster/internal/models"
)

type TodoEventHandler func(ctx context.Context, event models.TodoEvent) error

type EventSubscriber interface {
	SubscribeTodoEvents(ctx context.Context, handler TodoEventHandler) error
}
