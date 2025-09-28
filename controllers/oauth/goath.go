package oauth_controller

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/21TechLabs/factory-backend/controllers"
	"github.com/21TechLabs/factory-backend/dto"
	"github.com/21TechLabs/factory-backend/models"
	"github.com/21TechLabs/factory-backend/utils"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
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

func (oac *OAuthController) GothicCallback(w http.ResponseWriter, r *http.Request) {
	gothicUser, err := gothic.CompleteUserAuth(w, r)

	if err != nil {
		log.Printf("OAuth GothicCallback error go_fiber.CompleteUserAuth: %v\n", err)
		utils.ErrorResponse(oac.Logger, w, http.StatusInternalServerError, []byte(err.Error()))
		return
	}

	var provider = r.PathValue("provider")
	var userCreate dto.UserCreateDto

	switch provider {
	case "discord":
		discordUserWeb, err := discordUserGetDetail(gothicUser.AccessToken)
		var password = fmt.Sprintf("%s%s%s", discordUserWeb.Email, gothicUser.UserID, gothicUser.AccessToken)

		if err != nil {
			log.Printf("OAuth GothicCallback error discordUserGetDetail: %v\n", err)
			utils.ErrorResponse(oac.Logger, w, http.StatusInternalServerError, []byte(err.Error()))
			return
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
		utils.ErrorResponse(oac.Logger, w, http.StatusInternalServerError, []byte("Email not found"))
		return
	}

	user, err := oac.UserStore.UserGetByEmail(userCreate.Email)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			user, err = oac.UserStore.UserCreate(userCreate)
			if err != nil {

				log.Printf("OAuth GothicCallback error UserCreate: %v\n", err)
				utils.ErrorResponse(oac.Logger, w, http.StatusInternalServerError, []byte(err.Error()))
				return
			}
		} else {
			log.Printf("OAuth GothicCallback error UserGetByEmail: %v\n", err)
			utils.ErrorResponse(oac.Logger, w, http.StatusInternalServerError, []byte(err.Error()))
			return
		}
	}

	controllers.SetLoginTokenAndSendResponse(oac.Logger, r, w, user, false, oac.UserStore)
}
