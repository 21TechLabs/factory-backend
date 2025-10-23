package models

import (
	"time"

	"github.com/21TechLabs/factory-backend/dto"
	"github.com/21TechLabs/factory-backend/utils"
	"gorm.io/gorm"
)

type UserSubscriptionStore struct {
	db        *gorm.DB
	UserStore *UserStore
}

func NewUserSubscriptionStore(db *gorm.DB, userStore *UserStore) *UserSubscriptionStore {
	return &UserSubscriptionStore{
		db:        db,
		UserStore: userStore,
	}
}

type UserSubscription struct {
	ID                 uint                     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID             uint                     `gorm:"column:user_id" json:"userId"`
	User               User                     `gorm:"foreignKey:UserID;references:ID" json:"-"`
	StartDate          time.Time                `gorm:"column:start_date" json:"startDate"`
	EndDate            time.Time                `gorm:"column:end_date" json:"endDate"`
	IsActive           bool                     `gorm:"column:is_active" json:"isActive"`
	SubscriptionStatus utils.SubscriptionStatus `gorm:"column:status" json:"status"`
	PaymentGatewayName string                   `gorm:"column:payment_gateway_name" json:"paymentGatewayName"`
	SubscriptionID     string                   `gorm:"column:subscription_id;unique" json:"subscriptionId"`
	ProductPlanID      uint                     `gorm:"column:product_plan_id" json:"-"`
	ProductPlan        ProductPlan              `gorm:"foreignKey:ProductPlanID;references:ID" json:"-"`
	ChargedCount       int                      `gorm:"column:charged_count" json:"chargedCount"`
	TotalChargedCount  int                      `gorm:"column:total_charged_count" json:"totalChargedCount"`
	Suspended          bool                     `gorm:"column:suspended" json:"suspended"`
	CreatedAt          time.Time                `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt          time.Time                `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (UserSubscription) TableName() string {
	return "user_subscriptions"
}

func (uss *UserSubscriptionStore) Create(us UserSubscription) (*UserSubscription, error) {
	err := uss.db.Create(&us).Error
	if err != nil {
		return nil, err
	}
	return &us, nil
}

func (uss *UserSubscriptionStore) FindBy(filter dto.UserSubscriptionFilterDto, start, limit int) ([]UserSubscription, error) {
	var subscriptions []UserSubscription

	query := uss.db.Model(&UserSubscription{})

	if len(filter.UserIds) > 0 {
		query = query.Where("user_id IN (?)", filter.UserIds)
	}
	if filter.StartDate.Min != nil {
		query = query.Where("start_date >= ?", filter.StartDate.Min)
	}
	if filter.StartDate.Max != nil {
		query = query.Where("start_date <= ?", filter.StartDate.Max)
	}

	if filter.EndDate.Min != nil {
		query = query.Where("end_date >= ?", filter.EndDate.Min)
	}
	if filter.EndDate.Max != nil {
		query = query.Where("end_date <= ?", filter.EndDate.Max)
	}

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", filter.IsActive)
	}

	if len(filter.SubscriptionStatus) > 0 {
		query = query.Where("subscription_status = ?", filter.SubscriptionStatus)
	}

	if filter.ChargedCount.Min != nil {
		query = query.Where("charged_count >= ?", filter.ChargedCount.Min)
	}
	if filter.ChargedCount.Max != nil {
		query = query.Where("charged_count <= ?", filter.ChargedCount.Max)
	}

	if filter.TotalChargedCount.Min != nil {
		query = query.Where("total_charged_count >= ?", filter.TotalChargedCount.Min)
	}
	if filter.TotalChargedCount.Max != nil {
		query = query.Where("total_charged_count <= ?", filter.TotalChargedCount.Max)
	}

	if filter.CreatedAt.Min != nil {
		query = query.Where("created_at >= ?", filter.CreatedAt.Min)
	}

	if filter.CreatedAt.Max != nil {
		query = query.Where("created_at <= ?", filter.CreatedAt.Max)
	}

	if filter.UpdatedAt.Min != nil {
		query = query.Where("updated_at >= ?", filter.UpdatedAt.Min)
	}
	if filter.UpdatedAt.Max != nil {
		query = query.Where("updated_at <= ?", filter.UpdatedAt.Max)
	}

	if filter.PreloadUser {
		query = query.Preload("User")
	}

	if filter.PreloadProductPlan {
		query = query.Preload("ProductPlan")
	}

	tx := query.Offset(start).Limit(limit).Find(&subscriptions)

	if tx == nil {
		return subscriptions, gorm.ErrRecordNotFound
	}
	return subscriptions, tx.Error
}

func (uss *UserSubscriptionStore) FindOne(filter dto.UserSubscriptionFilterDto) (*UserSubscription, error) {
	filter.PreloadUser = true
	filter.PreloadProductPlan = true
	sub, err := uss.FindBy(filter, 0, 1)
	if err != nil {
		return nil, err
	}
	if sub == nil || len(sub) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &sub[0], nil
}

func (uss *UserSubscriptionStore) FindBySubscriptionID(subscriptionID string) (*UserSubscription, error) {
	var user UserSubscription
	err := uss.db.Where(&UserSubscription{SubscriptionID: subscriptionID}).First(&user).Error
	return &user, err
}

func (uss *UserSubscriptionStore) Save(us *UserSubscription) error {
	tx := uss.db.Save(&us)
	if tx == nil {
		return gorm.ErrRecordNotFound
	}
	return tx.Error
}
