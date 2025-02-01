package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
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

func GetToken(c *fiber.Ctx) (string, error) {
	authToken := c.Cookies("token")

	if authToken == "" {
		authToken = c.Get("Authorization")
		if authToken != "" {
			if len(authToken) > 10 {
				authToken = authToken[7:]
			}
		}

	}

	if authToken == "" {
		var payload = struct {
			Token string `json:"token"`
		}{}
		if err := c.BodyParser(&payload); err != nil {
			return "", RaiseError{Message: "Token not found in headers or body"}
		}
		authToken = payload.Token
	}

	if authToken == "" {
		return "", RaiseError{Message: "Token not found in headers or body"}
	}
	return authToken, nil
}

func IsValidObjectID(id string) bool {
	_, err := primitive.ObjectIDFromHex(id)
	return err == nil
}
