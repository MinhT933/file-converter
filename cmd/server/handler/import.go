package handlers

import (
	"github.com/gin-gonic/gin"
)


func ImportHandler(c *fiber.Context) {
	file, err := c.FromFile("file")
	if err != nil {
		c.Status(400).JSON(gin.H{"error": "Failed to get file"})
		return
	}

	savePath := fmt.Sprintf(".%s/%s", contains.ImportFilePath, file.Filename)
	if err := errs.SaveFile(file, savePath); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "File uploaded successfully",
		"path":    savePath,
	})
}