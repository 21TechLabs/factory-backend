package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/21TechLabs/factory-be/controllers"
	"github.com/21TechLabs/factory-be/dto"
	"github.com/21TechLabs/factory-be/models"
	"github.com/21TechLabs/factory-be/utils"
	"github.com/gofiber/fiber/v2"
)

type DiscordConfig struct {
	OAuthBaseURL string
	ClientID     string
	ClientSecret string
	ResponseType string
	RedirectURI  string
	Scope        string
	APIEndpoint  string
}

func LoadDiscordConfig() DiscordConfig {
	return DiscordConfig{
		OAuthBaseURL: utils.GetEnv("DISCORD_OAUTH_BASE_URL", false),
		ClientID:     utils.GetEnv("DISCORD_CLIENT_ID", false),
		ClientSecret: utils.GetEnv("DISCORD_CLIENT_SECRET", false),
		ResponseType: utils.GetEnv("DISCORD_RESPONSE_TYPE", false),
		RedirectURI:  utils.GetEnv("DISCORD_REDIRECT_URI", false),
		Scope:        utils.GetEnv("DISCORD_SCOPE", false),
		APIEndpoint:  utils.GetEnv("DISCORD_API_ENDPOINT", false),
	}
}

// var DISCORD_OAUTH_BASE_URL string
// var DISCORD_CLIENT_ID string
// var DISCORD_RESPONSE_TYPE string
// var DISCORD_REDIRECT_URI string
// var DISCORD_SCOPE string
// var DISCORD_CLIENT_SECRET string
// var DISCORD_API_ENDPOINT string

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

var discordConfig DiscordConfig = LoadDiscordConfig()

func init() {
	utils.LoadEnv()

	// DISCORD_OAUTH_BASE_URL = utils.GetEnv("DISCORD_OAUTH_BASE_URL", false)
	// DISCORD_CLIENT_ID = utils.GetEnv("DISCORD_CLIENT_ID", false)
	// DISCORD_RESPONSE_TYPE = utils.GetEnv("DISCORD_RESPONSE_TYPE", false)
	// DISCORD_REDIRECT_URI = utils.GetEnv("DISCORD_REDIRECT_URI", false)
	// DISCORD_SCOPE = utils.GetEnv("DISCORD_SCOPE", false)
	// DISCORD_CLIENT_SECRET = utils.GetEnv("DISCORD_CLIENT_SECRET", false)
	// DISCORD_API_ENDPOINT = utils.GetEnv("DISCORD_API_ENDPOINT", false)
}

func discordBuildRedirectURI() string {

	return fmt.Sprintf("%s/authorize?client_id=%s&response_type=%s&redirect_uri=%s&scope=%s", discordConfig.OAuthBaseURL, discordConfig.ClientID, discordConfig.ResponseType, discordConfig.RedirectURI, discordConfig.Scope)
}

func DiscordRedirectURI(ctx *fiber.Ctx) error {
	return ctx.Redirect(discordBuildRedirectURI())
}

func discordGetExchangeToken(code string) (dto.DiscordTokenExchangeResponse, error) {
	var body url.Values = url.Values{}

	body.Add("client_id", discordConfig.ClientID)
	body.Add("client_secret", discordConfig.ClientSecret)
	body.Add("grant_type", "authorization_code")
	body.Add("code", code)
	body.Add("redirect_uri", discordConfig.RedirectURI)

	req, err := http.NewRequest("POST", discordConfig.APIEndpoint+"/oauth2/token", strings.NewReader(body.Encode()))

	if err != nil {
		return dto.DiscordTokenExchangeResponse{}, fmt.Errorf("failed to exchange Discord token: %w", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := httpClient.Do(req)

	if err != nil {
		return dto.DiscordTokenExchangeResponse{}, fmt.Errorf("failed to exchange Discord token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return dto.DiscordTokenExchangeResponse{}, fmt.Errorf("failed to exchange token, status: %d", resp.StatusCode)
	}

	jsonByte, err := io.ReadAll(resp.Body)

	if err != nil {
		return dto.DiscordTokenExchangeResponse{}, fmt.Errorf("failed to exchange Discord token: %w", err)
	}

	var exchangeToken dto.DiscordTokenExchangeResponse = dto.DiscordTokenExchangeResponse{}
	err = json.Unmarshal(jsonByte, &exchangeToken)

	if err != nil {
		return dto.DiscordTokenExchangeResponse{}, fmt.Errorf("failed to exchange Discord token: %w", err)
	}

	return exchangeToken, nil
}

func discordUserGetDetail(accessToken string) (dto.DiscordUserWeb, error) {

	req, err := http.NewRequest("GET", discordConfig.APIEndpoint+"/users/@me", strings.NewReader(""))

	if err != nil {
		return dto.DiscordUserWeb{}, fmt.Errorf("failed to fetch Discord user details: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err := httpClient.Do(req)

	if err != nil {
		return dto.DiscordUserWeb{}, fmt.Errorf("failed to fetch Discord user details: %w", err)
	}
	defer resp.Body.Close()

	jsonByte, err := io.ReadAll(resp.Body)

	if err != nil {
		return dto.DiscordUserWeb{}, fmt.Errorf("failed to fetch Discord user details: %w", err)
	}

	var discordUserWeb dto.DiscordUserWeb = dto.DiscordUserWeb{}

	err = json.Unmarshal(jsonByte, &discordUserWeb)

	if err != nil {
		return dto.DiscordUserWeb{}, fmt.Errorf("failed to fetch Discord user details: %w", err)
	}

	return discordUserWeb, nil
}

func DiscordUserLogin(ctx *fiber.Ctx) error {
	var DiscordUserLoginBody dto.DiscordUserLoginBody

	if err := ctx.BodyParser(&DiscordUserLoginBody); err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	response, err := discordGetExchangeToken(DiscordUserLoginBody.Code)

	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	discordUserWeb, err := discordUserGetDetail(response.AccessToken)

	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	user, err := models.UserGetByEmail(discordUserWeb.Email)

	if err != nil {
		if err.Error() == "not found" {
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
				return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
			}
		} else {
			return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error())

		}
	}

	// set login token and send response
	return controllers.SetLoginTokenAndSendResponse(ctx, user, false)
}
