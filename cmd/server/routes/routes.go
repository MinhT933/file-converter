package routes

import (
	handlers "github.com/MinhT933/file-converter/cmd/server/handler"
	"github.com/MinhT933/file-converter/internal/db"
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

func RegisterRoutes(app *fiber.App, js jetstream.JetStream, log *zap.Logger) {
	grHealth := app.Group("health")
	grImport := app.Group("import")

	// check health
	grHealth.Get("/ping", handlers.PingHandler)

	// Connect to database
	database, err := db.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}

	//import
	importHandler := handlers.NewImportHandler(database, js, log)
	grImport.Post("/upload", importHandler.Upload)
}
