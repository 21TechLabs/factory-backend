package routes

import (
	"github.com/21TechLabs/factory-be/controllers"
	"github.com/21TechLabs/factory-be/middleware"
	fiber "github.com/gofiber/fiber/v2"
)

func SetupFile(app *fiber.App) {
	app.Post("/file", middleware.UserAuthMiddleware, controllers.FileUpload)

}
