package middleware

import (
	"log"
	"net/http"

	"github.com/21TechLabs/factory-backend/models"
)

type IMiddleware = func(next http.Handler) http.Handler

type Middleware struct {
	Logger    *log.Logger
	UserStore *models.UserStore
}

func NewMiddleware(logger *log.Logger, userStore *models.UserStore) *Middleware {
	return &Middleware{
		Logger:    logger,
		UserStore: userStore,
	}
}

type MiddlewareStack func(next http.Handler) http.Handler

func (m *Middleware) CreateStack(middlewares ...MiddlewareStack) MiddlewareStack {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

func (m *Middleware) CreateStackWithHandler(middlewares []MiddlewareStack, controller http.HandlerFunc) http.Handler {
	return m.CreateStack(middlewares...)(controller)
}
