package routes

import (
	"net/http"

	"github.com/21TechLabs/factory-backend/app"
	"github.com/21TechLabs/factory-backend/middleware"
	"github.com/markbates/goth/gothic"
)

func SetupOAuth(router *http.ServeMux, app *app.Application) {
	router.Handle("GET /user/oauth2/:provider/login", app.Middleware.CreateStackWithHandler(
		[]middleware.MiddlewareStack{},
		gothic.BeginAuthHandler,
	))

	router.Handle("GET /user/oauth2/:provider/login/callback", app.Middleware.CreateStackWithHandler(
		[]middleware.MiddlewareStack{app.Middleware.UserAuthMiddleware},
		app.OAuthController.GothicCallback,
	))
}
