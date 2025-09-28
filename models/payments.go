package models

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/21TechLabs/factory-backend/utils"
	"gorm.io/gorm"
)

type ProductPlanStore struct {
	DB *gorm.DB
}

func NewProductPlanStore(db *gorm.DB) *ProductPlanStore {
	return &ProductPlanStore{
		DB: db,
	}
}

type SubscriptionStatus string

const (
	SubscriptionStatusPending   SubscriptionStatus = "pending"
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusHalted    SubscriptionStatus = "halted"
	SubscriptionStatusCompleted SubscriptionStatus = "completed"
	SubscriptionStatusCharged   SubscriptionStatus = "charged"
	SubscriptionStatusCancelled SubscriptionStatus = "cancelled"
)

type gateways struct {
	Razorpay string
	Stripe   string
}

type productPlanType struct {
	Subscription string
	OneTime      string
}

var ProductPlanTypes = productPlanType{
	Subscription: "subscription",
	OneTime:      "one_time",
}

var PaymentGatewaysList = gateways{
	Razorpay: "razorpay",
	Stripe:   "stripe",
}

type PaymentGateway interface {
	CreatePayment(*UserSubscriptionStore, ProductPlan, int) (UserSubscription, error)
	VerifyPayment(string) (UserSubscription, error)
	UpdatePaymentStatus(*UserSubscriptionStore, string) (UserSubscription, error)
	CancelSubscription(string) (UserSubscription, error)
	GetUserSubscription() (interface{}, error)
	VerifyWebhookSignature(*http.Request) error
	GetOrderIdFromWebhookRequest([]byte) (string, error)
	SetUserId(uint)
	SetUserViaOrderId(*UserSubscriptionStore, string) error
}

type ProductPlan struct {
	gorm.Model
	ID            uint           `json:"id" gorm:"primaryKey"`
	AppCode       string         `json:"appCode" gorm:"column:app_code;not null;unique"`
	Name          string         `json:"name" gorm:"column:name;not null"`
	Description   string         `json:"description" gorm:"column:description;not null"`
	Type          string         `json:"type" gorm:"column:type;not null"`
	Amount        float64        `json:"amount" gorm:"column:amount;not null"`
	TokenReward   float64        `json:"tokenAmount" gorm:"column:token_reward;not null"`
	Subscriptions []Subscription `json:"subscriptions" gorm:"foreignKey:ProductPlanID"` // Has Many relationship
}

type Subscription struct {
	gorm.Model
	ID                      uint              `json:"id" gorm:"primaryKey"`
	ProductPlanID           uint              `json:"productPlanId" gorm:"column:product_plan_id;not null;"` // Foreign key for the ProductPlan
	ProductPlan             ProductPlan       `gorm:"foreignKey:ProductPlanID"`
	Name                    string            `json:"name"`
	Description             string            `json:"description"`
	PlanID                  utils.JSONMap     `json:"planId" gorm:"type:json"`
	Amount                  float64           `json:"amount"`
	Features                utils.StringSlice `json:"features" gorm:"type:json"`
	EveryNDays              int               `json:"everyNDays"`
	BillingCycles           int               `json:"billingCycle"`
	Level                   int               `json:"level"`
	TokenRewardEveryRenewal int64             `json:"tokenReward"`
	ExpiresAfterDayCount    int               `json:"expiresAfterDayCount"`
	CurrencySymbol          string            `json:"currencySymbol"`
	ForeignCurrencyPricing  utils.JSONMap     `json:"foreignCurrencyPricing" gorm:"type:json"`
	Recommended             bool              `json:"recommended"`
}

func SubscriptionSuccessStatusArr() []SubscriptionStatus {
	return []SubscriptionStatus{
		SubscriptionStatusActive,
		SubscriptionStatusCompleted,
		SubscriptionStatusCharged,
	}
}

func (pps *ProductPlanStore) ProductPlanGetBy(filter *ProductPlan) (ProductPlan, error) {
	var productPlan ProductPlan = ProductPlan{}

	coll := pps.DB.Model(&ProductPlan{})
	if coll == nil {
		return ProductPlan{}, errors.New("database connection not initialized")
	}

	query := coll.Where(filter)

	query = query.Joins("left join subscriptions on subscriptions.product_plan_id = product_plans.id")
	query = query.Preload("Subscriptions").Where("subscriptions.product_plan_id = product_plans.id")

	result := query.First(&productPlan)
	if result.Error != nil {
		return ProductPlan{}, result.Error
	}

	return productPlan, nil
}

func (pps *ProductPlanStore) ProductPlanGetByID(id int) (ProductPlan, error) {
	return pps.ProductPlanGetBy(&ProductPlan{ID: uint(id)})
}

func (pps *ProductPlanStore) GetPaymentGateway(paymentType string, userId uint) (PaymentGateway, error) {
	var payment PaymentGateway
	switch paymentType {
	case PaymentGatewaysList.Razorpay:
		payment = &Razorpay{
			UserId: userId,
		}
	case PaymentGatewaysList.Stripe:
		payment = &Stripe{UserId: userId}
	default:
		return nil, errors.New("invalid payment type")
	}
	return payment, nil
}

func (pps *ProductPlanStore) GetByAppCode(appCode string) (ProductPlan, error) {
	if len(appCode) == 0 {
		return ProductPlan{}, errors.New("invalid app code")
	}
	return pps.ProductPlanGetBy(&ProductPlan{AppCode: appCode})
}

func (pps *ProductPlanStore) SeedProductData() {
	// Seed the payment gateway
	/*Sample data:
		{
	  "_id": {
	    "$oid": "682c3e66b1d6c43ab847537a"
	  },
	  "created_at": {
	    "$date": "2025-05-20T08:33:42.359Z"
	  },
	  "updated_at": {
	    "$date": "2025-05-20T08:33:42.359Z"
	  },
	  "appCode": "app1",
	  "name": "Music LMS",
	  "description": "The Learning Management System for Musisicians",
	  "type": "subscription",
	  "amount": 0,
	  "tokenAmount": 0,
	  "subscriptions": [
	    {
	      "created_at": {
	        "$date": {
	          "$numberLong": "-62135596800000"
	        }
	      },
	      "updated_at": {
	        "$date": {
	          "$numberLong": "-62135596800000"
	        }
	      },
	      "name": "Monthly",
	      "description": "Stay connected and professional with our flexible Monthly Subscriptio",
	      "planId": {
	        "razorpay": "plan_QX7whUHFlZ1NDR",
	        "stripe": "plan_1"
	      },
	      "amount": 50,
	      "features": [
	        "Personalized email addresses",
	        "Advanced spam and phishing protection",
	        "24/7 customer support",
	        "Ample storage for emails and attachments",
	        "Mobile and desktop access",
	        "Easy migration from your current provider",
	        "Email scheduling, templates, and autoresponders"
	      ],
	      "everyNDays": 30,
	      "billingCycle": 36,
	      "level": 1,
	      "tokenReward": {
	        "$numberLong": "0"
	      },
	      "expiresAfterDayCount": 1095,
	      "currencySymbol": "₹",
	      "foreignCurrencyPricing": null,
	      "recommended": false,
	      "modules": [
	        "default"
	      ]
	    },
	    {
	      "created_at": {
	        "$date": {
	          "$numberLong": "-62135596800000"
	        }
	      },
	      "updated_at": {
	        "$date": {
	          "$numberLong": "-62135596800000"
	        }
	      },
	      "name": "Monthly",
	      "description": "Stay connected and professional with our flexible Monthly Subscriptio",
	      "planId": {
	        "razorpay": "plan_QX7xcYv1P0MGmW",
	        "stripe": "plan_1"
	      },
	      "amount": 549,
	      "features": [
	        "Personalized email addresses",
	        "Advanced spam and phishing protection",
	        "24/7 customer support",
	        "Ample storage for emails and attachments",
	        "Mobile and desktop access",
	        "Easy migration from your current provider",
	        "Email scheduling, templates, and autoresponders",
	        "Access across all devices",
	        "Hassle-free migration tools",
	        "Advanced tools: scheduling, templates, autoresponders"
	      ],
	      "everyNDays": 365,
	      "billingCycle": 10,
	      "level": 2,
	      "tokenReward": {
	        "$numberLong": "0"
	      },
	      "expiresAfterDayCount": 1095,
	      "currencySymbol": "₹",
	      "foreignCurrencyPricing": null,
	      "recommended": false,
	      "modules": [
	        "default"
	      ]
	    }
	  ]
	}
	*/
	var product ProductPlan = ProductPlan{
		AppCode:     "app1",
		Name:        "Music LMS",
		Description: "The Learning Management System for Musicians",
		Type:        ProductPlanTypes.Subscription,
		Amount:      0,
		TokenReward: 0,
	}

	// find product by app code
	existingProduct, err := pps.ProductPlanGetBy(&product)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		created := pps.DB.Create(&product)

		if created.Error != nil {
			panic(fmt.Sprintf("Failed to seed product data: %v", created.Error))
		}
		if created.RowsAffected == 0 {
			panic("No rows were affected while seeding product data")
		}
	} else {
		product = existingProduct
	}

	// Seed the subscription
	var subscriptions []Subscription = []Subscription{
		{
			Name:        "Monthly",
			Description: "Stay connected and professional with our flexible Monthly Subscription",
			PlanID: utils.JSONMap{
				"razorpay": "plan_QX7whUHFlZ1NDR",
				"stripe":   "plan_1",
			},
			Amount:                  50,
			Features:                utils.StringSlice{"Personalized email addresses", "Advanced spam and phishing protection", "24/7 customer support", "Ample storage for emails and attachments", "Mobile and desktop access", "Easy migration from your current provider", "Email scheduling, templates, and autoresponders"},
			EveryNDays:              30,
			BillingCycles:           36,
			Level:                   1,
			TokenRewardEveryRenewal: 0,
			ExpiresAfterDayCount:    1095,
			CurrencySymbol:          "₹",
			ForeignCurrencyPricing:  nil,
			Recommended:             false,
		},
		{
			Name:        "Yearly",
			Description: "Stay connected and professional with our flexible Yearly Subscription",
			PlanID: utils.JSONMap{
				"razorpay": "plan_QX7xcYv1P0MGmW",
				"stripe":   "plan_1",
			},
			Amount:                  549,
			Features:                utils.StringSlice{"Personalized email addresses", "Advanced spam and phishing protection", "24/7 customer support", "Ample storage for emails and attachments", "Mobile	 and desktop access", "Easy migration from your current provider", "Email scheduling, templates, and autoresponders", "Access across all devices", "Hassle-free migration tools", "Advanced tools: scheduling, templates, autoresponders"},
			EveryNDays:              365,
			BillingCycles:           10,
			Level:                   2,
			TokenRewardEveryRenewal: 0,
			ExpiresAfterDayCount:    1095,
			CurrencySymbol:          "₹",
			ForeignCurrencyPricing:  nil,
			Recommended:             false,
		},
	}

	for _, subscription := range subscriptions {
		subscription.ProductPlanID = product.ID
		result := pps.DB.Create(&subscription)
		if result.Error != nil {
			panic(fmt.Sprintf("Failed to seed subscription data: %v", result.Error))
		}
		// if result.RowsAffected == 0 {
		// 	panic("No rows were affected while seeding subscription data")
		// }
	}

	fmt.Println("Product data seeded successfully")
}
