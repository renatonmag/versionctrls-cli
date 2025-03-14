package main

import (
	"github.com/renatonmag/version-ctrls-cli/app/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	routes.SetupRoutes(app)
	app.Listen(":3333")
}
