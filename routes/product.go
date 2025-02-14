package routes

import (
	"github.com/21TechLabs/factory-be/controllers/products"
	"github.com/21TechLabs/factory-be/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupProduct(app *fiber.App) {
	app.Get("/product/:appCode", products.GetProductByAppCode)
	app.Get("/product/:appCode/@me", middleware.UserAuthMiddleware, products.GetUsersActiveProductSubsctiptionByAppCode)
}
