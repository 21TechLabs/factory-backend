package payments

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/21TechLabs/factory-be/dto"
	"github.com/21TechLabs/factory-be/models"
	"github.com/21TechLabs/factory-be/models/payments"
	"github.com/21TechLabs/factory-be/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreatePayment(c *fiber.Ctx) error {
	gateway := c.Params("paymentGateway")
	if gateway == "" {
		log.Printf("Payment gateway create error controller.payments.CreatePayment: %v", "Payment type or product id not found")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Payment type or product id is empty")
	}

	parsedBody := dto.CreateProductDto{}

	err := json.Unmarshal(c.Body(), &parsedBody)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	currentUser, ok := c.Locals("user").(models.User)

	if !ok {
		log.Printf("Payment gateway create error controller.payments.CreatePayment: %v", ok)
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not found")
	}

	// check if user has already subscribed to the product or not

	paymentGateway, err := payments.GetPaymentGateway(gateway, currentUser.ID.Hex())
	if err != nil {
		log.Printf("Payment gateway create error -- Get Payment Gateway controller.payments.CreatePayment: %v", err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	product, err := payments.ProductPlansGetByID(parsedBody.ProductId)

	if err != nil {
		log.Printf("Payment gateway create error -- ProductPlansGetByID controller.payments.CreatePayment: %v", err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	// check if user has an active subscription or not
	userSubscription, err := currentUser.GetActiveAppSubscriptionByAppCode(product.AppCode)

	if err != nil {
		if err != mongo.ErrNoDocuments {
			log.Printf("Payment gateway create error -- GetActiveAppSubscriptionByAppCode controller.payments.CreatePayment: %v", err)
			return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
		}
	}

	if !userSubscription.ID.IsZero() {
		log.Printf("Payment gateway create error -- GetActiveAppSubscriptionByAppCode controller.payments.CreatePayment: %v", "User already has an active subscription")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "User already has an active subscription")
	}

	userSubscription, err = paymentGateway.CreatePayment(product, parsedBody.PlanIdx)

	if err != nil {
		log.Printf("Payment gateway create error -- Create Payment controller.payments.CreatePayment: %v", err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    userSubscription,
	})
}

func UpdatePaymentStatusWebhook(c *fiber.Ctx) error {
	gateway := c.Params("paymentGateway")
	if gateway == "" {
		log.Printf("Payment gateway status update webhook error controller.payments.UpdatePaymentStatus: %v", "Payment type or product id not found")
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"ok": true,
		})
	}

	fmt.Println("received webhook from ", gateway)

	paymentGateway, err := payments.GetPaymentGateway(gateway, "")
	if err != nil {
		log.Printf("Payment gateway status update webhook error -- Get Payment Gateway controller.payments.CreatePayment: %v", err)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"ok": true,
		})
	}

	if paymentGateway.VerifyWebhookSignature(c) != nil {
		log.Printf("Payment gateway status update webhook error -- VerifyWebhookSignature controller.payments.UpdatePaymentStatus: %v", "Webhook signature not verified")
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"ok": true,
		})
	}

	var orderId string

	orderId, err = paymentGateway.GetOrderIdFromWebhookRequest(c.Body())

	if err != nil {
		log.Printf("Payment gateway status update webhook error -- GetOrderIdFromWebhookRequest controller.payments.UpdatePaymentStatus: %v", err)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"ok": true,
		})
	}

	if orderId == "" {
		log.Printf("Payment gateway status update webhook error -- GetOrderIdFromWebhookRequest controller.payments.UpdatePaymentStatus: %v", "Order id not found")
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"ok": true,
		})
	}

	if err = paymentGateway.SetUserViaOrderId(orderId); err != nil {
		log.Printf("Payment gateway status update webhook error -- SetUserViaOrderId controller.payments.UpdatePaymentStatus: %v", err)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"ok": true,
		})
	}

	_, err = paymentGateway.UpdatePaymentStatus(orderId)
	if err != nil {
		log.Printf("Payment gateway status update webhook error -- Update Payment Status controller.payments.UpdatePaymentStatus: %v", err)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"ok": true,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"ok": true,
	})
}
