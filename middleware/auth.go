package middleware

import (
	"github.com/21TechLabs/factory-be/models"
	"github.com/21TechLabs/factory-be/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

func UserAuthMiddleware(c *fiber.Ctx) error {
	authToken, err := utils.GetToken(c)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized: token not found",
		})
	}

	// Load your secret key from a secure source (e.g., environment variable)
	secretKey := []byte(utils.GetEnv("JWT_SECRET_KEY", false))

	user, err := models.JwtTokenVerifyAndGetUser(authToken, secretKey)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized: invalid token",
		})
	}

	if user.AccountBlocked {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized: account blocked",
			"user":    user.GetDetails(false),
		})
	}

	if user.AccountSuspended {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized: account suspended",
			"user":    user.GetDetails(false),
		})
	}

	if user.AccountDeleted {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized: account deleted",
			"user":    user.GetDetails(false),
		})
	}

	if user.MarkedForDeletion {
		user.MarkedForDeletion = false
		var ctx = mgm.Ctx()
		_, err := mgm.Coll(&user).UpdateOne(ctx, bson.M{
			"_id": user.ID,
		}, bson.M{
			"$set": user,
		})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error: failed to update user",
				"error":   err.Error(),
			})
		}
	}

	// Pass the user to the next handler (add to context)
	c.Locals("user", user)
	return c.Next()
}
