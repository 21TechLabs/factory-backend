package routes

import (
	"net/http"

	"github.com/21TechLabs/factory-backend/app"
	"github.com/21TechLabs/factory-backend/middleware"
)

// SetupRoutes configures and returns an http.ServeMux populated with the application's HTTP routes.
// 
// It registers the root (GET "/") and health (GET "/health") endpoints to the application's health check
// handler and sets up user, file, OAuth, and product-plan related routes by invoking the respective setup
// functions.
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
	SetupProductPlans(router, app)

	return router
}