package dto

import "github.com/21TechLabs/factory-backend/utils"

//	transaction := &Transaction{
//			UserID:             rpg.User.ID,
//			Token:              rpg.ProductPlan.Tokens,
//			Amount:             rpg.ProductPlan.PlanPrice,
//			Currency:           rpg.ProductPlan.PlanCurrency,
//			Status:             utils.TransactionStatusPending,
//			ProductPlanID:      &rpg.ProductPlan.ID,
//			PaymentGatewayName: PaymentGatewayRazorpay,
//		}
type TransactionCreateDto struct {
	Token                       int64                   `json:"token"`
	Amount                      float64                 `json:"amount"`
	Currency                    utils.Currency          `json:"currency"`
	Status                      utils.TransactionStatus `json:"status"`
	ProductPlanID               *uint                   `json:"productPlanId,omitempty"` // Optional, can be nil
	PaymentGatewayName          string                  `json:"paymentGatewayName"`
	PaymentGatewayRedirectURL   string                  `json:"paymentGatewayRedirectUrl"`
	PaymentGatewayTransactionID string                  `json:"transactionId"`
}
