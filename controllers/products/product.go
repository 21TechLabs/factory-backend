package products_controller

import (
	"log"
	"net/http"

	"github.com/21TechLabs/factory-backend/models"
	"github.com/21TechLabs/factory-backend/utils"
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

func (ppc *ProductPlanController) GetProductByAppCode(w http.ResponseWriter, r *http.Request) {
	// var appCode string = c.Params("appCode")
	appCode := r.PathValue("appCode")

	if len(appCode) == 0 {
		// return utils.ErrorResponse(c, fiber.StatusBadRequest, "App code is required!")
		log.Printf("App code is required!")
		utils.ErrorResponse(ppc.Logger, w, http.StatusBadRequest, []byte("App code is required!"))
		return
	}

	product, err := ppc.ProductPlanStore.GetByAppCode(appCode)

	if err != nil {
		log.Printf("Failed to fetch product with app code \"%v\" because an error occured: %v", appCode, err.Error())
		// return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to fetch product!")
		utils.ErrorResponse(ppc.Logger, w, http.StatusBadRequest, []byte("Failed to fetch product!"))
		return
	}

	utils.ResponseWithJSON(ppc.Logger, w, http.StatusOK, utils.Map{
		"product": product,
		"success": true,
	})
}

func (ppc *ProductPlanController) GetUsersActiveProductSubsctiptionByAppCode(w http.ResponseWriter, r *http.Request) {
	appCode := r.PathValue("appCode")

	if len(appCode) == 0 {
		log.Printf("App code is required!")
		utils.ErrorResponse(ppc.Logger, w, http.StatusBadRequest, []byte("App code is required!"))
		return
	}

	user, ok := r.Context().Value(utils.UserContextKey).(*models.User)

	if !ok {
		log.Printf("Failed to fetch user because user is not found in the context!")
		utils.ErrorResponse(ppc.Logger, w, http.StatusBadRequest, []byte("Failed to fetch user!"))
		return
	}

	subscription, err := ppc.UserStore.GetActiveAppSubscriptionByAppCode(user, appCode)

	if err != nil {
		log.Printf("Failed to fetch product with app code \"%v\" because an error occured: %v", appCode, err.Error())
		utils.ErrorResponse(ppc.Logger, w, http.StatusBadRequest, []byte("Failed to fetch product!"))
		return
	}

	utils.ResponseWithJSON(ppc.Logger, w, http.StatusOK, utils.Map{
		"subscription": subscription,
		"success":      true,
	})
}
