package routes

import (
	"net/http"

	"github.com/21TechLabs/factory-backend/app"
	"github.com/21TechLabs/factory-backend/middleware"
)

// app *app.Application)
func SetupFile(router *http.ServeMux, app *app.Application) {
	router.Handle("POST /file", app.Middleware.CreateStackWithHandler(
		[]middleware.MiddlewareStack{app.Middleware.UserAuthMiddleware},
		app.FileController.FileUpload,
	))

	router.Handle("GET /file/stream/:fileKey", app.Middleware.CreateStackWithHandler(
		[]middleware.MiddlewareStack{app.Middleware.UserAuthMiddleware},
		app.FileController.FileStreamS3,
	))

}
