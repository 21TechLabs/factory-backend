package controllers

import (
	"log"
	"time"

	"github.com/21TechLabs/factory-be/models"
	"github.com/21TechLabs/factory-be/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetLoginTokenAndSendResponse(ctx *fiber.Ctx, user models.User, allowPasswordAndResetToken bool) error {
	if len(user.Email) == 0 {
		log.Default().Printf("Failed to fetch user \"%v\" because token does not exists!", user.Email)
		ctx.Cookie(&fiber.Cookie{
			Name:     "token",
			Value:    "",
			Expires:  time.Now(),
			HTTPOnly: true,
			Secure:   true,
		})

		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, "failed to login user!")
	}

	var expiresAfter time.Time = time.Now().Add(time.Hour * 24 * 5)
	token, err := user.JwtTokenGet(expiresAfter, []byte(utils.GetEnv("JWT_SECRET_KEY", false)))

	if err != nil {
		log.Printf("Failed to create login token for the user %v because an error occured: %v", user.Email, err.Error())
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, "Failed to generate login token!")
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  expiresAfter,
		HTTPOnly: true,
		Secure:   true,
	})

	appCode := ctx.Query("appCode")

	var res = fiber.Map{
		"token":   token,
		"user":    user.GetDetails(allowPasswordAndResetToken),
		"success": true,
	}

	if len(appCode) > 0 {
		subscription, err := user.GetActiveAppSubscriptionByAppCode(appCode)

		if err != nil {
			if err != mongo.ErrNoDocuments {
				log.Printf("Failed to fetch product with app code \"%v\" because an error occured: %v", appCode, err.Error())
				return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
			}
		} else {
			res["subscription"] = subscription
		}
	}

	return ctx.Status(200).JSON(res)
}
