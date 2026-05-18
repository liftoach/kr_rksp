package handler

import (
	auth "kr/finance-app/pkg/jwt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(jwtManager *auth.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return fiber.ErrUnauthorized
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return fiber.ErrUnauthorized
		}

		tokenStr := parts[1]

		claims, err := jwtManager.Parse(tokenStr)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		c.Locals("user_id", claims.UserID.String())

		return c.Next()
	}
}

func CORSMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Method() == fiber.MethodOptions {
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Next()
	}
}
