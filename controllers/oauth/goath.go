package oauth

import (
	"fmt"
	"os"

	"github.com/21TechLabs/factory-be/controllers"
	"github.com/21TechLabs/factory-be/dto"
	"github.com/21TechLabs/factory-be/models"
	"github.com/21TechLabs/factory-be/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"github.com/shareed2k/goth_fiber"
)

func init() {
	utils.LoadEnv()

	goth.UseProviders(
		google.New(os.Getenv("OAUTH_GOOGLE_KEY"), os.Getenv("OAUTH_GOOGLE_SECRET"), "http://localhost:3000/auth/google/callback"),
	)
}

func GothicCallback(ctx *fiber.Ctx) error {
	gothicUser, err := goth_fiber.CompleteUserAuth(ctx)

	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	var password = fmt.Sprintf("%s%s%s", gothicUser.Email, gothicUser.UserID, gothicUser.AccessToken)
	var userCreate dto.UserCreateDto = dto.UserCreateDto{
		Name:            gothicUser.Name,
		Email:           gothicUser.Email,
		Password:        password,
		ConfirmPassword: password,
	}

	// find user with the email

	user, err := models.UserGetByEmail(gothicUser.Email)

	if err != nil {
		if err.Error() != "not found" {
			return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
		}
	}

	if len(user.Name) == 0 {
		// create a new user
		user, err = models.UserCreate(userCreate, models.Roles.Client)

		if err != nil {
			return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
		}
	}

	return controllers.SetLoginTokenAndSendResponse(ctx, user, false)
}
