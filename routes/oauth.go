package routes

import (
	"github.com/21TechLabs/musiclms-backend/app"
	"github.com/gofiber/fiber/v2"
	"github.com/shareed2k/goth_fiber"
)

func SetupOAuth(f *fiber.App, app *app.Application) {
	f.Get("/user/oauth2/:provider/login", goth_fiber.BeginAuthHandler)
	f.Get("/user/oauth2/:provider/login/callback", app.OAuthController.GothicCallback)
}
