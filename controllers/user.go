package controllers

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/21TechLabs/factory-be/dto"
	"github.com/21TechLabs/factory-be/models"
	"github.com/21TechLabs/factory-be/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func UserCreate(role string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		body := c.Body()

		usr := dto.UserCreateDto{}

		err := json.Unmarshal(body, &usr)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   err.Error(),
				"success": false,
			})
		}

		user, err := models.UserCreate(usr, role)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   err.Error(),
				"success": false,
			})
		}

		return SetLoginTokenAndSendResponse(c, user, false)
	}
}

func UserPasswordUpdate(c *fiber.Ctx) error {
	body := c.Body()

	parsedBody := dto.UserPasswordUpdateDto{}

	err := json.Unmarshal(body, &parsedBody)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	currentUser, err := models.UserGetByEmail(parsedBody.Email)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	err = currentUser.CompareAndUpdatePasswordWithToken(parsedBody.Token, parsedBody.Password)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}

func UserMarkForDeletion(c *fiber.Ctx) error {

	var currentUser = c.Locals("user").(models.User)

	err := currentUser.MarkAccountForDeletion()

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Account marked for deletion, account will be deleted in 30 days",
		"success": true,
		"user":    currentUser.GetDetails(false),
	})
}

func UserRequestPasswordResetLink(c *fiber.Ctx) error {
	email := c.Query("email")

	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Email is required",
			"success": false,
		})
	}

	user, err := models.UserGetByEmail(email)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	token, err := user.GeneratePasswordResetToken(true)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	user.PasswordResetToken = token

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}

func UserVerifyEmailToken(c *fiber.Ctx) error {
	token := c.Query("token")

	if token == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "token is required")
	}

	email := c.Query("email")

	if email == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "email is required")
	}

	user, err := models.UserVerifyEmailToken(email, token)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return SetLoginTokenAndSendResponse(c, user, false)
}

func UserLogin(c *fiber.Ctx) error {

	body := c.Body()

	var loginBody dto.UserLoginDto

	err := json.Unmarshal(body, &loginBody)

	if err != nil {
		fmt.Fprintf(os.Stdout, "UserLogin Error: Json Validation Failed.\n%v", err)
		return c.Status(fiber.ErrBadRequest.Code).JSON(bson.M{
			"success": false,
			"message": "UserLogin Error: Json Validation Failed.",
		})
	}

	user, err := models.UserLogin(loginBody)

	if err != nil {
		fmt.Fprintf(os.Stdout, "UserLogin Error: %v\n", err)
		return c.Status(fiber.ErrBadRequest.Code).JSON(bson.M{
			"success": false,
			"message": err.Error(),
		})
	}

	return SetLoginTokenAndSendResponse(c, user, false)
}

func UserLoginVerify(c *fiber.Ctx) error {
	var user = c.Locals("user").(models.User)
	return SetLoginTokenAndSendResponse(c, user, false)
}

func UserLogout(c *fiber.Ctx) error {

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now(),
		HTTPOnly: true,
		Secure:   true,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logged out successfully",
		"success": true,
	})
}
