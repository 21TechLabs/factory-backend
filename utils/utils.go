package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func LoadEnv() error {
	return godotenv.Load()
}

func GetEnv(envName string, allowEmpty bool) string {
	var env string = os.Getenv(envName)
	if env == "" && !allowEmpty {
		panic(envName + " name does not exists.")
	}
	return env
}

type StringFromWith map[string]int

type StringsToReplace struct {
	Str     string
	Replace StringFromWith
}

func ReplaceStringWith(val StringsToReplace) string {
	var newStr string = val.Str
	for key, value := range val.Replace {
		newStr = strings.ReplaceAll(val.Str, key, fmt.Sprintf("%v", value))
	}
	return newStr
}

func GetRegenTime(traitHealth int, health int) float64 {
	var regenPercentPerMinute = 0.2
	var regenRate = float64(traitHealth) / 100.0 * regenPercentPerMinute
	return float64(traitHealth-health) / regenRate
}

// GetToken extracts an authentication token from the HTTP request.
// It first attempts to read a cookie named "token"; if retrieving the cookie
// returns an error other than http.ErrNoCookie that error is returned.
// If the cookie is absent or empty, it falls back to the Authorization header
// and, if present, removes a leading "Bearer " prefix. If no token is found,
// an error is returned.
func GetToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie("token")

	if err != nil && err != http.ErrNoCookie {
		return "", fmt.Errorf("token not found in cookies: %w", err)
	}
	var authToken string

	if cookie != nil {
		authToken = cookie.Value
	}

	if authToken == "" {
		authToken = r.Header.Get("Authorization")
		if authToken != "" {
			if len(authToken) > 10 {
				authToken = authToken[7:]
			}
		}
	}

	if authToken == "" {
		return "", RaiseError{Message: "Token not found in headers or body"}
	}
	return authToken, nil
}

func ValidateHeaderHMACSha256(body []byte, secret string, signature string) bool {
	hmac := hmac.New(sha256.New, []byte(secret))
	hmac.Write(body)
	dataHmac := hmac.Sum(nil)
	return signature == hex.EncodeToString(dataHmac)
}

func IsValidOrigin(origin, whitelistOrigins string) bool {
	if origin == "" {
		return false
	}

	if strings.TrimSpace(whitelistOrigins) == "*" {
		return true
	}

	whitelist := strings.Split(whitelistOrigins, ",")
	for _, whitelistedOrigin := range whitelist {
		stringMatchRegex, err := regexp.MatchString(whitelistedOrigin, origin)
		if err != nil {
			log.Default().Printf("Error matching origin %s with regex %s: %v", origin, whitelistedOrigin, err)
			continue
		}
		if strings.TrimSpace(whitelistedOrigin) == "*" || stringMatchRegex {
			return true
		}
	}

	return false
}

func GenerateRandomUUID() string {
	return uuid.NewString()
}

func GetJwtSecret() []byte {
	secret := GetEnv("JWT_SECRET", false)
	return []byte(secret)
}

func GenerateJWT(claim jwt.Claims, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(secret)
}

func ParseJWT(tokenString string, secret []byte) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
}

func StructToBytes(data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func ParseBytesToStruct(data []byte, out interface{}) error {
	err := json.Unmarshal(data, out)
	if err != nil {
		return err
	}
	return nil
}

func ParseQueryParams(r *http.Request, out interface{}) error {
	_queryParams := r.URL.Query()

	queryParams := make(map[string]interface{})

	for key, values := range _queryParams {
		if len(values) > 0 {
			if len(values) == 1 {
				queryParams[key] = values[0]
			} else {
				queryParams[key] = values
			}
		}
	}

	jsonData, err := json.Marshal(queryParams)
	if err != nil {
		return fmt.Errorf("failed to marshal query params: %w", err)
	}

	fmt.Println("Parsed query params:", string(jsonData))

	err = json.Unmarshal(jsonData, out)
	if err != nil {
		return fmt.Errorf("failed to unmarshal query params: %w", err)
	}
	return nil
}
