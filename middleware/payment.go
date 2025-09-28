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
			user, err := utils.ReadContextValue[*models.User](r, utils.UserContextKey)

			if err != nil || user == nil {
				utils.ErrorResponse(m.Logger, w, http.StatusNotFound, []byte("User not found"))
				return
			}

			// get app code from query params
			appCode := r.URL.Query().Get("appCode")

			if appCode == "" {
				utils.ErrorResponse(m.Logger, w, http.StatusBadRequest, []byte("appCode is required"))
				return
			}

			activeSubscription, err := m.UserStore.GetActiveAppSubscriptionByAppCode(user, appCode)

			if err != nil {
				utils.ErrorResponse(m.Logger, w, http.StatusInternalServerError, []byte("Failed to get active subscription: "+err.Error()))
				return
			}

			if activeSubscription.PlanLevel < minLevel {
				utils.ErrorResponse(m.Logger, w, http.StatusBadRequest, []byte("User does not have the required plan level for this app"))
				return
			}

			// check if current subscription expired or not
			if time.Now().After(activeSubscription.SubscriptionEndsAt) {
				utils.ErrorResponse(m.Logger, w, http.StatusBadRequest, []byte("User's current subscription has expired"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
