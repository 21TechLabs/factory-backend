package middleware

import (
	"net/http"
	"time"

	"github.com/21TechLabs/factory-backend/models"
	"github.com/21TechLabs/factory-backend/utils"
)

func (m *Middleware) HasActivePlanAndLevel(minLevel int) IMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(utils.UserContextKey).(*models.User)

			if !ok {
				// return utils.ErrorResponse(c, fiber.StatusBadRequest, "User not found")
				utils.ErrorResponse(m.Logger, w, http.StatusNotFound, []byte("User not found"))
				return
			}

			// get app code from query params
			// appCode := c.Query("appCode")
			appCode := r.URL.Query().Get("appCode")

			if appCode == "" {
				// return utils.ErrorResponse(c, fiber.StatusBadRequest, "appCode is required")
				utils.ErrorResponse(m.Logger, w, http.StatusBadRequest, []byte("appCode is required"))
				return
			}

			activeSubscription, err := m.UserStore.GetActiveAppSubscriptionByAppCode(user, appCode)

			if err != nil {
				// return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get active subscription: "+err.Error())
				utils.ErrorResponse(m.Logger, w, http.StatusInternalServerError, []byte("Failed to get active subscription: "+err.Error()))
				return
			}

			if activeSubscription.PlanLevel < minLevel {
				// return utils.ErrorResponse(c, fiber.StatusBadRequest, "User does not have the required plan level for this app")
				utils.ErrorResponse(m.Logger, w, http.StatusBadRequest, []byte("User does not have the required plan level for this app"))
				return
			}

			// check if current subscription expired or not
			if time.Now().After(activeSubscription.SubscriptionEndsAt) {
				// return utils.ErrorResponse(c, fiber.StatusBadRequest, "User's current subscription has expired")4
				utils.ErrorResponse(m.Logger, w, http.StatusBadRequest, []byte("User's current subscription has expired"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
