package products

import (
	"log"

	"github.com/21TechLabs/factory-be/models"
	"github.com/21TechLabs/factory-be/models/payments"
	"github.com/21TechLabs/factory-be/utils"
	"github.com/gofiber/fiber/v2"
)

func GetProductByAppCode(c *fiber.Ctx) error {
	var appCode string = c.Params("appCode")

	if len(appCode) == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "App code is required!")
	}

	product, err := payments.ProductPlansGetByAppCode(appCode)

	if err != nil {
		log.Printf("Failed to fetch product with app code \"%v\" because an error occured: %v", appCode, err.Error())
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to fetch product!")
	}

	return c.Status(200).JSON(fiber.Map{
		"product": product,
		"success": true,
	})
}

func GetUsersActiveProductSubsctiptionByAppCode(c *fiber.Ctx) error {
	var appCode string = c.Params("appCode")

	if len(appCode) == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "App code is required!")
	}

	user, ok := c.Locals("user").(models.User)

	if !ok {
		log.Printf("Failed to fetch user because user is not found in the context!")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to fetch user!")
	}

	subscription, err := user.GetActiveAppSubscriptionByAppCode(appCode)

	if err != nil {
		log.Printf("Failed to fetch product with app code \"%v\" because an error occured: %v", appCode, err.Error())
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to fetch product!")
	}

	return c.Status(200).JSON(fiber.Map{
		"subscription": subscription,
		"success":      true,
	})
}
