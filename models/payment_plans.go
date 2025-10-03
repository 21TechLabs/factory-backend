package models

import (
	"errors"
	"time"

	"github.com/21TechLabs/factory-backend/dto"
	"github.com/21TechLabs/factory-backend/utils"
	"gorm.io/gorm"
)

type ProductPlanStore struct {
	DB               *gorm.DB
	TransactionStore *TransactionStore
}

func NewProductPlanStore(db *gorm.DB) *ProductPlanStore {
	return &ProductPlanStore{DB: db, TransactionStore: NewTransactionStore(db)}
}

type PaymentPlan struct {
	gorm.Model
	ID               uint              `gorm:"primaryKey;autoIncrement" json:"id"`
	PlanName         string            `gorm:"column:plan_name" json:"planName"`
	PlanDescription  string            `gorm:"column:plan_description" json:"planDescription"`
	PlanPrice        float64           `gorm:"column:plan_price" json:"planPrice"`
	PlanCurrency     utils.Currency    `gorm:"column:plan_currency" json:"planCurrency"`
	PlanDuration     time.Duration     `gorm:"column:plan_duration" json:"planDuration"`
	PlanType         utils.PlanType    `gorm:"column:plan_type" json:"planType"`
	Tokens           int64             `gorm:"column:tokens" json:"tokens"`
	IsActive         bool              `gorm:"column:is_active" json:"isActive"`
	Features         utils.StringSlice `gorm:"type:json;column:features" json:"features"`
	UpdatedBy        uint              `gorm:"column:updated_by" json:"updatedBy"`
	UpdatedByUser    User              `gorm:"foreignKey:UpdatedBy;references:ID" json:"user"`
	PaymentGatewayID utils.JSONMap     `gorm:"type:json column:payment_gateway_id" json:"paymentGatewayId"`
}

func (pps *ProductPlanStore) CreatePaymentPlan(plan *dto.PaymentPlanCreate, user *User) error {
	paymentPlan := &PaymentPlan{
		PlanName:         plan.PlanName,
		PlanDescription:  plan.PlanDescription,
		PlanPrice:        plan.PlanPrice,
		PlanCurrency:     plan.PlanCurrency,
		PlanDuration:     plan.PlanDuration,
		PlanType:         plan.PlanType,
		Tokens:           plan.Tokens,
		IsActive:         plan.IsActive,
		Features:         plan.Features,
		PaymentGatewayID: plan.PaymentGatewayID,
		UpdatedBy:        user.ID,
	}

	var tx = pps.DB.Create(paymentPlan)

	if tx == nil {
		return errors.New("failed to create payment plan")
	}
	return tx.Error
}

func (pps *ProductPlanStore) UpdatePaymentPlan(id uint, plan *dto.PaymentPlanCreate, user *User) error {
	var paymentPlan PaymentPlan
	if err := pps.DB.First(&paymentPlan, id).Error; err != nil {
		return err
	}

	paymentPlan.PlanName = plan.PlanName
	paymentPlan.PlanDescription = plan.PlanDescription
	paymentPlan.PlanPrice = plan.PlanPrice
	paymentPlan.PlanCurrency = plan.PlanCurrency
	paymentPlan.PlanDuration = plan.PlanDuration
	paymentPlan.PlanType = plan.PlanType
	paymentPlan.Tokens = plan.Tokens
	paymentPlan.IsActive = plan.IsActive
	paymentPlan.Features = plan.Features
	paymentPlan.PaymentGatewayID = plan.PaymentGatewayID
	paymentPlan.UpdatedBy = user.ID

	return pps.DB.Save(&paymentPlan).Error
}

func (pps *ProductPlanStore) GetPaymentPlanByID(id uint) (*PaymentPlan, error) {
	var paymentPlan PaymentPlan
	if err := pps.DB.First(&paymentPlan, id).Error; err != nil {
		return nil, err
	}
	return &paymentPlan, nil
}
