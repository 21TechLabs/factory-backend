package routes

import (
	"github.com/21TechLabs/musiclms-backend/app"
	fiber "github.com/gofiber/fiber/v2"
)

func SetupRoutes(f *fiber.App, app *app.Application) {
	SetupUser(f, app)
	SetupFile(f, app)
	SetupOAuth(f, app)
	SetupPayments(f, app)
	SetupProduct(f, app)

	f.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the Factory Backend API!")
	})

	f.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "UP",
			"message": "API is running smoothly",
		})
	})
}
