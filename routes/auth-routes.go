package routes

import (
	"github.com/ernestechie/cbt-genie-v2/controllers"
	"github.com/gofiber/fiber/v2"
)

var RegisterAuthRoutes  = func (c *fiber.App) {
	c.Post("/auth/get-started", controllers.GetStarted)
	c.Post("/auth/verify-otp", controllers.VerifyOtp)
}
