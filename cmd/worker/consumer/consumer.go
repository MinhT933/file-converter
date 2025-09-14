package consumer

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
)

type JobHandler[T any ] interface {
	Config() jetstream.ConsumerConfig
	Handle(ctx context.Context, payload T) error
	Parse(data []byte) (T, error)
}