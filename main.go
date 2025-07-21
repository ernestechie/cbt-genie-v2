package main

import (
	"log"

	"github.com/ernestechie/cbt-genie-v2/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var PORT_NUMBER string = ":10073"

func main() {
	app := fiber.New()
	routes.RegisterAuthRoutes(app)
	routes.RegisterUserRoutes(app)

	// Middleware
	app.Use(logger.New())

	log.Printf("App running in port %+v", PORT_NUMBER)
	app.Listen(PORT_NUMBER)
}
