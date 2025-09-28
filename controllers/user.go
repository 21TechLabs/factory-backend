package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/21TechLabs/factory-backend/dto"
	"github.com/21TechLabs/factory-backend/models"
	"github.com/21TechLabs/factory-backend/utils"
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

func (uc *UserController) UserCreate(w http.ResponseWriter, r *http.Request) {
	usr, err := utils.ReadContextValue[*dto.UserCreateDto](r, utils.SchemaValidatorContextKey)
	if err != nil {
		uc.Logger.Printf("UserCreate Error: %v\n", err)
		utils.ErrorResponse(uc.Logger, w, http.StatusBadRequest, []byte(err.Error()))
		return
	}
	user, err := uc.UserStore.UserCreate(*usr)

	if err != nil {
		uc.Logger.Printf("UserCreate Error: %v\n", err)
		utils.ErrorResponse(uc.Logger, w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	SetLoginTokenAndSendResponse(uc.Logger, r, w, user, false, uc.UserStore)
}

func (uc *UserController) UserUpdateDto(w http.ResponseWriter, r *http.Request) {

	parsedBody, err := utils.ReadContextValue[*dto.UserUpdateDto](r, utils.SchemaValidatorContextKey)
	if err != nil {
		uc.Logger.Printf("UserUpdateDto Error: %v\n", err)
		utils.ErrorResponse(uc.Logger, w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	currentUser, err := utils.ReadContextValue[*models.User](r, utils.UserContextKey)
	if err != nil || currentUser == nil {
		utils.ErrorResponse(uc.Logger, w, http.StatusUnauthorized, []byte("User not found"))
		return
	}

	currentUser.Name = parsedBody.Name
	// currentUser.Email = parsedBody.Email

	err = uc.UserStore.Update(currentUser)

	if err != nil {
		utils.ErrorResponse(uc.Logger, w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	utils.ResponseWithJSON(uc.Logger, w, http.StatusOK, utils.Map{
		"success": true,
		"user":    uc.UserStore.GetDetails(currentUser, false),
	})
}

func (uc *UserController) UserPasswordUpdate(w http.ResponseWriter, r *http.Request) {
	parsedBody, err := utils.ReadContextValue[*dto.UserPasswordUpdateDto](r, utils.SchemaValidatorContextKey)
	if err != nil {
		uc.Logger.Printf("UserPasswordUpdate Error: %v\n", err)
		utils.ErrorResponse(uc.Logger, w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	currentUser, err := uc.UserStore.UserGetByEmail(parsedBody.Email)

	if err != nil {
		utils.ErrorResponse(uc.Logger, w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	err = uc.UserStore.CompareAndUpdatePasswordWithToken(&currentUser, parsedBody.Token, parsedBody.Password)

	if err != nil {
		utils.ErrorResponse(uc.Logger, w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	// return c.Status(fiber.StatusOK).JSON(fiber.Map{
	// 	"success": true,
	// })
	utils.ResponseWithJSON(uc.Logger, w, http.StatusOK, utils.Map{
		"success": true,
		"user":    uc.UserStore.GetDetails(&currentUser, false),
	})
}

func (uc *UserController) UserMarkForDeletion(w http.ResponseWriter, r *http.Request) {
	currentUser, err := utils.ReadContextValue[*models.User](r, utils.UserContextKey)
	if err != nil || currentUser == nil {
		utils.ErrorResponse(uc.Logger, w, http.StatusUnauthorized, []byte("User not found"))
		return
	}

	err = uc.UserStore.MarkAccountForDeletion(currentUser)

	if err != nil {
		utils.ErrorResponse(uc.Logger, w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	utils.ResponseWithJSON(uc.Logger, w, http.StatusOK, utils.Map{
		"message": "Account marked for deletion, account will be deleted in 30 days",
		"success": true,
		"user":    uc.UserStore.GetDetails(currentUser, false),
	})
}

func (uc *UserController) UserRequestPasswordResetLink(w http.ResponseWriter, r *http.Request) {
	// email := c.Query("email")
	email := r.URL.Query().Get("email")

	if email == "" {
		utils.ErrorResponse(uc.Logger, w, http.StatusBadRequest, []byte("email is required"))
		return
	}

	user, err := uc.UserStore.UserGetByEmail(email)

	if err != nil {
		utils.ErrorResponse(uc.Logger, w, http.StatusBadRequest, []byte("user not found"))
		return
	}

	token, err := uc.UserStore.GeneratePasswordResetToken(&user, true)

	if err != nil {
		utils.ErrorResponse(uc.Logger, w, http.StatusInternalServerError, []byte("Failed to generate password reset token!"))
		return
	}

	user.PasswordResetToken = token

	utils.ResponseWithJSON(uc.Logger, w, http.StatusOK, utils.Map{
		"success": true,
		"message": "Password reset link sent to your email",
		"user":    uc.UserStore.GetDetails(&user, false),
	})
}

func (uc *UserController) UserVerifyEmailToken(w http.ResponseWriter, r *http.Request) {
	// token := c.Query("token")
	token := r.URL.Query().Get("token")

	if token == "" {
		utils.ErrorResponse(uc.Logger, w, http.StatusBadRequest, []byte("token is required"))
		return
	}

	email := r.URL.Query().Get("email")

	if email == "" {
		utils.ErrorResponse(uc.Logger, w, http.StatusBadRequest, []byte("email is required"))
		return
	}

	user, err := uc.UserStore.UserVerifyEmailToken(email, token)

	if err != nil {
		utils.ErrorResponse(uc.Logger, w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	SetLoginTokenAndSendResponse(uc.Logger, r, w, user, false, uc.UserStore)
}

func (uc *UserController) UserLogin(w http.ResponseWriter, r *http.Request) {

	loginBody, err := utils.ReadContextValue[*dto.UserLoginDto](r, utils.SchemaValidatorContextKey)
	if err != nil {
		uc.Logger.Printf("UserLogin Error: %v\n", err)
		utils.ErrorResponse(uc.Logger, w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	user, err := uc.UserStore.UserLogin(*loginBody)

	if err != nil {
		uc.Logger.Printf("UserLogin Error: %v\n", err)
		utils.ErrorResponse(uc.Logger, w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	SetLoginTokenAndSendResponse(uc.Logger, r, w, user, false, uc.UserStore)
}

func (uc *UserController) UserLoginVerify(w http.ResponseWriter, r *http.Request) {
	user, err := utils.ReadContextValue[*models.User](r, utils.UserContextKey)
	if err != nil || user == nil {
		utils.ErrorResponse(uc.Logger, w, http.StatusUnauthorized, []byte("User not found"))
		return
	}

	SetLoginTokenAndSendResponse(uc.Logger, r, w, *user, false, uc.UserStore)
}

func (uc *UserController) UserLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
		Secure:   true,
	})

	utils.ResponseWithJSON(uc.Logger, w, http.StatusOK, utils.Map{
		"success": true,
		"message": "User logged out successfully",
	})
}
