package oauth_controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/21TechLabs/factory-backend/dto"
	"github.com/21TechLabs/factory-backend/utils"
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

var discordConfig DiscordConfig = LoadDiscordConfig()

func discordUserGetDetail(accessToken string) (dto.DiscordUserWeb, error) {
	var httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}

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
