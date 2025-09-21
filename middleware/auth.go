package middleware

import (
	"github.com/21TechLabs/musiclms-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) UserAuthMiddleware(c *fiber.Ctx) error {
	authToken, err := utils.GetToken(c)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized: "+err.Error())
	}

	secretKey := []byte(utils.GetEnv("JWT_SECRET_KEY", false))

	user, err := m.UserStore.JwtTokenVerifyAndGetUser(authToken, secretKey)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized: "+err.Error())
	}

	if user.AccountBlocked {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized: account blocked")
	}

	if user.AccountSuspended {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized: account suspended")
	}

	if user.AccountDeleted {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized: account deleted")
	}

	if user.MarkedForDeletion {
		user.MarkedForDeletion = false
		err := m.UserStore.Update(&user)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update user: "+err.Error())
		}
	}

	// Pass the user to the next handler (add to context)
	c.Locals("user", user)
	return c.Next()
}
