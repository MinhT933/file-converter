package nats

import (
	"context"

	"github.com/MinhT933/file-converter/cmd/worker/consumer"
	"github.com/MinhT933/file-converter/internal/contextx"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

func StartConsumer[T any](ctx context.Context, js jetstream.JetStream, stream string, h consumer.JobHandler[T]) error {
	log := contextx.Logger(ctx)
	
	cfg := h.Config()
	log.Info("Starting consumer",
		// zap.String("stream", stream),
		zap.String("stream", stream),
		zap.String("durable", cfg.Durable),
		zap.Any("filter_subject", cfg.FilterSubject),
	)

	cs, err := js.CreateOrUpdateConsumer(ctx, stream, h.Config())
	if err != nil {
		log.Error("Failed to create or update consumer",
			zap.String("stream", stream),
			zap.String("durable", h.Config().Durable),
			zap.Error(err),
		)
		return err
	}

	_, err = cs.Consume(func(msg jetstream.Msg) {
		payload, err := h.Parse(msg.Data())
		if err != nil {
			log.Error("Failed to parse message",
				zap.String("subject", msg.Subject()),
				zap.ByteString("data", msg.Data()),
				zap.Error(err),
			)
			return
		}

		if err := h.Handle(ctx, payload); err != nil {
			log.Error("Failed to handle message",
				zap.String("subject", msg.Subject()),
				zap.Any("payload", payload),
				zap.Error(err),
			)
			return
		}

		log.Info("Message processed",
			zap.String("subject", msg.Subject()),
			zap.Any("payload", payload),
		)

		msg.Ack()
	})
	return err
}