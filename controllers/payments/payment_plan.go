package payments_controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/21TechLabs/factory-backend/dto"
	"github.com/21TechLabs/factory-backend/models"
	"github.com/21TechLabs/factory-backend/utils"
)

type PaymentPlanController struct {
	Logger           *log.Logger
	ProductPlanStore *models.ProductPlanStore
	TransactionStore *models.TransactionStore
	UserStore        *models.UserStore
}

var RazorpayHMECSecret = ""

func init() {
	RazorpayHMECSecret = utils.GetEnv("PAYMENTS_HMEC_SECRET", true)
}

// NewPaymentPlanController creates a PaymentPlanController configured with the provided logger and product plan store.
// The logger is used for request-related logging; store provides persistence operations for payment plans.
func NewPaymentPlanController(logger *log.Logger, store *models.ProductPlanStore) *PaymentPlanController {
	fs := models.NewFileStore(store.DB)
	us := models.NewUserStore(store.DB, fs)
	return &PaymentPlanController{ProductPlanStore: store, Logger: logger, TransactionStore: models.NewTransactionStore(store.DB), UserStore: us}
}

func (ppc *PaymentPlanController) CreateProductPlan(w http.ResponseWriter, r *http.Request) {
	plan, err := utils.ReadContextValue[*dto.ProductPlanCreate](r, utils.SchemaValidatorContextKey)
	if err != nil {
		utils.ErrorResponse(ppc.Logger, w, http.StatusBadRequest, []byte("Invalid payment plan data"))
		return
	}

	user, err := utils.ReadContextValue[*models.User](r, utils.UserContextKey)
	if err != nil {
		fmt.Println("User not authenticated:", err)
		utils.ErrorResponse(ppc.Logger, w, http.StatusUnauthorized, []byte("User not authenticated"))
		return
	}

	err = ppc.ProductPlanStore.CreateProductPlan(plan, user)
	if err != nil {
		utils.ErrorResponse(ppc.Logger, w, http.StatusInternalServerError, []byte(err.Error()))
		return
	}

	utils.ResponseWithJSON(ppc.Logger, w, http.StatusCreated, utils.Map{
		"message": "Payment plan created successfully",
		"plan":    plan,
	})
}

func (ppc *PaymentPlanController) UpdatePaymentPlan(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		utils.ErrorResponse(ppc.Logger, w, http.StatusBadRequest, []byte("Payment plan ID is required"))
		return
	}

	// convert id to uint
	planID, err := strconv.Atoi(id)
	if err != nil {
		utils.ErrorResponse(ppc.Logger, w, http.StatusBadRequest, []byte("Invalid payment plan ID"))
		return
	}

	plan, err := utils.ReadContextValue[*dto.ProductPlanCreate](r, utils.SchemaValidatorContextKey)
	if err != nil {
		utils.ErrorResponse(ppc.Logger, w, http.StatusBadRequest, []byte("Invalid payment plan data"))
		return
	}

	user, err := utils.ReadContextValue[*models.User](r, utils.UserContextKey)
	if err != nil {
		utils.ErrorResponse(ppc.Logger, w, http.StatusUnauthorized, []byte("User not authenticated"))
		return
	}

	err = ppc.ProductPlanStore.UpdateProductPlan(uint(planID), plan, user)
	if err != nil {
		utils.ErrorResponse(ppc.Logger, w, http.StatusInternalServerError, []byte(err.Error()))
		return
	}

	utils.ResponseWithJSON(ppc.Logger, w, http.StatusOK, utils.Map{
		"message": "Payment plan updated successfully",
		"plan":    plan,
	})
}

func (ppc *PaymentPlanController) GetPaymentPlanByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Fetching payment plan by ID")
	id := r.PathValue("id")
	if id == "" {
		utils.ErrorResponse(ppc.Logger, w, http.StatusBadRequest, []byte("Payment plan ID is required"))
		return
	}

	// convert id to uint
	planID, err := strconv.Atoi(id)
	if err != nil {
		utils.ErrorResponse(ppc.Logger, w, http.StatusBadRequest, []byte("Invalid payment plan ID"))
		return
	}

	paymentPlan, err := ppc.ProductPlanStore.GetProductPlanByID(uint(planID))
	if err != nil {
		utils.ErrorResponse(ppc.Logger, w, http.StatusInternalServerError, []byte(err.Error()))
		return
	}

	if paymentPlan == nil {
		utils.ErrorResponse(ppc.Logger, w, http.StatusNotFound, []byte("Payment plan not found"))
		return
	}

	utils.ResponseWithJSON(ppc.Logger, w, http.StatusOK, utils.Map{
		"plan": paymentPlan,
	})
}

func (ppc *PaymentPlanController) GetProductPlans(w http.ResponseWriter, r *http.Request) {
	// load query parameters as body
	body := &dto.ProductPlanFetchDto{}
	if err := utils.ParseQueryParams(r, body); err != nil {
		fmt.Println("Error parsing query parameters:", err)
		utils.ErrorResponse(ppc.Logger, w, http.StatusBadRequest, []byte("Invalid query parameters"))
		return
	}

	// validate body
	if err := utils.ValidateStruct(body); err != nil {
		utils.ErrorResponse(ppc.Logger, w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	plans, err := ppc.ProductPlanStore.GetProductPlans(*body)
	if err != nil {
		utils.ErrorResponse(ppc.Logger, w, http.StatusInternalServerError, []byte(err.Error()))
		return
	}

	utils.ResponseWithJSON(ppc.Logger, w, http.StatusOK, utils.Map{
		"plans": plans,
	})
}

func (ppc *PaymentPlanController) DeletePaymentPlan(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		utils.ErrorResponse(ppc.Logger, w, http.StatusBadRequest, []byte("Payment plan ID is required"))
		return
	}

	planID, err := strconv.Atoi(id)
	if err != nil {
		utils.ErrorResponse(ppc.Logger, w, http.StatusBadRequest, []byte("Invalid payment plan ID"))
		return
	}

	_, err = utils.ReadContextValue[*models.User](r, utils.UserContextKey)
	if err != nil {
		utils.ErrorResponse(ppc.Logger, w, http.StatusUnauthorized, []byte("User not authenticated"))
		return
	}

	err = ppc.ProductPlanStore.DeleteProductPlan(uint(planID))
	if err != nil {
		utils.ErrorResponse(ppc.Logger, w, http.StatusInternalServerError, []byte(err.Error()))
		return
	}

	utils.ResponseWithJSON(ppc.Logger, w, http.StatusOK, utils.Map{
		"message": "Payment plan deleted successfully",
	})
}

func (ppc *PaymentPlanController) ProcessOrderPaid(w http.ResponseWriter, r *http.Request) {
	paymentGateway := r.PathValue("paymentGateway")
	// always respond with 200 OK (but log the error if any)

	gateway, err := models.GetPaymentGateway(paymentGateway, ppc.Logger, ppc.TransactionStore, ppc.UserStore)

	if err != nil {
		ppc.Logger.Printf("Error getting payment gateway: %v", err)
		utils.ResponseWithJSON(ppc.Logger, w, http.StatusOK, utils.Map{
			"message": "Payment gateway not found",
		})
		return
	}

	var rawBody []byte = make([]byte, r.ContentLength)

	for {
		n, err := r.Body.Read(rawBody)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			ppc.Logger.Printf("Error reading request body: %v", err)
			utils.ResponseWithJSON(ppc.Logger, w, http.StatusInternalServerError, utils.Map{
				"message": "Failed to read request body",
			})
			return
		}
		if n == 0 {
			break
		}
	}
	defer r.Body.Close()

	isValid := utils.ValidateHeaderHMACSha256(rawBody, RazorpayHMECSecret, r.Header.Get("X-Razorpay-Signature"))
	if !isValid {
		utils.ResponseWithJSON(ppc.Logger, w, http.StatusBadRequest, utils.Map{
			"message": "Invalid HMAC signature",
		})
		return
	}

	// body
	var body models.RazorpayBaseEvent[models.RazorpayOrderPaidPayload]

	if err := json.Unmarshal(rawBody, &body); err != nil {
		ppc.Logger.Printf("Error unmarshalling request body: %v", err)
		utils.ResponseWithJSON(ppc.Logger, w, http.StatusBadRequest, utils.Map{
			"message": "Invalid request body",
		})
		return
	}

	// validate body
	if err := utils.ValidateStruct(body); err != nil {
		ppc.Logger.Printf("Validation error: %v", err)
		utils.ResponseWithJSON(ppc.Logger, w, http.StatusBadRequest, utils.Map{
			"message": err.Error(),
		})
		return
	}

	ppc.Logger.Printf("Processing order paid event for order ID: %s", body.Payload.Order.Entity.ID)

	txn, err := gateway.CaptureOrderPaid(body)

	if err != nil {
		ppc.Logger.Printf("Error capturing order paid: %v", err)
		utils.ResponseWithJSON(ppc.Logger, w, http.StatusInternalServerError, utils.Map{
			"message": "Failed to capture order paid",
		})
		return
	}

	ppc.Logger.Printf("Order %s captured successfully for transaction ID %d", body.Payload.Order.Entity.ID, txn.ID)

	utils.ResponseWithJSON(ppc.Logger, w, http.StatusOK, utils.Map{
		"success": "Ok!",
	})
}

func (ppc *PaymentPlanController) ProductBuy(w http.ResponseWriter, r *http.Request) {
	user, err := utils.ReadContextValue[*models.User](r, utils.UserContextKey)
	if err != nil {
		utils.ErrorResponse(ppc.Logger, w, http.StatusUnauthorized, []byte("User not authenticated"))
		return
	}

	_productId := r.PathValue("productId")
	if _productId == "" {
		utils.ErrorResponse(ppc.Logger, w, http.StatusBadRequest, []byte("Product ID is required"))
		return
	}

	productId, err := strconv.Atoi(_productId)
	if err != nil {
		utils.ErrorResponse(ppc.Logger, w, http.StatusBadRequest, []byte("Invalid product ID"))
		return
	}

	_count := r.URL.Query().Get("count")
	count := 1
	if _count != "" {
		count, err = strconv.Atoi(_count)
		if err != nil {
			utils.ErrorResponse(ppc.Logger, w, http.StatusBadRequest, []byte("Invalid count parameter"))
			return
		}
	}

	paymentGateway := r.PathValue("paymentGateway")
	if paymentGateway == "" {
		utils.ErrorResponse(ppc.Logger, w, http.StatusBadRequest, []byte("Payment gateway is required"))
		return
	}

	gateway, err := models.GetPaymentGateway(paymentGateway, ppc.Logger, ppc.TransactionStore, ppc.UserStore)
	if err != nil {
		utils.ErrorResponse(ppc.Logger, w, http.StatusInternalServerError, []byte(err.Error()))
		return
	}

	// get product plan by ID
	plan, err := ppc.ProductPlanStore.GetProductPlanByID(uint(productId))
	if err != nil {
		utils.ErrorResponse(ppc.Logger, w, http.StatusInternalServerError, []byte(err.Error()))
		return
	}

	txn, err := gateway.InitiatePayment(plan, user, count)

	if err != nil {
		utils.ErrorResponse(ppc.Logger, w, http.StatusInternalServerError, []byte(err.Error()))
		return
	}

	ppc.Logger.Printf("Payment initiated successfully for user %s with transaction ID %d", user.ID, txn.ID)

	utils.ResponseWithJSON(ppc.Logger, w, http.StatusOK, utils.Map{
		"message": "Payment initiated successfully",
		"txn":     txn,
	})
}
