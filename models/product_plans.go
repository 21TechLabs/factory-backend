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
	UserStore        *UserStore
}

func NewProductPlanStore(db *gorm.DB, us *UserStore) *ProductPlanStore {
	return &ProductPlanStore{DB: db, TransactionStore: NewTransactionStore(db), UserStore: us}
}

type ProductPlan struct {
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
	UpdatedByUser    User              `gorm:"foreignKey:UpdatedBy;references:ID" json:"-"`
	PaymentGatewayID utils.JSONMap     `gorm:"type:json;column:payment_gateway_id" json:"paymentGatewayId"`
	CreatedAt        time.Time         `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt        time.Time         `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (pps *ProductPlanStore) CreateProductPlan(plan *dto.ProductPlanCreate, user *User) error {
	ProductPlan := &ProductPlan{
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

	var tx = pps.DB.Create(ProductPlan)

	if tx == nil {
		return errors.New("failed to create payment plan")
	}
	return tx.Error
}

func (pps *ProductPlanStore) UpdateProductPlan(id uint, plan *dto.ProductPlanCreate, user *User) error {
	var ProductPlan ProductPlan
	if err := pps.DB.First(&ProductPlan, id).Error; err != nil {
		return err
	}

	ProductPlan.PlanName = plan.PlanName
	ProductPlan.PlanDescription = plan.PlanDescription
	ProductPlan.PlanPrice = plan.PlanPrice
	ProductPlan.PlanCurrency = plan.PlanCurrency
	ProductPlan.PlanDuration = plan.PlanDuration
	ProductPlan.PlanType = plan.PlanType
	ProductPlan.Tokens = plan.Tokens
	ProductPlan.IsActive = plan.IsActive
	ProductPlan.Features = plan.Features
	ProductPlan.PaymentGatewayID = plan.PaymentGatewayID
	ProductPlan.UpdatedBy = user.ID

	return pps.DB.Save(&ProductPlan).Error
}

func (pps *ProductPlanStore) GetProductPlanByID(id uint) (*ProductPlan, error) {
	var ProductPlan ProductPlan
	if err := pps.DB.First(&ProductPlan, id).Error; err != nil {
		return nil, err
	}
	return &ProductPlan, nil
}

func (pps *ProductPlanStore) GetProductPlans(filter dto.ProductPlanFetchDto) ([]ProductPlan, error) {
	var plans []ProductPlan
	query := pps.DB.Model(&ProductPlan{})

	if filter.PlanType != "" {
		query = query.Where("plan_type = ?", filter.PlanType)
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	if filter.PlanCurrency != "" {
		query = query.Where("plan_currency = ?", filter.PlanCurrency)
	}

	if filter.SortBy != nil {
		for _, sort := range filter.SortBy {
			switch sort.Field {
			case "planName":
				query = query.Order("plan_name " + sort.Direction)
			case "planPrice":
				query = query.Order("plan_price " + sort.Direction)
			case "createdAt":
				query = query.Order("created_at " + sort.Direction)
			case "updatedAt":
				query = query.Order("updated_at " + sort.Direction)
			}
		}
	}

	if filter.Start != "" {
		start, err := filter.Start.Int64()
		if err != nil {
			return nil, err
		}
		query = query.Offset(int(start))
	}
	if filter.Limit != "" {
		limit, err := filter.Limit.Int64()
		if err != nil {
			return nil, err
		}
		if limit <= 0 {
			return nil, utils.ErrInvalidLimit
		}
		query = query.Limit(int(limit))
	}

	if err := query.Find(&plans).Error; err != nil {
		return nil, err
	}

	return plans, nil
}

func (pps *ProductPlanStore) DeleteProductPlan(id uint) error {
	var ProductPlan ProductPlan
	if err := pps.DB.First(&ProductPlan, id).Error; err != nil {
		return err
	}
	if err := pps.DB.Delete(&ProductPlan).Error; err != nil {
		return err
	}
	return nil
}
