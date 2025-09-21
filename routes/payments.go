package routes

import (
	"github.com/21TechLabs/musiclms-backend/app"
	"github.com/21TechLabs/musiclms-backend/dto"
	"github.com/21TechLabs/musiclms-backend/models"
	"github.com/gofiber/fiber/v2"
)

func SetupPayments(f *fiber.App, app *app.Application) {
	f.Post("/payments/:paymentGateway", app.Middleware.SchemaValidatorMiddleware(func() interface{} {
		return &dto.CreateProductDto{}
	}), app.Middleware.UserAuthMiddleware, app.Middleware.HasRoleMiddleware([]models.UserRole{models.UserRoleClient}), app.PaymentsController.CreatePayment)

	f.Post("/payments/:paymentGateway/verify", app.PaymentsController.UpdatePaymentStatusWebhook)
}
