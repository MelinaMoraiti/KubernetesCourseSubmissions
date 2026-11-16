package messaging

import (
    "context"

    "todo-backend/internal/models"
)

type EventPublisher interface {
    PublishTodoEvent(ctx context.Context, event models.TodoEvent) error
}