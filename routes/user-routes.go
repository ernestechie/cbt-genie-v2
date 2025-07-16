package routes

import (
	"github.com/ernestechie/cbt-genie-v2/controllers"
	"github.com/gofiber/fiber/v2"
)

var RegisterUserRoutes  = func (c *fiber.App) {
	c.Get("/users", controllers.GetAllUsers)
	c.Get("/users/:userId", controllers.GetUser)
	c.Post("/users", controllers.RegisterUser)
}
