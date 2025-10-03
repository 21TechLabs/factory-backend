package middleware

import (
	"context"
	"net/http"

	"github.com/21TechLabs/factory-backend/utils"
)

func (m *Middleware) UserAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken, err := utils.GetToken(r)

		if err != nil {
			utils.ErrorResponse(m.Logger, w, http.StatusUnauthorized, []byte("Unauthorized: "+err.Error()))
			return
		}

		secretKey := []byte(utils.GetEnv("JWT_SECRET_KEY", false))

		user, err := m.UserStore.JwtTokenVerifyAndGetUser(authToken, secretKey)
		if err != nil {
			utils.ErrorResponse(m.Logger, w, http.StatusUnauthorized, []byte("Unauthorized: "+err.Error()))
			return
		}

		if user.AccountBlocked {
			utils.ErrorResponse(m.Logger, w, http.StatusUnauthorized, []byte("Unauthorized: account blocked"))
			return
		}

		if user.AccountSuspended {
			utils.ErrorResponse(m.Logger, w, http.StatusUnauthorized, []byte("Unauthorized: account suspended"))
			return
		}

		if user.AccountDeleted {
			utils.ErrorResponse(m.Logger, w, http.StatusUnauthorized, []byte("Unauthorized: account deleted"))
			return
		}

		if user.MarkedForDeletion {
			user.MarkedForDeletion = false
			err := m.UserStore.Update(&user)
			if err != nil {
				utils.ErrorResponse(m.Logger, w, http.StatusInternalServerError, []byte("Failed to update user: "+err.Error()))
				return
			}
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, utils.UserContextKey, &user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
