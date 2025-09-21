package products_controller

import (
	"log"

	"github.com/21TechLabs/musiclms-backend/models"
	"github.com/21TechLabs/musiclms-backend/utils"
	"github.com/gofiber/fiber/v2"
)

type ProductPlanController struct {
	Logger           *log.Logger
	ProductPlanStore *models.ProductPlanStore
	UserStore        *models.UserStore
}

func NewProductPlanController(log *log.Logger, store *models.ProductPlanStore, us *models.UserStore) *ProductPlanController {
	return &ProductPlanController{
		Logger:           log,
		ProductPlanStore: store,
		UserStore:        us,
	}
}

func (ppc *ProductPlanController) GetProductByAppCode(c *fiber.Ctx) error {
	var appCode string = c.Params("appCode")

	if len(appCode) == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "App code is required!")
	}

	product, err := ppc.ProductPlanStore.GetByAppCode(appCode)

	if err != nil {
		log.Printf("Failed to fetch product with app code \"%v\" because an error occured: %v", appCode, err.Error())
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to fetch product!")
	}

	return c.Status(200).JSON(fiber.Map{
		"product": product,
		"success": true,
	})
}

func (ppc *ProductPlanController) GetUsersActiveProductSubsctiptionByAppCode(c *fiber.Ctx) error {
	var appCode string = c.Params("appCode")

	if len(appCode) == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "App code is required!")
	}

	user, ok := c.Locals("user").(models.User)

	if !ok {
		log.Printf("Failed to fetch user because user is not found in the context!")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to fetch user!")
	}

	subscription, err := ppc.UserStore.GetActiveAppSubscriptionByAppCode(&user, appCode)

	if err != nil {
		log.Printf("Failed to fetch product with app code \"%v\" because an error occured: %v", appCode, err.Error())
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to fetch product!")
	}

	return c.Status(200).JSON(fiber.Map{
		"subscription": subscription,
		"success":      true,
	})
}
