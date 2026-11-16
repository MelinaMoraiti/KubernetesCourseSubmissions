package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	"broadcaster/internal/config"
	"broadcaster/internal/models"
)

const (
	BroadcasterConsumer = "BROADCASTER"
	StatusSubject       = models.SubjectTodoStatus
)

type NatsService struct {
	nc     *nats.Conn
	js     jetstream.JetStream
	stream string
	cc     jetstream.ConsumeContext
}

var (
	_ EventSubscriber = (*NatsService)(nil)
	_ StatusPublisher = (*NatsService)(nil)
)

func NewNatsService(ctx context.Context, cfg *config.Config) (*NatsService, error) {
	nc, err := nats.Connect(
		cfg.NatsURL,
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
		nats.RetryOnFailedConnect(true),
	)
	if err != nil {
		return nil, fmt.Errorf("connect to nats: %w", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("initialize jetstream: %w", err)
	}

	// Wait until the stream exists.
	for {
		_, err = js.Stream(ctx, cfg.NatsStreamName)
		if err == nil {
			break
		}

		select {
		case <-ctx.Done():
			nc.Close()
			return nil, fmt.Errorf("get stream %q: %w", cfg.NatsStreamName, ctx.Err())
		case <-time.After(time.Second):
		}
	}

	return &NatsService{
		nc:     nc,
		js:     js,
		stream: cfg.NatsStreamName,
	}, nil
}

// SubscribeTodoEvents attaches this replica to the shared durable consumer.
//
// Every replica uses the same durable name, so JetStream gives each message
// to exactly one replica. Delivery is at-most-once on purpose:
//   - MaxDeliver=1: a message is never redelivered, even if this replica
//     crashes, the ack is lost or the handler is slow.
//   - If the handler fails the message is terminated, not Nak'ed.
//   - DeliverNew: when the consumer is first created it does not replay the
//     stream's history.
//
// A message may be lost when something is broken, but it is never sent twice.
func (s *NatsService) SubscribeTodoEvents(ctx context.Context, handler TodoEventHandler) error {
	consumer, err := s.js.CreateOrUpdateConsumer(ctx, s.stream, jetstream.ConsumerConfig{
		Durable:       BroadcasterConsumer,
		AckPolicy:     jetstream.AckExplicitPolicy,
		DeliverPolicy: jetstream.DeliverNewPolicy,
		// The stream captures "todos.>", which includes todos.status. Only
		// consume the event subjects so we never read our own status messages.
		FilterSubjects: []string{
			models.SubjectTodoCreated,
			models.SubjectTodoCompleted,
		},
		AckWait:       30 * time.Second,
		MaxAckPending: 100,
		MaxDeliver:    1,
	})
	if err != nil {
		return fmt.Errorf("create broadcaster consumer: %w", err)
	}

	cc, err := consumer.Consume(
		func(msg jetstream.Msg) {
			var event models.TodoEvent
			if err := json.Unmarshal(msg.Data(), &event); err != nil {
				slog.Error("invalid todo event, dropping", "err", err)
				_ = msg.Term()
				return
			}

			if err := handler(ctx, event); err != nil {
				slog.Error("handling todo event failed, dropping", "event_id", event.ID, "err", err)
				_ = msg.Term()
				return
			}

			if err := msg.Ack(); err != nil {
				slog.Warn("ack failed", "event_id", event.ID, "err", err)
			}
		},
		jetstream.PullMaxMessages(1),
		jetstream.ConsumeErrHandler(func(_ jetstream.ConsumeContext, err error) {
			slog.Warn("consume error", "err", err)
		}),
	)
	if err != nil {
		return fmt.Errorf("start consuming todo events: %w", err)
	}
	s.cc = cc
	return nil
}

// PublishStatusMessage publishes a status message to NATS.
func (s *NatsService) PublishStatusMessage(ctx context.Context, message models.StatusMessage) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal status message: %w", err)
	}

	publishCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if _, err := s.js.Publish(publishCtx, StatusSubject, data); err != nil {
		return fmt.Errorf("publish status message: %w", err)
	}
	return nil
}

// Close stops pulling new messages, lets in-flight ones finish, then drains
// the connection.
func (s *NatsService) Close() error {
	if s.cc != nil {
		s.cc.Drain()
		select {
		case <-s.cc.Closed():
		case <-time.After(15 * time.Second):
		}
	}
	return s.nc.Drain()
}
