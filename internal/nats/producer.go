package nats

import (
	"context"
	"fmt"

	"github.com/MinhT933/file-converter/internal/contextx"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

func PublishSyncData(ctx context.Context, js jetstream.JetStream, subject string, payload []byte) error {
	log := contextx.Logger(ctx)
	ack, err := js.Publish(ctx, subject, payload)
	if err != nil {
		log.Error("failed to publish to subject",
			zap.String("subject", subject),
			zap.Error(err),
		)
		return fmt.Errorf("failed to publish to subject %s: %w", subject, err)
	}

	if ack == nil || ack.Duplicate {
		log.Error("publish ack failed or duplicate for subject",
			zap.String("subject", subject),
			zap.Error(err),
		)
		return fmt.Errorf("publish ack failed or duplicate for subject %s", subject)
	}

	log.Info("published sync", zap.String("subject", subject), zap.Uint64("seq", ack.Sequence))

	return nil
}

func PublishSyncMessage(ctx context.Context, js jetstream.JetStream, msg *nats.Msg, opts ...jetstream.PublishOpt) error {
	log := contextx.Logger(ctx)
	ack, err := js.PublishMsg(ctx, msg, opts...)
	if err != nil {
		log.Error("failed to publish to subject",
			zap.String("subject", msg.Subject),
			zap.Error(err),
		)
		return fmt.Errorf("failed to publish to subject %s: %w", msg.Subject, err)
	}

	if ack == nil || ack.Duplicate {
		log.Error("publish ack failed or duplicate for subject",
			zap.String("subject", msg.Subject),
			zap.Error(err),
		)
		return fmt.Errorf("publish ack failed or duplicate for subject %s", msg.Subject)
	}

	return nil
}

func PublishAsyncData(js jetstream.JetStream, subject string, payload []byte, opts ...jetstream.PublishOpt) (jetstream.PubAckFuture, error) {
	log := contextx.Logger(context.Background())

	future, err := js.PublishAsync(subject, payload, opts...)
	if err != nil {
		log.Error("failed to publish async",
			zap.String("subject", subject),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to publish async to subject %s: %w", subject, err)
	}

	log.Info("async publish sent",
		zap.String("subject", subject),
		zap.Int("payload_size", len(payload)),
	)

	return future, nil
}

func PublishAsyncMessage(js jetstream.JetStream, msg *nats.Msg, opts ...jetstream.PublishOpt) (jetstream.PubAckFuture, error) {
	log := contextx.Logger(context.Background())

	future, err := js.PublishMsgAsync(msg, opts...)
	if err != nil {
		log.Error("failed to publish async msg",
			zap.String("subject", msg.Subject),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to publish async msg to subject %s: %w", msg.Subject, err)
	}

	log.Info("async publish msg sent",
		zap.String("subject", msg.Subject),
		zap.Int("payload_size", len(msg.Data)),
	)

	return future, nil
}

func LogAsyncAck(future jetstream.PubAckFuture, logger *zap.Logger) {
	go func() {
		select {
		case ack := <-future.Ok():
			if ack != nil {
				logger.Info("async ack received",
					zap.String("stream", ack.Stream),
					zap.Uint64("seq", ack.Sequence))
			}
		case err := <-future.Err():
			if err != nil {
				logger.Error("async publish failed", zap.Error(err))
			}
		}
	}()
}

