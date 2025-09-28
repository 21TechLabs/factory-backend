package payments_controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/21TechLabs/factory-backend/dto"
	"github.com/21TechLabs/factory-backend/models"
	"github.com/21TechLabs/factory-backend/utils"
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

func (pc *PaymentsController) CreatePayment(w http.ResponseWriter, r *http.Request) {
	gateway := r.PathValue("paymentGateway")
	if gateway == "" {
		log.Printf("Payment gateway create error controller.payments.CreatePayment: %v", "Payment type or product id not found")
		utils.ErrorResponse(pc.Logger, w, http.StatusBadRequest, []byte("Payment type or product id is empty"))
		return
	}

	// var parsedBody *dto.CreateProductDto = r.Context().Value(utils.SchemaValidatorContextKey).(*dto.CreateProductDto)
	parsedBody, err := utils.ReadContextValue[dto.CreateProductDto](r, utils.SchemaValidatorContextKey)

	if err != nil {
		log.Printf("Payment gateway create error controller.payments.CreatePayment: %v", err)
		utils.ErrorResponse(pc.Logger, w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	currentUser, ok := r.Context().Value("user").(models.User)

	if !ok {
		log.Printf("Payment gateway create error controller.payments.CreatePayment: %v", ok)
		utils.ErrorResponse(pc.Logger, w, http.StatusUnauthorized, []byte("User not found"))
		return
	}

	// check if user has already subscribed to the product or not

	paymentGateway, err := pc.ProductPlanStore.GetPaymentGateway(gateway, currentUser.ID)
	if err != nil {
		log.Printf("Payment gateway create error -- Get Payment Gateway controller.payments.CreatePayment: %v", err)
		utils.ErrorResponse(pc.Logger, w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	product, err := pc.ProductPlanStore.ProductPlanGetByID(parsedBody.ProductId)

	if err != nil {
		log.Printf("Payment gateway create error -- ProductPlansGetByID controller.payments.CreatePayment: %v", err)
		utils.ErrorResponse(pc.Logger, w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	// check if user has an active subscription or not
	userSubscription, err := pc.UserStore.GetActiveAppSubscriptionByAppCode(&currentUser, product.AppCode)

	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Printf("Payment gateway create error -- GetActiveAppSubscriptionByAppCode controller.payments.CreatePayment: %v", err)
			utils.ErrorResponse(pc.Logger, w, http.StatusBadRequest, []byte(err.Error()))
			return
		}
	}

	if userSubscription.Status == models.SubscriptionStatusActive {
		log.Printf("Payment gateway create error -- GetActiveAppSubscriptionByAppCode controller.payments.CreatePayment: %v", "User already has an active subscription")
		utils.ErrorResponse(pc.Logger, w, http.StatusBadRequest, []byte("User already has an active subscription"))
		return
	}

	userSubscription, err = paymentGateway.CreatePayment(pc.UserSubscriptionStore, product, parsedBody.PlanIdx)

	if err != nil {
		log.Printf("Payment gateway create error -- Create Payment controller.payments.CreatePayment: %v", err)
		utils.ErrorResponse(pc.Logger, w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	utils.ResponseWithJSON(pc.Logger, w, http.StatusOK, utils.Map{
		"success": true,
		"data":    userSubscription,
	})

}

func (pc *PaymentsController) UpdatePaymentStatusWebhook(w http.ResponseWriter, r *http.Request) {
	// gateway := c.Params("paymentGateway")
	gateway := r.PathValue("paymentGateway")
	if gateway == "" {
		log.Printf("Payment gateway status update webhook error controller.payments.UpdatePaymentStatus: %v", "Payment type or product id not found")
		utils.ResponseWithJSON(pc.Logger, w, http.StatusBadRequest, utils.Map{
			"ok": true,
		})
		return
	}

	fmt.Println("received webhook from ", gateway)

	// create a zero object of the payment gateway
	paymentGateway, err := pc.ProductPlanStore.GetPaymentGateway(gateway, 0)
	if err != nil {
		log.Printf("Payment gateway status update webhook error -- Get Payment Gateway controller.payments.CreatePayment: %v", err)
		utils.ResponseWithJSON(pc.Logger, w, http.StatusBadRequest, utils.Map{
			"ok": true,
		})
		return
	}

	if paymentGateway.VerifyWebhookSignature(r) != nil {
		log.Printf("Payment gateway status update webhook error -- VerifyWebhookSignature controller.payments.UpdatePaymentStatus: %v", "Webhook signature not verified")
		utils.ResponseWithJSON(pc.Logger, w, http.StatusBadRequest, utils.Map{
			"ok": true,
		})
		return
	}

	var orderId string

	bodyBytes := make([]byte, r.ContentLength)
	_, err = r.Body.Read(bodyBytes)
	if err != nil {
		log.Printf("Payment gateway status update webhook error -- Read body controller.payments.UpdatePaymentStatus: %v", err)
		utils.ResponseWithJSON(pc.Logger, w, http.StatusBadRequest, utils.Map{
			"ok": true,
		})
		return
	}
	defer r.Body.Close()

	orderId, err = paymentGateway.GetOrderIdFromWebhookRequest(bodyBytes)

	if err != nil {
		log.Printf("Payment gateway status update webhook error -- GetOrderIdFromWebhookRequest controller.payments.UpdatePaymentStatus: %v", err)
		utils.ResponseWithJSON(pc.Logger, w, http.StatusBadRequest, utils.Map{
			"ok": true,
		})
		return
	}

	if orderId == "" {
		log.Printf("Payment gateway status update webhook error -- GetOrderIdFromWebhookRequest controller.payments.UpdatePaymentStatus: %v", "Order id not found")
		utils.ResponseWithJSON(pc.Logger, w, http.StatusBadRequest, utils.Map{
			"ok": true,
		})
		return
	}

	if err = paymentGateway.SetUserViaOrderId(pc.UserSubscriptionStore, orderId); err != nil {
		log.Printf("Payment gateway status update webhook error -- SetUserViaOrderId controller.payments.UpdatePaymentStatus: %v", err)
		utils.ResponseWithJSON(pc.Logger, w, http.StatusBadRequest, utils.Map{
			"ok": true,
		})
		return
	}

	_, err = paymentGateway.UpdatePaymentStatus(pc.UserSubscriptionStore, orderId)
	if err != nil {
		log.Printf("Payment gateway status update webhook error -- Update Payment Status controller.payments.UpdatePaymentStatus: %v", err)
		utils.ResponseWithJSON(pc.Logger, w, http.StatusBadRequest, utils.Map{
			"ok": true,
		})
		return
	}

	utils.ResponseWithJSON(pc.Logger, w, http.StatusBadRequest, utils.Map{
		"ok": true,
	})
}

func (pc *PaymentsController) GetProductPlansByAppCode(w http.ResponseWriter, r *http.Request) {
	appCode := r.URL.Query().Get("appCode")
	if appCode == "" {
		log.Printf("Get Product Plans by App Code error controller.payments.GetProductPlansByAppCode: %v", "App code not found")
		utils.ErrorResponse(pc.Logger, w, http.StatusBadRequest, []byte("App code is empty"))
		return
	}

	plans, err := pc.ProductPlanStore.GetByAppCode(appCode)
	if err != nil {
		log.Printf("Get Product Plans by App Code error controller.payments.GetProductPlansByAppCode: %v", err)
		// return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
		utils.ErrorResponse(pc.Logger, w, http.StatusInternalServerError, []byte(err.Error()))
		return
	}

	// return c.JSON(plans)
	utils.ResponseWithJSON(pc.Logger, w, http.StatusOK, utils.Map{
		"success": true,
		"plans":   plans,
	})
}
