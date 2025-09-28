package routes

import (
	"net/http"

	"github.com/21TechLabs/factory-backend/app"
	"github.com/21TechLabs/factory-backend/middleware"
)

func SetupRoutes(app *app.Application) *http.ServeMux {
	router := http.NewServeMux()

	router.Handle("GET /", app.Middleware.CreateStackWithHandler(
		[]middleware.MiddlewareStack{},
		app.HealthCheckController.HealthCheck,
	))

	router.Handle("GET /health", app.Middleware.CreateStackWithHandler(
		[]middleware.MiddlewareStack{},
		app.HealthCheckController.HealthCheck,
	))

	SetupUser(router, app)
	SetupFile(router, app)
	SetupOAuth(router, app)
	SetupPayments(router, app)
	SetupProduct(router, app)

	return router
}
