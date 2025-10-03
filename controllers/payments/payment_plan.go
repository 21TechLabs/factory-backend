package payments_controller

import (
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
}

func NewPaymentPlanController(logger *log.Logger, store *models.ProductPlanStore) *PaymentPlanController {
	return &PaymentPlanController{ProductPlanStore: store, Logger: logger}
}

func (ppc *PaymentPlanController) CreatePaymentPlan(w http.ResponseWriter, r *http.Request) {
	plan, err := utils.ReadContextValue[*dto.PaymentPlanCreate](r, utils.SchemaValidatorContextKey)
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

	err = ppc.ProductPlanStore.CreatePaymentPlan(plan, user)
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

	plan, err := utils.ReadContextValue[dto.PaymentPlanCreate](r, utils.SchemaValidatorContextKey)
	if err != nil {
		utils.ErrorResponse(ppc.Logger, w, http.StatusBadRequest, []byte("Invalid payment plan data"))
		return
	}

	user, err := utils.ReadContextValue[*models.User](r, utils.UserContextKey)
	if err != nil {
		utils.ErrorResponse(ppc.Logger, w, http.StatusUnauthorized, []byte("User not authenticated"))
		return
	}

	err = ppc.ProductPlanStore.UpdatePaymentPlan(uint(planID), &plan, user)
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

	paymentPlan, err := ppc.ProductPlanStore.GetPaymentPlanByID(uint(planID))
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
