package main

import (
	"context"
	"fmt"

	"github.com/MinhT933/file-converter/cmd/server/routes"
	_ "github.com/MinhT933/file-converter/docs"
	"github.com/MinhT933/file-converter/internal/config"
	"github.com/MinhT933/file-converter/internal/contextx"
	"github.com/MinhT933/file-converter/internal/nats"
	"github.com/MinhT933/file-converter/pkg/logger"
	fiberSwagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

// @title           File Converter API
// @version         2.0
// @description     Upload & convert files asynchronously via Asynq queue.
// @contact.name    minht12
// @contact.email   phammanhtoanhht933@gmail.com
// @host      localhost:8081
// @schemes   http
// @BasePath  /api
func main() {
	app := fiber.New()
	cfg := config.Load()

	logger.Init(cfg.Env)

	ctx := contextx.WithLogger(context.Background(), logger.L())
	log :=  contextx.Logger(ctx)
	nc := nats.Connect(cfg.Nats.URL)
	js := nats.NewJetStream(nc)

	// Asynq client
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
	})

	defer asynqClient.Close()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://127.0.0.1:8081, http://localhost:8081, https://localhost:3000/, http://localhost:3000/",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowCredentials: true,
		AllowHeaders:     "Content-Type, Authorization",
		ExposeHeaders:    "Content-Disposition",
	}))


	app.Get("/swagger/*", fiberSwagger.HandlerDefault)

	routes.RegisterRoutes(app, js, log)
		log.Info("Application started",
		zap.String("swagger_url", fmt.Sprintf("http://localhost:%s/swagger/index.html", cfg.PortHTTP)),
	)

	if err := app.Listen(":" + cfg.PortHTTP); err != nil {
		panic(err)
	}
}
