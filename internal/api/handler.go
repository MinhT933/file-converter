package api

import (
	"github.com/MinhT933/file-converter/internal/config"
	"github.com/MinhT933/file-converter/internal/infra/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
)

// Handler chứa dependencies cho API handlers
type Handler struct {
	Cfg          *config.Config
	AsynqClient  *asynq.Client
	AuthProvider auth.Provider
}

// Upload handles file upload
func (h *Handler) Upload(c *fiber.Ctx) error {
	return c.SendString("Upload functionality - TODO")
}

// Status handles job status check
func (h *Handler) Status(c *fiber.Ctx) error {
	jobID := c.Params("job_id")
	return c.JSON(fiber.Map{
		"job_id": jobID,
		"status": "processing",
	})
}
