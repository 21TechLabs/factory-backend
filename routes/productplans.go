package routes

import (
	"net/http"

	"github.com/21TechLabs/factory-backend/app"
	"github.com/21TechLabs/factory-backend/dto"
	"github.com/21TechLabs/factory-backend/middleware"
)

func SetupProductPlans(router *http.ServeMux, app *app.Application) {

	router.Handle("POST /products/create", app.Middleware.CreateStackWithHandler(
		[]middleware.MiddlewareStack{app.Middleware.SchemaValidatorMiddleware(dto.DtoMapKeyPaymentPlanCreate), app.Middleware.UserAuthMiddleware},
		app.PaymentPlanController.CreatePaymentPlan,
	))

}
