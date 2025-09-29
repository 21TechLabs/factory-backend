package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/21TechLabs/factory-backend/models"
	"github.com/21TechLabs/factory-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func SetLoginTokenAndSendResponse(log *log.Logger, r *http.Request, w http.ResponseWriter, user models.User, allowPasswordAndResetToken bool, us *models.UserStore) {

	if len(user.Email) == 0 {
		log.Printf("Failed to fetch user \"%v\" because token does not exists!", user.Email)
		w.Header().Set("Content-Type", "application/json")
		// clear the cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    "",
			Expires:  time.Now(),
			HttpOnly: true,
			Secure:   true,
		})

		utils.ErrorResponse(log, w, http.StatusBadRequest, []byte("User not found!"))
		return
	}

	var expiresAfter time.Time = time.Now().Add(time.Hour * 24 * 5)
	token, err := user.JwtTokenGet(expiresAfter, []byte(utils.GetEnv("JWT_SECRET_KEY", false)))

	if err != nil {
		log.Printf("Failed to create login token for the user %v because an error occured: %v", user.Email, err.Error())
		utils.ErrorResponse(log, w, http.StatusBadRequest, []byte("Failed to generate login token!"))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  expiresAfter,
		HttpOnly: true,
		Secure:   true,
	})

	var res = fiber.Map{
		"token":   token,
		"user":    us.GetDetails(&user, allowPasswordAndResetToken),
		"success": true,
	}

	utils.ResponseWithJSON(log, w, http.StatusOK, res)
}
