package main

import (
	"context"
	"fmt"
	"os"

	"github.com/MinhT933/file-converter/cmd/server/routes"
	docs "github.com/MinhT933/file-converter/docs"
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

	// Configure swagger doc host/schemes at runtime so Swagger UI works in dev/prod
	// Choose scheme via APP_SCHEME env var (http or https). Default mapping:
	// - dev/local -> http
	// - prod -> https
	appScheme := os.Getenv("APP_SCHEME")
	if appScheme == "" {
		if cfg.Env == "dev" || cfg.Env == "local" {
			appScheme = "http"
		} else {
			appScheme = "https"
		}
	}

	if cfg.Env == "dev" || cfg.Env == "local" {
		docs.SwaggerInfo.Host = "localhost:" + cfg.PortHTTP
		docs.SwaggerInfo.Schemes = []string{"http"}
	} else {
		// allow override with SWAGGER_HOST env var if needed
		swaggerHost := os.Getenv("SWAGGER_HOST")
		if swaggerHost == "" {
			swaggerHost = "api-convert-file.minht.io.vn"
		}
		docs.SwaggerInfo.Host = swaggerHost
		docs.SwaggerInfo.Schemes = []string{appScheme}
	}

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
		AllowOrigins:     cfg.CORSOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowCredentials: true,
		AllowHeaders:     "Content-Type, Authorization",
		ExposeHeaders:    "Content-Disposition",
	}))


	app.Get("/swagger/*", fiberSwagger.HandlerDefault)

	routes.RegisterRoutes(app, js, log)
		// Log resolved swagger url (respecting scheme)
		swaggerURL := fmt.Sprintf("%s://%s/swagger/index.html", appScheme, docs.SwaggerInfo.Host)
		log.Info("Application started",
			zap.String("swagger_url", swaggerURL),
		)

	if err := app.Listen(":" + cfg.PortHTTP); err != nil {
		panic(err)
	}
}
