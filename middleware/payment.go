package middleware

import (
	"time"

	"github.com/21TechLabs/musiclms-backend/models"
	"github.com/21TechLabs/musiclms-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) HasActivePlanAndLevel(minLevel int) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(*models.User)

		if !ok {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "User not found")
		}

		// get app code from query params
		appCode := c.Query("appCode")

		if appCode == "" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "appCode is required")
		}

		activeSubscription, err := m.UserStore.GetActiveAppSubscriptionByAppCode(user, appCode)

		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get active subscription: "+err.Error())
		}

		if activeSubscription.PlanLevel < minLevel {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "User does not have the required plan level for this app")
		}

		// check if current subscription expired or not
		if time.Now().After(activeSubscription.SubscriptionEndsAt) {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "User's current subscription has expired")
		}

		return c.Next()
	}

}
