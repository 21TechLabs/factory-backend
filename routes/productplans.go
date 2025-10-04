package routes

import (
	"net/http"

	"github.com/21TechLabs/factory-backend/app"
	"github.com/21TechLabs/factory-backend/dto"
	"github.com/21TechLabs/factory-backend/middleware"
)

// SetupProductPlans registers product plan routes on the provided router.
// It registers POST /products/create, applies schema validation for PaymentPlanCreate,
// enforces user authentication, and dispatches requests to app.PaymentPlanController.CreatePaymentPlan.
func SetupProductPlans(router *http.ServeMux, app *app.Application) {

	router.Handle("POST /products/create", app.Middleware.CreateStackWithHandler(
		[]middleware.MiddlewareStack{app.Middleware.SchemaValidatorMiddleware(dto.DtoMapKeyPaymentPlanCreate), app.Middleware.UserAuthMiddleware},
		app.PaymentPlanController.CreatePaymentPlan,
	))

}