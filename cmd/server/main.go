//go:generate swag init -g cmd/server/main.go --parseDependency --parseInternal
package main

import (
	_ "github.com/MinhT933/file-converter/docs"
	"github.com/MinhT933/file-converter/internal/api"
	"github.com/MinhT933/file-converter/internal/config"
	fiberSwagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/hibiken/asynq"
	"log"
)

// @title           File Converter API
// @version         1.0
// @description     Upload & convert files asynchronously via Asynq queue.
// @contact.name    minht
// @contact.email   phammanhtoanhht933@gmail.com
// @host      127.0.0.1:8080
// @schemes   https
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


	//thêm
	app.Use(cors.New(cors.Config{
	    AllowOrigins: "http://127.0.0.1:8080, http://localhost:8080",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowCredentials: true,
        AllowHeaders:     "Content-Type, Authorization",
		ExposeHeaders: "Content-Disposition",
	}))

	app.Use(cors.New()) 

	app.Get("/swagger/*", fiberSwagger.HandlerDefault)

	api.RegisterRoutes(app, cfg, asynqClient)

    log.Fatal(app.ListenTLS(
    ":8080",
    "127.0.0.1.pem",
    "127.0.0.1-key.pem",
   ))


	if err := app.Listen(":" + cfg.PortHTTP); err != nil {
		panic(err)
	}
}
