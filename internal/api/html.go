package api

import (
	"encoding/json"
	"io"
	"log"
	"strings"

	"github.com/MinhT933/file-converter/internal/converter"          // <-- đường dẫn thật
	_ "github.com/MinhT933/file-converter/internal/converter/htmlwk" // blank-import để init() plugin
	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
)

// ConvertHTMLPDF godoc
// @Summary  Convert HTML to PDF
// @Accept   multipart/form-data
// @Produce  application/pdf
// @Param    file formData file true "HTML file"
// @Success 200 {file} file
// @Router   /convert/html_pdf [post]
func (h *Handler) ConvertHTMLPDF(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return fiber.ErrBadRequest
	}
	src, _ := file.Open()
	defer src.Close()

	// Stream → Pipe để không ngốn RAM
	pr, pw := io.Pipe()
	errCh := make(chan error, 1)
	go func() {
		defer pw.Close()
		errCh <- converter.Registry["html_pdf"].Convert(c.Context(), src, pw)
		pw.Close()
	}()

	// Xây dựng tên file lưu (tuỳ dự án, bạn có thể lưu lên disk nếu muốn gửi link qua mail)
	// Ví dụ lưu ra file, hoặc upload lên cloud, hoặc tạo link tạm thời
	savedFilePath := "/converted/" + strings.TrimSuffix(file.Filename, ".html") + ".pdf"

	// Trả file PDF cho client
	c.Type("application/pdf")
	c.Attachment(strings.TrimSuffix(file.Filename, ".html") + ".pdf")
	err = c.SendStream(pr)

	// --- DRY: Gửi email báo thành công bằng Asynq ---

	userEmail := c.FormValue("email")
	log.Printf("EMAIL1: %s", userEmail)
	if userEmail == "" {
		userEmail = "phammanhtoanhht@gmail.com"
	}
	log.Printf("EMAIL: %s", userEmail)
	if userEmail != "" {
		// Dùng link file (nếu bạn lưu file lại) hoặc ghi chú "file đã convert thành công"
		go func() {
			// Có thể là link download thực tế nếu bạn cho phép tải lại file sau này
			fileURL := savedFilePath // (hoặc link public nếu bạn upload lên cloud)
			payload, _ := json.Marshal(map[string]string{
				"email":    userEmail,
				"file_url": fileURL,
			})
			_, err := h.AsynqClient.Enqueue(
				asynq.NewTask("email:notify", payload),
				asynq.Queue("default"),
			)
			log.Printf("payload", payload)
			if err != nil {
				// Không ảnh hưởng user, chỉ log lỗi để dev biết
				log.Println("Gửi email báo thành công thất bại:", err.Error())
			}
		}()
	}
	// -------------------------------------------------

	return err
}
