package api

import (
	"github.com/MinhT933/file-converter/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
)

type Handler struct {
	Cfg         *config.Config
	AsynqClient *asynq.Client
}

func RegisterRoutes(app *fiber.App, cfg *config.Config, client *asynq.Client) {
	h := &Handler{Cfg: cfg, AsynqClient: client}

	v1 := app.Group("/api")
	v1.Post("/upload", h.Upload)
	v1.Get("/status/:job_id", h.Status)
}

// Upload       godoc
// @Summary     Upload file and enqueue converting job
// @Accept      multipart/form-data
// @Param       file formData file true "File to upload"
// @Success     200 {object} map[string]string
// @Router      /upload [post]
func (h *Handler) Upload(c *fiber.Ctx) error { return c.SendString("todo") }

// Status godoc
// @Summary  Get job status
// @Param    job_id path string true "Job ID"
// @Success  200 {object} map[string]string
// @Router   /status/{job_id} [get]
func (h *Handler) Status(c *fiber.Ctx) error { return c.SendString("todo") }
