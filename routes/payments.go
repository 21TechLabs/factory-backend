package routes

import (
	"github.com/21TechLabs/factory-be/controllers/payments"
	"github.com/21TechLabs/factory-be/dto"
	"github.com/21TechLabs/factory-be/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupPayments(app *fiber.App) {
	app.Post("/payments/:paymentType", middleware.SchemaValidatorMiddleware(func() interface{} {
		return &dto.CreateProductDto{}
	}), middleware.UserAuthMiddleware, payments.CreatePayment)
}
