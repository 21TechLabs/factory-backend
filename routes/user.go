package routes

import (
	"github.com/21TechLabs/factory-be/controllers"
	"github.com/21TechLabs/factory-be/dto"
	"github.com/21TechLabs/factory-be/middleware"
	"github.com/21TechLabs/factory-be/models"
	fiber "github.com/gofiber/fiber/v2"
)

func SetupUser(app *fiber.App) {
	app.Post("/user/create", middleware.SchemaValidatorMiddleware(func() interface{} {
		return &dto.UserCreateDto{}
	}), controllers.UserCreate(models.Roles.Admin))

	app.Post("/user/login", middleware.SchemaValidatorMiddleware(func() interface{} {
		return &dto.UserLoginDto{}
	}), controllers.UserLogin)

	app.Post("/user/login/verify", middleware.UserAuthMiddleware, controllers.UserLoginVerify)

	app.Get("/user/reset-password", controllers.UserRequestPasswordResetLink)
	app.Post("/user/reset-password", middleware.SchemaValidatorMiddleware(func() interface{} { return &dto.UserPasswordUpdateDto{} }), controllers.UserPasswordUpdate)

	app.Get("/user/verify-email", controllers.UserVerifyEmailToken)

	app.Delete("/user/:id", middleware.UserAuthMiddleware, controllers.UserMarkForDeletion)

	app.Patch("/user/:id/password", middleware.SchemaValidatorMiddleware(func() interface{} {
		return &dto.UserPasswordUpdateDto{}
	}), middleware.UserAuthMiddleware, controllers.UserPasswordUpdate)

	app.Get("/user/logout", middleware.UserAuthMiddleware, controllers.UserLogout)

}
