package middleware

import (
	"context"
	"net/http"

	"github.com/21TechLabs/factory-backend/models"
	"github.com/21TechLabs/factory-backend/utils"
)

func (m *Middleware) HasRoleMiddleware(whiteListedRoles []models.UserRole) IMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := utils.ReadContextValue[*models.User](r, utils.UserContextKey)

			if err != nil || user == nil {
				utils.ErrorResponse(m.Logger, w, http.StatusUnauthorized, []byte("Unauthorized: user not found"))
				return
			}

			for _, role := range whiteListedRoles {
				if user.Role == role {
					r = r.WithContext(context.WithValue(r.Context(), utils.UserContextKey, user)) // Store user in context for further use
					next.ServeHTTP(w, r)
					return
				}
			}
			utils.ErrorResponse(m.Logger, w, http.StatusForbidden, []byte("Forbidden: insufficient permissions"))
		})
	}
}
