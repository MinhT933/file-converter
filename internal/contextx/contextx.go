package contextx

import (
	"context"

	"go.uber.org/zap"
)

type key int

const (
	loggerKey key = iota
)

func WithLogger(ctx context.Context, log *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, log)
}

func Logger(ctx context.Context) *zap.Logger {
	if v := ctx.Value(loggerKey); v != nil {
		return v.(*zap.Logger)
	}
	return zap.NewNop()
}
