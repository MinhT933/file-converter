package api

import (
	"github.com/MinhT933/file-converter/internal/config"
	"github.com/MinhT933/file-converter/internal/infra/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
)

func RegisterRoutes(app *fiber.App, cfg *config.Config, client *asynq.Client, authProvider auth.Provider) {
	h := &Handler{Cfg: cfg, AsynqClient: client, AuthProvider: authProvider}

	v1 := app.Group("/api")
	v1.Post("/upload", h.Upload)
	// v1.Get("/status/:job_id", h.Status)

	v1.Post("/convert/html_pdf", h.ConvertHTMLPDF)
	v1.Post("/auth/social/login", h.SocialLogin)

}

// Upload       godoc
// @Summary     Upload file and enqueue converting job
// @Accept      multipart/form-data
// @Param       file formData file true "File to upload"
// @Success     200 {object} map[string]string
// @Router      /upload [post]
// func (h *Handler) Upload(c *fiber.Ctx) error { return c.SendString("todo") }

// Status godoc
// @Summary  Get job status
// @Param    job_id path string true "Job ID"
// @Success  200 {object} map[string]string
// @Router   /status/{job_id} [get]
// func (h *Handler) Status(c *fiber.Ctx) error { return c.SendString("todo") }
