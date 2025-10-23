package routes

import (
	"net/http"

	"github.com/21TechLabs/factory-backend/app"
	"github.com/21TechLabs/factory-backend/dto"
	"github.com/21TechLabs/factory-backend/middleware"
	"github.com/21TechLabs/factory-backend/models"
)

func SetupProductPlans(router *http.ServeMux, app *app.Application) {

	router.Handle("GET /products", app.Middleware.CreateStackWithHandler(
		[]middleware.MiddlewareStack{},
		app.PaymentPlanController.GetProductPlans,
	))

	router.Handle("POST /products/create", app.Middleware.CreateStackWithHandler(
		[]middleware.MiddlewareStack{
			app.Middleware.SchemaValidatorMiddleware(dto.DtoMapKeyPaymentPlanCreate),
			app.Middleware.UserAuthMiddleware,
			app.Middleware.HasRoleMiddleware([]models.UserRole{models.UserRoleAdmin}),
		},
		app.PaymentPlanController.CreateProductPlan,
	))

	router.Handle("GET /products/{id}", app.Middleware.CreateStackWithHandler(
		[]middleware.MiddlewareStack{},
		app.PaymentPlanController.GetPaymentPlanByID,
	))

	router.Handle("PUT /products/{id}", app.Middleware.CreateStackWithHandler(
		[]middleware.MiddlewareStack{
			app.Middleware.SchemaValidatorMiddleware(dto.DtoMapKeyPaymentPlanCreate),
			app.Middleware.UserAuthMiddleware,
			app.Middleware.HasRoleMiddleware([]models.UserRole{models.UserRoleAdmin, models.UserRoleClient}),
		},
		app.PaymentPlanController.UpdatePaymentPlan,
	))

	router.Handle("DELETE /products/{id}", app.Middleware.CreateStackWithHandler(
		[]middleware.MiddlewareStack{
			app.Middleware.UserAuthMiddleware,
			app.Middleware.HasRoleMiddleware([]models.UserRole{models.UserRoleAdmin}),
		},
		app.PaymentPlanController.DeletePaymentPlan,
	))

	router.Handle("GET /products/{productId}/buy/{paymentGateway}", app.Middleware.CreateStackWithHandler(
		[]middleware.MiddlewareStack{
			app.Middleware.UserAuthMiddleware,
			app.Middleware.HasRoleMiddleware([]models.UserRole{models.UserRoleAdmin, models.UserRoleClient}),
		},
		app.PaymentPlanController.ProductBuy,
	))

	router.Handle("POST /products/{paymentGateway}/{webhook}", app.Middleware.CreateStackWithHandler(
		[]middleware.MiddlewareStack{},
		app.PaymentPlanController.ProcessWebhook,
	))

}
