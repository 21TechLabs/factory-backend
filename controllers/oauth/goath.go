package oauth_controller

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/21TechLabs/musiclms-backend/controllers"
	"github.com/21TechLabs/musiclms-backend/dto"
	"github.com/21TechLabs/musiclms-backend/models"
	"github.com/21TechLabs/musiclms-backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/discord"
	"github.com/shareed2k/goth_fiber"
	"gorm.io/gorm"
)

type OAuthController struct {
	Logger    *log.Logger
	UserStore *models.UserStore
}

func NewOAuthController(log *log.Logger, userStore *models.UserStore) *OAuthController {
	return &OAuthController{
		Logger:    log,
		UserStore: userStore,
	}
}

func init() {
	utils.LoadEnv()

	goth.UseProviders(
		discord.New(os.Getenv("DISCORD_CLIENT_ID"), os.Getenv("DISCORD_CLIENT_SECRET"), os.Getenv("DISCORD_REDIRECT_URI")),
	)
}

func (oac *OAuthController) GothicCallback(ctx *fiber.Ctx) error {
	gothicUser, err := goth_fiber.CompleteUserAuth(ctx)

	if err != nil {
		log.Printf("OAuth GothicCallback error go_fiber.CompleteUserAuth: %v\n", err)
		return utils.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	var provider = ctx.Params("provider")
	var userCreate dto.UserCreateDto

	switch provider {
	case "discord":
		discordUserWeb, err := discordUserGetDetail(gothicUser.AccessToken)
		var password = fmt.Sprintf("%s%s%s", discordUserWeb.Email, gothicUser.UserID, gothicUser.AccessToken)

		if err != nil {
			log.Printf("OAuth GothicCallback error discordUserGetDetail: %v\n", err)
			return utils.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}

		userCreate = dto.UserCreateDto{
			Name:            discordUserWeb.Username,
			Email:           discordUserWeb.Email,
			Password:        password,
			ConfirmPassword: password,
		}

	default:
		var password = fmt.Sprintf("%s%s%s", gothicUser.Email, gothicUser.UserID, gothicUser.AccessToken)
		userCreate = dto.UserCreateDto{
			Name:            gothicUser.Name,
			Email:           gothicUser.Email,
			Password:        password,
			ConfirmPassword: password,
		}
	}

	if len(userCreate.Email) == 0 {
		log.Printf("OAuth GothicCallback error Email not found for provider %s\n", provider)
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, "Email not found")
	}

	user, err := oac.UserStore.UserGetByEmail(userCreate.Email)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// create a new user
			// user role

			var isSubdomain = len(strings.Split(ctx.Get("Origin"), ".")) == 3

			user, err = oac.UserStore.UserCreate(userCreate, isSubdomain)
			if err != nil {
				return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
			}
		} else {
			log.Printf("OAuth GothicCallback error UserGetByEmail: %v\n", err)
			return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
		}
	}

	return controllers.SetLoginTokenAndSendResponse(ctx, user, false, oac.UserStore)
}
