package routes

import (
	"net/http"

	"github.com/21TechLabs/factory-backend/app"
	"github.com/21TechLabs/factory-backend/dto"
	"github.com/21TechLabs/factory-backend/middleware"
)

func SetupUser(router *http.ServeMux, app *app.Application) {

	router.Handle("POST /user/create", app.Middleware.CreateStackWithHandler(
		[]middleware.MiddlewareStack{
			app.Middleware.SchemaValidatorMiddleware(func() interface{} {
				return &dto.UserCreateDto{}
			}),
		},
		app.UserController.UserCreate,
	))

	router.Handle("POST /user/login", app.Middleware.CreateStackWithHandler(
		[]middleware.MiddlewareStack{
			app.Middleware.SchemaValidatorMiddleware(func() interface{} {
				return &dto.UserLoginDto{}
			}),
		},
		app.UserController.UserLogin,
	))

	router.Handle("PATCH /user/update", app.Middleware.CreateStackWithHandler(
		[]middleware.MiddlewareStack{
			app.Middleware.SchemaValidatorMiddleware(func() interface{} {
				return &dto.UserUpdateDto{}
			}),
			app.Middleware.UserAuthMiddleware,
		},
		app.UserController.UserUpdateDto,
	))

	router.Handle("POST /user/login/verify", app.Middleware.CreateStackWithHandler(
		[]middleware.MiddlewareStack{app.Middleware.UserAuthMiddleware},
		app.UserController.UserLoginVerify,
	))

	router.Handle("GET /user/reset-password", app.Middleware.CreateStackWithHandler(
		[]middleware.MiddlewareStack{},
		app.UserController.UserRequestPasswordResetLink,
	))

	router.Handle("POST /user/reset-password", app.Middleware.CreateStackWithHandler(
		[]middleware.MiddlewareStack{
			app.Middleware.SchemaValidatorMiddleware(func() interface{} {
				return &dto.UserPasswordUpdateDto{}
			}),
		},
		app.UserController.UserPasswordUpdate,
	))

	router.Handle("GET /user/verify-email", app.Middleware.CreateStackWithHandler(
		[]middleware.MiddlewareStack{},
		app.UserController.UserVerifyEmailToken,
	))

	router.Handle("DELETE /user/:id", app.Middleware.CreateStackWithHandler(
		[]middleware.MiddlewareStack{app.Middleware.UserAuthMiddleware},
		app.UserController.UserMarkForDeletion,
	))

	router.Handle("PATCH /user/:id/password", app.Middleware.CreateStackWithHandler(
		[]middleware.MiddlewareStack{
			app.Middleware.SchemaValidatorMiddleware(func() interface{} {
				return &dto.UserPasswordUpdateDto{}
			}),
			app.Middleware.UserAuthMiddleware,
		},
		app.UserController.UserPasswordUpdate,
	))

	router.Handle("GET /user/logout", app.Middleware.CreateStackWithHandler(
		[]middleware.MiddlewareStack{app.Middleware.UserAuthMiddleware},
		app.UserController.UserLogout,
	))

}
