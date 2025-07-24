//go:generate swag init -g cmd/server/main.go --parseDependency --parseInternal
package main

import (
	"context"
	"log"

	_ "github.com/MinhT933/file-converter/docs"
	"github.com/MinhT933/file-converter/internal/api"
	"github.com/MinhT933/file-converter/internal/config"
	"github.com/MinhT933/file-converter/internal/infra/auth"
	"github.com/MinhT933/file-converter/internal/infra/firebase"
	fiberSwagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/hibiken/asynq"
	"github.com/MinhT933/file-converter/internal/repositories"
	"github.com/MinhT933/file-converter/internal/services"
)

// @title           File Converter API
// @version         1.0
// @description     Upload & convert files asynchronously via Asynq queue.
// @contact.name    minht
// @contact.email   phammanhtoanhht933@gmail.com
// @host      localhost:8080
// @schemes   https
// @BasePath  /api
func main() {
	cfg := config.Load()
	ctx := context.Background()

	fb := firebase.NewClient(ctx, cfg.FirebaseCredFile)





	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("💥 Failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)

	authProvider := auth.NewFirebaseProvider(fb.Auth)

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
		AllowOrigins:     "http://127.0.0.1:8080, http://localhost:8080, https://localhost:3000/, http://localhost:3000/",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowCredentials: true,
		AllowHeaders:     "Content-Type, Authorization",
		ExposeHeaders:    "Content-Disposition",
	}))

	app.Get("/swagger/*", fiberSwagger.HandlerDefault)

	api.RegisterRoutes(app, cfg, asynqClient, authProvider, authService)

	log.Fatal(app.ListenTLS(
		":8080",
		"127.0.0.1.pem",
		"127.0.0.1-key.pem",
	))

	if err := app.Listen(":" + cfg.PortHTTP); err != nil {
		panic(err)
	}
}
