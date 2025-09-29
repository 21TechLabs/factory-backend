package payments

import (
	"github.com/21TechLabs/factory-backend/utils"
	"gorm.io/gorm"
)

type PaymentPlans struct {
	gorm.Model
	PlanName         string            `gorm:"column:plan_name" json:"planName"`
	PlanDescription  string            `gorm:"column:plan_description" json:"planDescription"`
	PlanPrice        int64             `gorm:"column:plan_price" json:"planPrice"`
	PlanDuration     int64             `gorm:"column:plan_duration" json:"planDuration"`
	IsActive         bool              `gorm:"column:is_active" json:"isActive"`
	Features         utils.StringSlice `gorm:"type:json column:features" json:"features"`
	CreatedBy        uint              `gorm:"column:created_by" json:"createdBy"`
	UpdatedBy        uint              `gorm:"column:updated_by" json:"updatedBy"`
	PaymentGatewayID utils.JSONMap     `gorm:"type:json column:payment_gateway_id" json:"paymentGatewayId"`
}
