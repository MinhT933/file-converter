package api

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/MinhT933/file-converter/internal/converter"          // <-- đường dẫn thật
	_ "github.com/MinhT933/file-converter/internal/converter/htmlwk" // blank-import để init() plugin
	"github.com/MinhT933/file-converter/internal/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

	outputName := strings.TrimSuffix(file.Filename, ".html") + ".pdf"

	userID := c.Locals("user_id")
	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		userIDStr = "9dab2d9f-47e9-4e74-b468-5042de616651"
	}
	outputPath := "./storage/" + userIDStr + "/" + outputName

	conversion := &domain.Conversion{
		ConversionID:      uuid.NewString(),
		OriginalFilename:  file.Filename,
		ConvertedFilename: outputName,
		Status:            "pending",
		UserID:            userIDStr,
		ExpiresAt:         time.Now().Add(7 * 24 * time.Hour),
	}

	_, conversionID, err := h.FileService.SaveConvertedFile(
		c.Context(),
		userIDStr,
		conversion,
	)
	log.Println("Conversion ID:", conversionID)

	if err != nil {
		log.Println("Error creating conversion in DB:", err)
		return fiber.ErrInternalServerError
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		log.Println("Error creating output file:", err)
		_ = h.FileService.UpdateConversionStatus(c.Context(), conversionID, "failed")
		return fiber.ErrInternalServerError
	}
	defer outFile.Close()

	// Stream → Pipe để không ngốn RAM
	err = converter.Registry["html_pdf"].Convert(c.Context(), src, outFile)
	if err != nil {
		log.Println("Error converting HTML to PDF:", err)
		_ = h.FileService.UpdateConversionStatus(c.Context(), conversionID, "failed")
		// Nếu lỗi, có thể là do file không hợp lệ hoặc plugin không hoạt động
		return fiber.ErrInternalServerError
	}
	if err := h.FileService.UpdateConversionStatus(c.Context(), conversionID, "success"); err != nil {
	}
	// Cập nhật trạng thái thành công
	_ = h.FileService.UpdateConversionStatus(c.Context(), conversionID, "success")

	userEmail := c.FormValue("email")
	if userEmail == "" {
		userEmail = "phammanhtoanhht@gmail.com"
	}

	// --- DRY: Gửi email báo thành công bằng Asynq ---
	if userEmail != "" {
		// Dùng link file (nếu bạn lưu file lại) hoặc ghi chú "file đã convert thành công"
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Println("Panic in async email goroutine:", r)
				}
			}()
			// Có thể là link download thực tế nếu bạn cho phép tải lại file sau này
			fileURL := outputPath // (hoặc link public nếu bạn upload lên cloud)
			payload, _ := json.Marshal(map[string]string{
				"email":    userEmail,
				"file_url": fileURL,
			})
			_, err := h.AsynqClient.Enqueue(
				asynq.NewTask("email:notify", payload),
				asynq.Queue("default"),
			)
			if err != nil {
				// Không ảnh hưởng user, chỉ log lỗi để dev biết
				log.Println("Gửi email thất bại:", err.Error())
			}
		}()
	}
	// Trả file PDF cho client và return luôn
	c.Type("application/pdf")
	c.Attachment(outputName)
	err = c.SendFile(outputPath)
	if err != nil {
		log.Println("SendFile error:", err)
		return fiber.ErrInternalServerError
	}
	// -------------------------------------------------

	return nil
}
