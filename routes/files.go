package routes

import (
	"github.com/21TechLabs/musiclms-backend/app"
	fiber "github.com/gofiber/fiber/v2"
)

func SetupFile(f *fiber.App, app *app.Application) {
	f.Post("/file", app.Middleware.UserAuthMiddleware, app.FileController.FileUpload)
	f.Get("/file/stream/:fileKey", app.FileController.FileStreamS3)

}
