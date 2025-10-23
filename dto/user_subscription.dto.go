package dto

import (
	"time"

	"github.com/21TechLabs/factory-backend/utils"
)

type UserSubscriptionFilterDto struct {
	UserIds            []uint                     `json:"user_ids,omitempty"`
	StartDate          MinMax[time.Time]          `json:"start_date,omitempty"`
	EndDate            MinMax[time.Time]          `json:"end_date,omitempty"`
	IsActive           *bool                      `json:"is_active,omitempty"`
	SubscriptionStatus []utils.SubscriptionStatus `json:"subscription_status,omitempty"`
	ChargedCount       MinMax[uint]               `json:"charged_count,omitempty"`
	TotalChargedCount  MinMax[uint]               `json:"total_charged_count,omitempty"`
	CreatedAt          MinMax[time.Time]          `json:"created_at,omitempty"`
	UpdatedAt          MinMax[time.Time]          `json:"updated_at,omitempty"`
	PreloadUser        bool
	PreloadProductPlan bool
}
