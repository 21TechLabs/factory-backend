package controllers

import (
	"log"
	"time"

	"github.com/21TechLabs/musiclms-backend/models"
	"github.com/21TechLabs/musiclms-backend/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetLoginTokenAndSendResponse(ctx *fiber.Ctx, user models.User, allowPasswordAndResetToken bool, us *models.UserStore) error {

	appCode := ctx.Query("appCode")

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

	var res = fiber.Map{
		"token":   token,
		"user":    us.GetDetails(&user, allowPasswordAndResetToken),
		"success": true,
	}

	if len(appCode) > 0 {
		subscription, err := us.GetActiveAppSubscriptionByAppCode(&user, appCode)

		if err != nil {
			if err != gorm.ErrRecordNotFound {
				log.Printf("Failed to fetch product with app code \"%v\" because an error occured: %v", appCode, err.Error())
				return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
			}
		} else {
			res["subscription"] = subscription
		}
	}

	return ctx.Status(200).JSON(res)
}
