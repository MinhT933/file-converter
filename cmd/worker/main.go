package main

import (
	"github.com/MinhT933/file-converter/internal/config"
	"github.com/MinhT933/file-converter/internal/tasks"
	wk "github.com/MinhT933/file-converter/internal/worker"
	"github.com/hibiken/asynq"
)

func main() {
	cfg := config.Load()

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.RedisAddr, Password: cfg.RedisPass},
		asynq.Config{Concurrency: 4},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeConvert, wk.ConvertHandler())

	if err := srv.Run(mux); err != nil {
		panic(err)
	}
}
