package handlers

import "github.com/gofiber/fiber/v2"

// Ping godoc
// @Summary Health Check
// @Description Simple ping endpoint
// @Tags Health
// @Success 200 {string} string "pong"
// @Router /health/ping [get]
func PingHandler(c *fiber.Ctx) error {
	return c.SendString("pong 1212")
}
