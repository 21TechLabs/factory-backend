package payments

import (
	"encoding/json"
	"log"

	"github.com/21TechLabs/factory-be/dto"
	"github.com/21TechLabs/factory-be/models"
	"github.com/21TechLabs/factory-be/models/payments"
	"github.com/21TechLabs/factory-be/utils"
	"github.com/gofiber/fiber/v2"
)

func CreatePayment(c *fiber.Ctx) error {
	paymentType := c.Params("paymentType")
	if paymentType == "" {
		log.Printf("Payment gateway create error: %v", "Payment type or product id not found")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Payment type or product id is empty")
	}

	parsedBody := dto.CreateProductDto{}

	err := json.Unmarshal(c.Body(), &parsedBody)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	currentUser, ok := c.Locals("user").(models.User)

	if !ok {
		log.Printf("Payment gateway create error: %v", ok)
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not found")
	}

	payment, err := payments.GetPaymentGateway(paymentType, currentUser.ID.Hex())

	if err != nil {
		log.Printf("Payment gateway create error -- Get Payment Gateway: %v", err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	product, err := payments.ProductPlansGetByID(parsedBody.ProductId)

	if err != nil {
		log.Printf("Payment gateway create error -- ProductPlansGetByID: %v", err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	userSubscription, err := payment.CreatePayment(product, parsedBody.PlanIdx)

	if err != nil {
		log.Printf("Payment gateway create error -- Create Payment: %v", err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    userSubscription,
	})
}
