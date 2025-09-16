package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func ProtectRoute (c *fiber.Ctx) error {
	c.Next()
	return nil
}
