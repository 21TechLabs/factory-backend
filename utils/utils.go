package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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

func MapToStruct[T any](d any) (*T, error) {
	jsonBytes, err := json.Marshal(d)
	if err != nil {
		fmt.Println("Error marshaling map:", err)
		return nil, err
	}

	var result T
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		fmt.Println("Error unmarshaling map:", err)
		return nil, err
	}
	return &result, nil
}
