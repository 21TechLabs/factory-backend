package routes

import (
	"github.com/21TechLabs/musiclms-backend/app"
	"github.com/gofiber/fiber/v2"
)

func SetupProduct(f *fiber.App, app *app.Application) {
	f.Get("/product/:appCode", app.ProductPlanController.GetProductByAppCode)
	f.Get("/product/:appCode/@me", app.Middleware.UserAuthMiddleware, app.ProductPlanController.GetUsersActiveProductSubsctiptionByAppCode)
	f.Get("/product/:appCode/plans", app.PaymentsController.GetProductPlansByAppCode)
}
