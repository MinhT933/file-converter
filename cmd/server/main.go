package main

import (
	_ "github.com/MinhT933/file-converter/docs"
	"github.com/MinhT933/file-converter/internal/api"
	"github.com/MinhT933/file-converter/internal/config"
	fiberSwagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
)

// @title           File Converter API
// @version         1.0
// @description     Upload & convert files asynchronously via Asynq queue.
// @contact.name    Your Name
// @contact.email   you@example.com
// @host      localhost:8080
// @BasePath  /api
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

	 app.Get("/swagger/*", fiberSwagger.HandlerDefault)

	api.RegisterRoutes(app, cfg, asynqClient)

	if err := app.Listen(":" + cfg.PortHTTP); err != nil {
		panic(err)
	}
}
