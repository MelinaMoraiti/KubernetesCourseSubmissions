package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
 
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
 
	"todo-backend/internal/config"
	"todo-backend/internal/models"
)


type NatsService struct {
	js jetstream.JetStream
	nc *nats.Conn
}


var _ EventPublisher = (*NatsService)(nil)

// NewNatsService connects to NATS and makes sure the stream exists.
// The NATS server must be started with JetStream enabled (nats-server -js).
// It keeps retrying until ctx expires, so the app can start before NATS is ready.

func NewNatsService(ctx context.Context, cfg *config.Config) (*NatsService, error) {
	nc, err := nats.Connect(cfg.NatsURL,
		nats.MaxReconnects(-1),
		nats.RetryOnFailedConnect(true),
	)
	if err != nil {
		return nil, fmt.Errorf("connect to nats: %w", err)
	}
 
	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("init jetstream: %w", err)
	}
 
	streamCfg := jetstream.StreamConfig{
		Name:     cfg.NatsStreamName,
		Subjects: cfg.NatsSubjects,
		Storage:  jetstream.FileStorage,
		MaxAge:   7 * 24 * time.Hour,
	}
	for {
		_, err = js.CreateOrUpdateStream(ctx, streamCfg)
		if err == nil {
			break
		}
		select {
		case <-ctx.Done():
			nc.Close()
			return nil, fmt.Errorf("create stream %q: %w (last error: %v)", cfg.NatsStreamName, ctx.Err(), err)
		case <-time.After(time.Second):
		}
	}
 
	return &NatsService{nc: nc, js: js}, nil
}
// PublishTodoEvent publishes and waits for the JetStream acknowledgement.
// The event ID is used as Nats-Msg-Id so a retried publish isn't stored twice.
func (s *NatsService) PublishTodoEvent(ctx context.Context, event models.TodoEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}
 
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
 
	if _, err := s.js.Publish(ctx, event.Subject(), data, jetstream.WithMsgID(event.ID)); err != nil {
		return fmt.Errorf("publish %s: %w", event.Subject(), err)
	}
	return nil
}
 
// Close drains the connection so buffered messages are flushed.
func (s *NatsService) Close() error {
	return s.nc.Drain()
}
 
