package api

import (
    "encoding/json"
    "path/filepath"

    "github.com/gofiber/fiber/v2"
    "github.com/hibiken/asynq"
    "github.com/MinhT933/file-converter/internal/tasks"
	"github.com/MinhT933/file-converter/internal/config"
)

type Handler struct {
    Cfg         *config.Config
    AsynqClient *asynq.Client
}

// Upload       godoc
// @Summary     Upload file and enqueue converting job
// @Accept      multipart/form-data
// @Param       file formData file true "File to upload"
// @Success     200 {object} map[string]string
// @Router      /upload [post]
func (h *Handler) Upload(c *fiber.Ctx) error {
    // 1) Nhận file
    file, err := c.FormFile("file")
    if err != nil {
        return fiber.ErrBadRequest
    }

    // 2) Save lên disk
    inputPath := filepath.Join("/tmp", file.Filename)
    if err := c.SaveFile(file, inputPath); err != nil {
        return fiber.ErrInternalServerError
    }

    // 3) Tạo outputPath và đọc email
    outputPath := filepath.Join("/converted", file.Filename)
    userEmail := c.FormValue("email")

    if userEmail == "" {
        userEmail = "phammanhtoanhht@gmail.com"
    }


    // 4) Đóng gói payload
    payload, err := json.Marshal(map[string]string{
        "email":       userEmail,
        "input_path":  inputPath,
        "output_path": outputPath,
    })
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "cannot encode payload",
        })
    }

    // 5) Enqueue task chuyển file (phải nằm trong hàm này)
    info, err := h.AsynqClient.Enqueue(
        asynq.NewTask(tasks.TypeConvertFile, payload),
        asynq.Queue("default"),
    )
    if err != nil {
      return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }
    // 6) Trả về client bao gồm job ID nếu cần
    return c.JSON(fiber.Map{
        "status": "queued",
        "job_id": info.ID,
    })
}
