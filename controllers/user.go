package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/21TechLabs/musiclms-backend/dto"
	"github.com/21TechLabs/musiclms-backend/models"
	"github.com/21TechLabs/musiclms-backend/utils"
	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	Logger    *log.Logger
	UserStore *models.UserStore
}

func NewUserController(logger *log.Logger, store *models.UserStore) *UserController {
	return &UserController{
		Logger:    logger,
		UserStore: store,
	}
}

func (uc *UserController) UserCreate(c *fiber.Ctx) error {
	body := c.Body()

	usr := dto.UserCreateDto{}

	err := json.Unmarshal(body, &usr)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	var isSubdomain = len(strings.Split(c.Get("Origin"), ".")) == 3

	user, err := uc.UserStore.UserCreate(usr, isSubdomain)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return SetLoginTokenAndSendResponse(c, user, false, uc.UserStore)
}

func (uc *UserController) UserUpdateDto(c *fiber.Ctx) error {
	body := c.Body()

	parsedBody := dto.UserUpdateDto{}

	if err := json.Unmarshal(body, &parsedBody); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	currentUser, ok := c.Locals("user").(models.User)

	if !ok {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "User not found")
	}

	currentUser.Name = parsedBody.Name
	// currentUser.Email = parsedBody.Email

	err := uc.UserStore.Update(&currentUser)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"user":    uc.UserStore.GetDetails(&currentUser, false),
	})
}

func (uc *UserController) UserPasswordUpdate(c *fiber.Ctx) error {
	body := c.Body()

	parsedBody := dto.UserPasswordUpdateDto{}

	err := json.Unmarshal(body, &parsedBody)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	currentUser, err := uc.UserStore.UserGetByEmail(parsedBody.Email)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	err = uc.UserStore.CompareAndUpdatePasswordWithToken(&currentUser, parsedBody.Token, parsedBody.Password)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}

func (uc *UserController) UserMarkForDeletion(c *fiber.Ctx) error {

	var currentUser = c.Locals("user").(models.User)

	err := uc.UserStore.MarkAccountForDeletion(&currentUser)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Account marked for deletion, account will be deleted in 30 days",
		"success": true,
		"user":    uc.UserStore.GetDetails(&currentUser, false),
	})
}

func (uc *UserController) UserRequestPasswordResetLink(c *fiber.Ctx) error {
	email := c.Query("email")

	if email == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "email is required")
	}

	user, err := uc.UserStore.UserGetByEmail(email)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "user not found")
	}

	token, err := uc.UserStore.GeneratePasswordResetToken(&user, true)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to generate password reset token")
	}

	user.PasswordResetToken = token

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}

func (uc *UserController) UserVerifyEmailToken(c *fiber.Ctx) error {
	token := c.Query("token")

	if token == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "token is required")
	}

	email := c.Query("email")

	if email == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "email is required")
	}

	user, err := uc.UserStore.UserVerifyEmailToken(email, token)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return SetLoginTokenAndSendResponse(c, user, false, uc.UserStore)
}

func (uc *UserController) UserLogin(c *fiber.Ctx) error {

	body := c.Body()

	var loginBody dto.UserLoginDto

	err := json.Unmarshal(body, &loginBody)

	if err != nil {
		fmt.Fprintf(os.Stdout, "UserLogin Error: Json Validation Failed.\n%v", err)
		return c.Status(fiber.ErrBadRequest.Code).JSON(map[string]interface{}{
			"success": false,
			"message": "UserLogin Error: Json Validation Failed.",
		})
	}

	user, err := uc.UserStore.UserLogin(loginBody)

	if err != nil {
		fmt.Fprintf(os.Stdout, "UserLogin Error: %v\n", err)
		return c.Status(fiber.ErrBadRequest.Code).JSON(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return SetLoginTokenAndSendResponse(c, user, false, uc.UserStore)
}

func (uc *UserController) UserLoginVerify(c *fiber.Ctx) error {
	var user = c.Locals("user").(models.User)
	return SetLoginTokenAndSendResponse(c, user, false, uc.UserStore)
}

func (uc *UserController) UserLogout(c *fiber.Ctx) error {

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
