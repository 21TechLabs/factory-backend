package routes

import (
	"net/http"

	"github.com/21TechLabs/factory-backend/app"
	"github.com/21TechLabs/factory-backend/middleware"
)

func SetupProduct(router *http.ServeMux, app *app.Application) {

	router.Handle("GET /product/:appCode", app.Middleware.CreateStackWithHandler(
		[]middleware.MiddlewareStack{},
		app.ProductPlanController.GetProductByAppCode,
	))

	router.Handle("GET /product/:appCode/@me", app.Middleware.CreateStackWithHandler(
		[]middleware.MiddlewareStack{app.Middleware.UserAuthMiddleware},
		app.ProductPlanController.GetUsersActiveProductSubsctiptionByAppCode,
	))

	router.Handle("GET /product/:appCode/plans", app.Middleware.CreateStackWithHandler(
		[]middleware.MiddlewareStack{},
		app.PaymentsController.GetProductPlansByAppCode,
	))
}
