package middleware

import (
	"github.com/21TechLabs/musiclms-backend/models"
	"github.com/21TechLabs/musiclms-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) HasRoleMiddleware(whiteListedRoles []models.UserRole) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(models.User)

		if !ok {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized: user not found")
		}

		for _, role := range whiteListedRoles {
			if user.Role == role {
				c.Locals("user", user) // Store user in context for further use
				return c.Next()
			}
		}
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Forbidden: insufficient permissions")
	}
}
