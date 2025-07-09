package main

import (
	"github.com/MinhT933/file-converter/internal/api"
	"github.com/MinhT933/file-converter/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
)

func main() {
	cfg := config.Load()

	// Asynq client
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
	})
	defer asynqClient.Close()

	app := fiber.New(fiber.Config{
		BodyLimit: cfg.MaxUploadMB * 1024 * 1024,
	})

	api.RegisterRoutes(app, cfg, asynqClient)

	if err := app.Listen(":" + cfg.PortHTTP); err != nil {
		panic(err)
	}
}
