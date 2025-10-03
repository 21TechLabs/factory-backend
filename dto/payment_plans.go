package dto

import (
	"time"

	"github.com/21TechLabs/factory-backend/utils"
)

type PaymentPlanCreate struct {
	PlanName         string            `json:"planName" validate:"required"`
	PlanDescription  string            `json:"planDescription" validate:"required"`
	PlanPrice        float64           `json:"planPrice" validate:"required,gt=0"`
	PlanCurrency     utils.Currency    `json:"planCurrency" validate:"required,oneof='USD' 'EUR' 'INR'"`
	PlanDuration     time.Duration     `json:"planDuration" validate:"required"`
	PlanType         utils.PlanType    `json:"planType" validate:"required,oneof='subscription' 'one_time'"`
	Tokens           int64             `json:"tokens" validate:"required,gte=0"`
	IsActive         bool              `json:"isActive" validate:"required"`
	Features         utils.StringSlice `json:"features" validate:"required"`
	PaymentGatewayID utils.JSONMap     `json:"paymentGatewayId" validate:"required"`
}
