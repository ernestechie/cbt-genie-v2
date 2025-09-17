package routes

import (
	"github.com/ernestechie/cbt-genie-v2/controllers"
	"github.com/ernestechie/cbt-genie-v2/middleware"
	"github.com/gofiber/fiber/v2"
)

var RegisterUserRoutes  = func (c *fiber.App) {
	c.Get("/users", middleware.ProtectRoute(), controllers.GetAllUsers)
	c.Get("/users/:userId",  middleware.ProtectRoute(), controllers.GetUser)
	c.Post("/users",  middleware.ProtectRoute(), controllers.RegisterUser)
}
