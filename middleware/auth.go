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
			m.Logger.Println("Error getting auth token:", err)
			utils.ErrorResponse(m.Logger, w, http.StatusUnauthorized, []byte("Unauthorized: "+err.Error()))
			return
		}

		secretKey := []byte(utils.GetEnv("JWT_SECRET_KEY", false))

		user, err := m.UserStore.JwtTokenVerifyAndGetUser(authToken, secretKey)
		if err != nil {
			m.Logger.Println("Error verifying JWT token:", err)
			utils.ErrorResponse(m.Logger, w, http.StatusUnauthorized, []byte("Unauthorized: "+err.Error()))
			return
		}

		if user.AccountBlocked {
			m.Logger.Println("User account is blocked:", user.Name, "ID:", user.ID)
			utils.ErrorResponse(m.Logger, w, http.StatusUnauthorized, []byte("Unauthorized: account blocked"))
			return
		}

		if user.AccountSuspended {
			m.Logger.Println("User account is suspended:", user.Name, "ID:", user.ID)
			utils.ErrorResponse(m.Logger, w, http.StatusUnauthorized, []byte("Unauthorized: account suspended"))
			return
		}

		if user.AccountDeleted {
			m.Logger.Println("User account is deleted:", user.Name, "ID:", user.ID)
			utils.ErrorResponse(m.Logger, w, http.StatusUnauthorized, []byte("Unauthorized: account deleted"))
			return
		}

		if user.MarkedForDeletion {
			user.MarkedForDeletion = false
			err := m.UserStore.Update(&user)
			if err != nil {
				m.Logger.Println("Error updating user marked for deletion:", err)
				utils.ErrorResponse(m.Logger, w, http.StatusInternalServerError, []byte("Failed to update user: "+err.Error()))
				return
			}
		}

		r = r.WithContext(
			context.WithValue(r.Context(), utils.UserContextKey, &user),
		)

		next.ServeHTTP(w, r)
	})
}
