package routes

import (
	"github.com/21TechLabs/factory-be/controllers/oauth"
	"github.com/gofiber/fiber/v2"
	"github.com/shareed2k/goth_fiber"
)

func SetupOAuth(app *fiber.App) {
	app.Get("/user/oauth2/:provider/login", goth_fiber.BeginAuthHandler)
	app.Get("/user/oauth2/:provider/login/callback", oauth.GothicCallback)
}
