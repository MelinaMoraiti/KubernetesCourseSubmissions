package messaging

import (
	"context"

	"broadcaster/internal/models"
)

type StatusPublisher interface {
	PublishStatusMessage(ctx context.Context, message models.StatusMessage) error
}
