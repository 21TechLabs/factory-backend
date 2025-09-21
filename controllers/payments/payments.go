package payments_controller

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/21TechLabs/musiclms-backend/dto"
	"github.com/21TechLabs/musiclms-backend/models"
	"github.com/21TechLabs/musiclms-backend/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PaymentsController struct {
	Logger                *log.Logger
	ProductPlanStore      *models.ProductPlanStore
	UserStore             *models.UserStore
	UserSubscriptionStore *models.UserSubscriptionStore
}

func NewPaymentsController(log *log.Logger, store *models.ProductPlanStore, us *models.UserStore, uss *models.UserSubscriptionStore) *PaymentsController {
	return &PaymentsController{
		Logger:                log,
		ProductPlanStore:      store,
		UserStore:             us,
		UserSubscriptionStore: uss,
	}
}

func (pc *PaymentsController) CreatePayment(c *fiber.Ctx) error {
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

	paymentGateway, err := pc.ProductPlanStore.GetPaymentGateway(gateway, currentUser.ID)
	if err != nil {
		log.Printf("Payment gateway create error -- Get Payment Gateway controller.payments.CreatePayment: %v", err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	product, err := pc.ProductPlanStore.ProductPlanGetByID(parsedBody.ProductId)

	if err != nil {
		log.Printf("Payment gateway create error -- ProductPlansGetByID controller.payments.CreatePayment: %v", err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	// check if user has an active subscription or not
	userSubscription, err := pc.UserStore.GetActiveAppSubscriptionByAppCode(&currentUser, product.AppCode)

	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Printf("Payment gateway create error -- GetActiveAppSubscriptionByAppCode controller.payments.CreatePayment: %v", err)
			return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
		}
	}

	if userSubscription.Status == models.SubscriptionStatusActive {
		log.Printf("Payment gateway create error -- GetActiveAppSubscriptionByAppCode controller.payments.CreatePayment: %v", "User already has an active subscription")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "User already has an active subscription")
	}

	userSubscription, err = paymentGateway.CreatePayment(pc.UserSubscriptionStore, product, parsedBody.PlanIdx)

	if err != nil {
		log.Printf("Payment gateway create error -- Create Payment controller.payments.CreatePayment: %v", err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    userSubscription,
	})
}

func (pc *PaymentsController) UpdatePaymentStatusWebhook(c *fiber.Ctx) error {
	gateway := c.Params("paymentGateway")
	if gateway == "" {
		log.Printf("Payment gateway status update webhook error controller.payments.UpdatePaymentStatus: %v", "Payment type or product id not found")
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"ok": true,
		})
	}

	fmt.Println("received webhook from ", gateway)

	// create a zero object of the payment gateway
	paymentGateway, err := pc.ProductPlanStore.GetPaymentGateway(gateway, 0)
	if err != nil {
		log.Printf("Payment gateway status update webhook error -- Get Payment Gateway controller.payments.CreatePayment: %v", err)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"ok": true,
		})
	}

	// if paymentGateway.VerifyWebhookSignature(c) != nil {
	// 	log.Printf("Payment gateway status update webhook error -- VerifyWebhookSignature controller.payments.UpdatePaymentStatus: %v", "Webhook signature not verified")
	// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
	// 		"ok": true,
	// 	})
	// }

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

	if err = paymentGateway.SetUserViaOrderId(pc.UserSubscriptionStore, orderId); err != nil {
		log.Printf("Payment gateway status update webhook error -- SetUserViaOrderId controller.payments.UpdatePaymentStatus: %v", err)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"ok": true,
		})
	}

	_, err = paymentGateway.UpdatePaymentStatus(pc.UserSubscriptionStore, orderId)
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

func (pc *PaymentsController) GetProductPlansByAppCode(c *fiber.Ctx) error {
	appCode := c.Params("appCode")
	if appCode == "" {
		log.Printf("Get Product Plans by App Code error controller.payments.GetProductPlansByAppCode: %v", "App code not found")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "App code is empty")
	}

	plans, err := pc.ProductPlanStore.GetByAppCode(appCode)
	if err != nil {
		log.Printf("Get Product Plans by App Code error controller.payments.GetProductPlansByAppCode: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(plans)
}
