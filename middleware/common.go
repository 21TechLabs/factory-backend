package middleware

import (
	"log"

	"github.com/21TechLabs/musiclms-backend/models"
)

type Middleware struct {
	Logger    *log.Logger
	UserStore *models.UserStore
}

func NewMiddleware(log *log.Logger, userStore *models.UserStore) *Middleware {
	return &Middleware{
		Logger:    log,
		UserStore: userStore,
	}
}
