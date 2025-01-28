package dto

type DiscordTokenExchangeResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type DiscordGetExchangeTokenBody struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	Code         string `json:"code"`
	RedirectUri  string `json:"redirect_uri"`
}

type DiscordUserLoginBody struct {
	Code string `json:"code" required`
}

type DiscordUserWeb struct {
	ID                   string `json:"id"`
	Username             string `json:"username"`
	Avatar               any    `json:"avatar"`
	Discriminator        string `json:"discriminator"`
	PublicFlags          int    `json:"public_flags"`
	PremiumType          int    `json:"premium_type"`
	Flags                int    `json:"flags"`
	Banner               any    `json:"banner"`
	AccentColor          int    `json:"accent_color"`
	GlobalName           string `json:"global_name"`
	AvatarDecorationData any    `json:"avatar_decoration_data"`
	BannerColor          string `json:"banner_color"`
	MfaEnabled           bool   `json:"mfa_enabled"`
	Locale               string `json:"locale"`
	Email                string `json:"email"`
}
