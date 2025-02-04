package middleware

import (
	"time"

	"github.com/21TechLabs/factory-be/models"
	"github.com/gofiber/fiber/v2"
)

func HasActivePlanAndLevel(minLevel int) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(*models.User)

		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized: user not found",
			})
		}

		// get app code from query params
		appCode := c.Query("appCode")

		if appCode == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "appCode is required",
			})
		}

		activeSubscription, err := user.GetActiveAppSubscriptionByAppCode(appCode)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		if activeSubscription.ID.IsZero() {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "No active subscription found",
			})
		}

		if activeSubscription.PlanLevel < minLevel {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "User requires a higher plan level",
			})
		}

		// check if current subscription expired or not
		if time.Now().After(activeSubscription.SubscriptionEndsAt) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Subscription expired",
			})
		}

		return c.Next()
	}

}
