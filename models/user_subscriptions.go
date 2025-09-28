package models

import (
	"time"

	"github.com/21TechLabs/factory-backend/utils"
	"gorm.io/gorm"
)

type UserSubscriptionStore struct {
	DB *gorm.DB
}

func NewUserSubscriptionStore(db *gorm.DB) *UserSubscriptionStore {
	return &UserSubscriptionStore{
		DB: db,
	}
}

type UserSubscription struct {
	gorm.Model
	UserID                   uint               `json:"user_id" gorm:"column:user_id;not null;"`
	User                     User               `json:"-" gorm:"foreignKey:UserID"`
	Plan                     string             `json:"plan" gorm:"plan"`
	SubscriptionStartsAt     time.Time          `json:"subscriptionStartsAt" gorm:"subscription_starts_at"`
	SubscriptionEndsAt       time.Time          `json:"subscriptionEndsAt" gorm:"subscription_ends_at"`
	ReSubscribeOn            time.Time          `json:"reSubscribeOn" gorm:"re_subscribe_on"`
	SubscriptionTenureInDays int                `json:"subscriptionTenureInDays" gorm:"subscription_tenure_in_days"`
	Status                   SubscriptionStatus `json:"status" gorm:"status"`
	PlanLevel                int                `json:"planLevel" gorm:"plan_level"`
	TokenReward              int64              `json:"tokenReward" gorm:"token_reward"`
	PaymentId                utils.JSONMap      `json:"paymentId" gorm:"payment_id"`
	AppCode                  string             `json:"appCode" gorm:"app_code"`
	HasExpired               bool               `json:"hasExpired" gorm:"has_expired"`
	PaymentMethod            string             `json:"paymentMethod" gorm:"payment_method"`
	ExtraOptions             utils.JSONMap      `json:"extraOptions" gorm:"extra_options"`
	UsageStats               utils.JSONMap      `json:"usageStats" gorm:"usage_stats"`
	MaxUsageStats            utils.JSONMap      `json:"maxUsageStats" gorm:"max_usage_stats"`
}

func (us *UserSubscription) TableName() string {
	return "user_subscriptions"
}

func (uss *UserSubscriptionStore) Save(us *UserSubscription, update bool) error {
	if update {
		result := uss.DB.Model(us).Updates(us)
		return result.Error
	}

	result := uss.DB.Create(us)
	return result.Error
}

func (uss *UserSubscriptionStore) GetUserSubscriptionById(userId uint) (*UserSubscription, error) {
	var userSubscription UserSubscription
	if result := uss.DB.Where("user_id = ? and status = ?", userId, SubscriptionStatusActive).First(&userSubscription); result.Error != nil {
		return nil, result.Error
	}
	return &userSubscription, nil
}

func (uss *UserSubscriptionStore) GetUserSubscriptionBy(filter map[string]interface{}) (*UserSubscription, error) {
	var userSubscription UserSubscription
	if result := uss.DB.Where(filter).First(&userSubscription); result.Error != nil {
		return nil, result.Error
	}
	return &userSubscription, nil
}
