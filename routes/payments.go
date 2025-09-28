package routes

import (
	"net/http"

	"github.com/21TechLabs/factory-backend/app"
	"github.com/21TechLabs/factory-backend/dto"
	"github.com/21TechLabs/factory-backend/middleware"
	"github.com/21TechLabs/factory-backend/models"
)

func SetupPayments(router *http.ServeMux, app *app.Application) {

	router.Handle("POST /payments/:paymentGateway", app.Middleware.CreateStackWithHandler(
		[]middleware.MiddlewareStack{
			app.Middleware.SchemaValidatorMiddleware(func() interface{} {
				return &dto.CreateProductDto{}
			}),
			app.Middleware.UserAuthMiddleware,
			app.Middleware.HasRoleMiddleware([]models.UserRole{models.UserRoleClient}),
		},
		app.PaymentsController.CreatePayment,
	))

	router.Handle("POST /payments/:paymentGateway/verify", app.Middleware.CreateStackWithHandler(
		[]middleware.MiddlewareStack{},
		app.PaymentsController.UpdatePaymentStatusWebhook,
	))

}
