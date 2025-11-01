package dto

import (
	"encoding/json"
	"time"

	"github.com/21TechLabs/factory-backend/utils"
)

type ProductPlanCreate struct {
	PlanName         string                `json:"planName" validate:"required"`
	PlanDescription  string                `json:"planDescription" validate:"required"`
	PlanPrice        float64               `json:"planPrice" validate:"required,gt=0"`
	PlanCurrency     utils.Currency        `json:"planCurrency" validate:"required,oneof=USD EUR INR"`
	PlanDuration     time.Duration         `json:"planDuration" validate:"required"`
	PlanType         utils.PlanType        `json:"planType" validate:"required,oneof=subscription one_time"`
	Tokens           int64                 `json:"tokens" validate:"required,gte=0"`
	IsActive         bool                  `json:"isActive" validate:"required"`
	Features         utils.StringSlice     `json:"features" validate:"required"`
	PaymentGatewayID utils.JSONMap[string] `json:"paymentGatewayId" validate:"required"`
}

type ProductPlanFetchDto struct {
	PlanType     utils.PlanType `json:"planType" validate:"omitempty,oneof=subscription one_time"`
	IsActive     *bool          `json:"isActive" validate:"omitempty"`
	PlanCurrency utils.Currency `json:"planCurrency" validate:"omitempty,oneof=USD EUR INR"`
	Start        json.Number    `json:"start" validate:"omitempty,gte=0"`
	Limit        json.Number    `json:"limit" validate:"omitempty,gte=1,lte=100"`
	SortBy       []utils.SortBy `json:"sortBy" validate:"omitempty,dive,required"`
}
