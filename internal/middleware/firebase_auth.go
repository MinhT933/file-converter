package middleware

import (
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/gofiber/fiber/v2"
)

func NewFirebaseClient(authClient *auth.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Lấy token từ header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return fiber.ErrUnauthorized
		}

		// Tách token ra
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader { // Nếu không có "Bearer "
			token = authHeader
		}

		// Verify token
		ctx := c.Context()
		user, err := authClient.VerifyIDToken(ctx, token)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		c.Locals("user", user) // Lưu thông tin user vào context
		return c.Next()
	}
}
