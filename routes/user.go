package routes

import (
	"github.com/21TechLabs/musiclms-backend/app"
	"github.com/21TechLabs/musiclms-backend/dto"
	fiber "github.com/gofiber/fiber/v2"
)

func SetupUser(f *fiber.App, app *app.Application) {
	f.Post("/user/create", app.Middleware.SchemaValidatorMiddleware(func() interface{} {
		return &dto.UserCreateDto{}
	}), app.UserController.UserCreate)

	f.Post("/user/login", app.Middleware.SchemaValidatorMiddleware(func() interface{} {
		return &dto.UserLoginDto{}
	}), app.UserController.UserLogin)

	f.Patch("/user/update", app.Middleware.SchemaValidatorMiddleware(func() interface{} {
		return &dto.UserUpdateDto{}
	}), app.Middleware.UserAuthMiddleware, app.UserController.UserUpdateDto)

	f.Post("/user/login/verify", app.Middleware.UserAuthMiddleware, app.UserController.UserLoginVerify)

	f.Get("/user/reset-password", app.UserController.UserRequestPasswordResetLink)
	f.Post("/user/reset-password", app.Middleware.SchemaValidatorMiddleware(func() interface{} { return &dto.UserPasswordUpdateDto{} }), app.UserController.UserPasswordUpdate)

	f.Get("/user/verify-email", app.UserController.UserVerifyEmailToken)

	f.Delete("/user/:id", app.Middleware.UserAuthMiddleware, app.UserController.UserMarkForDeletion)

	f.Patch("/user/:id/password", app.Middleware.SchemaValidatorMiddleware(func() interface{} {
		return &dto.UserPasswordUpdateDto{}
	}), app.Middleware.UserAuthMiddleware, app.UserController.UserPasswordUpdate)

	f.Get("/user/logout", app.Middleware.UserAuthMiddleware, app.UserController.UserLogout)

}
