package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/21TechLabs/factory-be/controllers"
	"github.com/21TechLabs/factory-be/dto"
	"github.com/21TechLabs/factory-be/models"
	"github.com/21TechLabs/factory-be/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

var DISCORD_OAUTH_BASE_URL string
var DISCORD_CLIENT_ID string
var DISCORD_RESPONSE_TYPE string
var DISCORD_REDIRECT_URI string
var DISCORD_SCOPE string

func init() {
	utils.LoadEnv()

	DISCORD_OAUTH_BASE_URL = utils.GetEnv("DISCORD_OAUTH_BASE_URL", false)
	DISCORD_CLIENT_ID = utils.GetEnv("DISCORD_CLIENT_ID", false)
	DISCORD_RESPONSE_TYPE = utils.GetEnv("DISCORD_RESPONSE_TYPE", false)
	DISCORD_REDIRECT_URI = utils.GetEnv("DISCORD_REDIRECT_URI", false)
	DISCORD_SCOPE = utils.GetEnv("DISCORD_SCOPE", false)
}

func discordBuildRedirectURI() string {
	return fmt.Sprintf("%s/authorize?client_id=%s&response_type=%s&redirect_uri=%s&scope=%s", DISCORD_OAUTH_BASE_URL, DISCORD_CLIENT_ID, DISCORD_RESPONSE_TYPE, DISCORD_REDIRECT_URI, DISCORD_SCOPE)
}

func DiscordRedirectURI(ctx *fiber.Ctx) error {
	return ctx.Redirect(discordBuildRedirectURI())
}

func discordGetExchangeToken(code string) (dto.DiscordTokenExchangeResponse, error) {
	var body url.Values = url.Values{}
	body.Add("client_id", utils.GetEnv("DISCORD_CLIENT_ID", false))
	body.Add("client_secret", utils.GetEnv("DISCORD_CLIENT_SECRET", false))
	body.Add("grant_type", "authorization_code")
	body.Add("code", code)
	body.Add("redirect_uri", DISCORD_REDIRECT_URI)

	var DiscordApiEndpoint = utils.GetEnv("DISCORD_API_ENDPOINT", false) + "/oauth2/token"

	req, err := http.NewRequest("POST", DiscordApiEndpoint, strings.NewReader(body.Encode()))

	if err != nil {
		return dto.DiscordTokenExchangeResponse{}, utils.RaiseError{Message: err.Error()}
	}

	client := &http.Client{}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)

	if err != nil {
		return dto.DiscordTokenExchangeResponse{}, utils.RaiseError{Message: err.Error()}
	}

	jsonByte, err := io.ReadAll(resp.Body)

	if err != nil {
		return dto.DiscordTokenExchangeResponse{}, utils.RaiseError{Message: err.Error()}
	}
	defer resp.Body.Close()

	var exhcangeToken dto.DiscordTokenExchangeResponse = dto.DiscordTokenExchangeResponse{}
	err = json.Unmarshal(jsonByte, &exhcangeToken)

	if err != nil {
		return dto.DiscordTokenExchangeResponse{}, utils.RaiseError{Message: err.Error()}
	}

	return exhcangeToken, nil
}

func discordUserGetDetail(accessToken string) (dto.DiscordUserWeb, error) {
	var DISCORD_API_ENDPOINT string = utils.GetEnv("DISCORD_API_ENDPOINT", false) + "/users/@me"
	req, err := http.NewRequest("GET", DISCORD_API_ENDPOINT, strings.NewReader(""))

	if err != nil {
		return dto.DiscordUserWeb{}, utils.RaiseError{Message: err.Error()}
	}

	var client = &http.Client{}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err := client.Do(req)

	if err != nil {
		return dto.DiscordUserWeb{}, utils.RaiseError{Message: err.Error()}
	}

	jsonByte, err := io.ReadAll(resp.Body)

	if err != nil {
		return dto.DiscordUserWeb{}, utils.RaiseError{Message: err.Error()}
	}
	defer resp.Body.Close()

	var discordUserWeb dto.DiscordUserWeb = dto.DiscordUserWeb{}

	err = json.Unmarshal(jsonByte, &discordUserWeb)

	if err != nil {
		return dto.DiscordUserWeb{}, utils.RaiseError{Message: err.Error()}
	}

	return discordUserWeb, nil
}

func DiscordUserLogin(ctx *fiber.Ctx) error {
	var DiscordUserLoginBody dto.DiscordUserLoginBody
	err := json.Unmarshal(ctx.Body(), &DiscordUserLoginBody)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(bson.M{
			"message": err.Error(),
		})
	}

	response, err := discordGetExchangeToken(DiscordUserLoginBody.Code)

	if err != nil {

		return ctx.Status(fiber.StatusBadRequest).JSON(bson.M{
			"message": err.Error(),
		})
	}

	discordUserWeb, err := discordUserGetDetail(response.AccessToken)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(bson.M{
			"message": err.Error(),
		})
	}

	user, err := models.UserGetByEmail(discordUserWeb.Email)

	if err != nil {
		if err.Error() != "not found" {
			return ctx.Status(fiber.StatusBadRequest).JSON(bson.M{
				"message": err.Error(),
			})
		}
	}

	if len(user.Name) == 0 {
		// create a new user
		var password = fmt.Sprintf("%s%s%s", discordUserWeb.Email, discordUserWeb.Username, response.AccessToken)
		var userCreate dto.UserCreateDto = dto.UserCreateDto{
			Name:            discordUserWeb.GlobalName,
			Email:           discordUserWeb.Email,
			Password:        password,
			ConfirmPassword: password,
		}

		user, err = models.UserCreate(userCreate, models.Roles.Client)

		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(bson.M{
				"message": err.Error(),
			})
		}
	}

	// set login token and send response
	controllers.SetLoginTokenAndSendResponse(ctx, user, false)
	return nil
}
