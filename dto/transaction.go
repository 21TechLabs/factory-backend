package dto

import (
	"github.com/21TechLabs/factory-backend/utils"
	"github.com/google/uuid"
)

type TransactionCreateDto struct {
	Token                       int64                   `json:"token"`
	Amount                      float64                 `json:"amount"`
	Currency                    utils.Currency          `json:"currency"`
	Status                      utils.TransactionStatus `json:"status"`
	ProductPlanID               *uuid.UUID              `json:"productPlanId,omitempty"` // Optional, can be nil
	PaymentGatewayName          string                  `json:"paymentGatewayName"`
	PaymentGatewayRedirectURL   string                  `json:"paymentGatewayRedirectUrl"`
	PaymentGatewayTransactionID string                  `json:"transactionId"`
	UserSubscriptionID          *uuid.UUID              `json:"userSubscriptionId,omitempty"`
}
