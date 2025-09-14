package logger

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log  *zap.Logger
	once sync.Once
)

func Init(env string) {
	once.Do(func() {
		var cfg zap.Config
		if env == "prod" {
			cfg = zap.NewProductionConfig()
			cfg.EncoderConfig.TimeKey = "timestamp"
			cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		} else {
			cfg = zap.NewDevelopmentConfig()
			cfg.EncoderConfig.TimeKey = "ts"
			cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		}

		l, err := cfg.Build()
		if err != nil {
			panic(err)
		}
		log = l
	})
}

func L() *zap.Logger {
	if log == nil {
		panic("logger not initialized, call logger.Init() first")
	}
	return log
}
