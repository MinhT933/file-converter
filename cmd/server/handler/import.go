package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	contants "github.com/MinhT933/file-converter/internal/contants"
	natsCus "github.com/MinhT933/file-converter/internal/nats"
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

type ImportHandler struct {
	DB  *sql.DB
	JS  jetstream.JetStream
	Log *zap.Logger
}

func NewImportHandler(db *sql.DB, js jetstream.JetStream, log *zap.Logger) *ImportHandler {
	return &ImportHandler{
		DB:  db,
		JS:  js,
		Log: log,
	}
}

// UploadFile godoc
// @Summary Upload a CSV file
// @Tags Import
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV file"
// @Success 200 {object} map[string]string
// @Router /import/upload [post]
func (h *ImportHandler) Upload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		h.Log.Error("failed to read file", zap.Error(err))
		return fiber.NewError(fiber.StatusBadRequest, "invalid file")
	}

	// Ensure upload directory exists and create an absolute path for the saved file
	relPath := filepath.Join(contants.ImportFilePath, file.Filename)
	absPath, _ := filepath.Abs(relPath)
	savePath := absPath
	fmt.Printf("Save file to: %s\n", savePath)

	if err := os.MkdirAll(filepath.Dir(savePath), os.ModePerm); err != nil {
		h.Log.Error("failed to create directory", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "cannot create directory")
	}

	if err := c.SaveFile(file, savePath); err != nil {
		h.Log.Error("failed to save file", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "cannot save file")
	}

	// publish absolute path so worker can open it reliably
	payload := map[string]string{"path": savePath}
	res, err := json.Marshal(payload)
	if err != nil {
		h.Log.Error("failed to marshal payload", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "marshal error")
	}

	ack, err := natsCus.PublishAsyncData(h.JS, contants.SubjectExcel, res)
	if err != nil {
		h.Log.Error("failed to publish", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "publish error")
	}

	go func() {
		select {
		case err := <-ack.Err():
			h.Log.Error("nats ack error", zap.Error(err))
		case pa := <-ack.Ok():
			h.Log.Info("nats ack ok",
				zap.String("stream", pa.Stream),
				zap.String("subject", strconv.FormatUint(uint64(pa.Sequence), 10)),
			)
		}
	}()

	return c.JSON(fiber.Map{
		"message": "upload success",
		"path":    savePath,
	})
}
