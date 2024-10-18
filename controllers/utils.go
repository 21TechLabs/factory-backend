package controllers

import (
	"log"
	"time"

	"github.com/21TechLabs/factory-be/models"
	"github.com/21TechLabs/factory-be/utils"
	"github.com/gofiber/fiber/v2"
)

func SetLoginTokenAndSendResponse(ctx *fiber.Ctx, user models.User, passwordResetToekn bool) {
	if len(user.Email) == 0 {
		log.Default().Printf("Failed to fetch user \"%v\" because token does not exists!", user.Email)
		ctx.Cookie(&fiber.Cookie{
			Name:     "token",
			Value:    "",
			Expires:  time.Now(),
			HTTPOnly: true,
			Secure:   true,
		})
		ctx.Status(401).SendString("Failed to login user!")
		return
	}

	var expiresAfter time.Time = time.Now().Add(time.Hour * 24 * 5)
	token, err := user.JwtTokenGet(expiresAfter, []byte(utils.GetEnv("JWT_SECRET_KEY", false)))

	if err != nil {
		log.Default().Printf("Failed to create login token for the user %v because an error occured: %v", user.Email, err.Error())
		ctx.Status(401).SendString("Failed to generate login token")
		return
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  expiresAfter,
		HTTPOnly: true,
		Secure:   true,
	})

	ctx.Status(200).JSON(fiber.Map{
		"token":   token,
		"user":    user.GetDetails(passwordResetToekn),
		"success": true,
	})
}
