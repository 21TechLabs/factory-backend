package controllers

import (
	"encoding/json"
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

	usr := dto.UserCreateDto{}

	body := r.Body
	err := json.NewDecoder(body).Decode(&usr)

	if err != nil {
		utils.ErrorResponse(uc.Logger, w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	user, err := uc.UserStore.UserCreate(usr)

	if err != nil {
		// return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
		utils.ErrorResponse(uc.Logger, w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	SetLoginTokenAndSendResponse(uc.Logger, r, w, user, false, uc.UserStore)
}

func (uc *UserController) UserUpdateDto(w http.ResponseWriter, r *http.Request) {

	parsedBody := dto.UserUpdateDto{}

	err := json.NewDecoder(r.Body).Decode(&parsedBody)

	if err != nil {
		utils.ErrorResponse(uc.Logger, w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	currentUser, ok := r.Context().Value(utils.UserContextKey).(*models.User)

	if !ok {
		utils.ErrorResponse(uc.Logger, w, http.StatusBadRequest, []byte("User not found"))
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
	// body := c.Body()

	parsedBody := dto.UserPasswordUpdateDto{}

	err := json.NewDecoder(r.Body).Decode(&parsedBody)

	if err != nil {
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

	currentUser, ok := r.Context().Value(utils.UserContextKey).(*models.User)

	if !ok {
		utils.ErrorResponse(uc.Logger, w, http.StatusBadRequest, []byte("User not found"))
		return
	}

	err := uc.UserStore.MarkAccountForDeletion(currentUser)

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

	var loginBody dto.UserLoginDto

	err := json.NewDecoder(r.Body).Decode(&loginBody)

	if err != nil {
		uc.Logger.Printf("UserLogin Error: Json Validation Failed: %v\n", err)
		utils.ErrorResponse(uc.Logger, w, http.StatusBadRequest, []byte("UserLogin Error: Json Validation Failed."))
		return
	}

	user, err := uc.UserStore.UserLogin(loginBody)

	if err != nil {
		uc.Logger.Printf("UserLogin Error: %v\n", err)
		utils.ErrorResponse(uc.Logger, w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	SetLoginTokenAndSendResponse(uc.Logger, r, w, user, false, uc.UserStore)
}

func (uc *UserController) UserLoginVerify(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(utils.UserContextKey).(*models.User)
	if !ok {
		utils.ErrorResponse(uc.Logger, w, http.StatusBadRequest, []byte("User not found in context"))
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
