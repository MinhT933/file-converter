package main

import (
	"context"
	"log"

	"github.com/MinhT933/file-converter/cmd/worker/consumer"
	"github.com/MinhT933/file-converter/internal/config"
	"github.com/MinhT933/file-converter/internal/contextx"
	"github.com/MinhT933/file-converter/internal/db"
	"github.com/MinhT933/file-converter/internal/nats"
	"github.com/MinhT933/file-converter/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	cfg := config.LoadConfigWorker()
    _, err := db.ConnectDB()  // Không truyền cfg.DB, giống server
    if err != nil {
        log.Fatal("Failed to connect to DB", zap.Error(err))
    }
	logger.Init(cfg.App.Env)
	ctx := contextx.WithLogger(context.Background(), logger.L())
	log := contextx.Logger(ctx)

	stream := "import"

	nc := nats.Connect(cfg.Nats.URL)
	js := nats.NewJetStream(nc)


	workerType := cfg.App.WorkerType

	switch workerType {
	case "pdf":
		// nats.StartConsumer(ctx, js, stream, consumer.NewPdfHandler())
	case "excel":
		nats.StartConsumer(ctx, js, stream, consumer.NewUserExcelConsumer())
	case "all":
		// nats.StartConsumer(ctx, js, stream, consumer.NewPdfHandler())
		nats.StartConsumer(ctx, js, stream, consumer.NewUserExcelConsumer())
	default:
		log.Fatal("Unknown WORKER_TYPE", zap.String("WORKER_TYPE", workerType))
	}

	select {}
}
