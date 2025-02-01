package payments

import (
	"errors"

	"github.com/21TechLabs/factory-be/utils"
	"github.com/kamva/mgm/v3"
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

type ProductPlans struct {
	mgm.DefaultModel `bson:",inline"`
	AppCode          string         `json:"appCode" bson:"appCode"`
	Name             string         `json:"name" bson:"name"`
	Description      string         `json:"description" bson:"description"`
	Type             string         `json:"type" bson:"type"`
	Amount           float64        `json:"amount" bson:"amount"`
	TokenReward      float64        `json:"tokenAmount" bson:"tokenAmount"`
	Subscriptions    []Subscription `json:"subscriptions" bson:"subscriptions"`
}

type Subscription struct {
	mgm.DefaultModel        `bson:",inline"`
	Name                    string            `json:"name" bson:"name"`
	Description             string            `json:"description" bson:"description"`
	PlanID                  map[string]string `json:"planId" bson:"planId"`
	Amount                  float64           `json:"amount" bson:"amount"`
	Features                []string          `json:"features" bson:"features"`
	EveryNDays              int               `json:"everyNDays" bson:"everyNDays"`
	BillingCycles           int               `json:"billingCycle" bson:"billingCycle"`
	Level                   int               `json:"level" bson:"level"`
	TokenRewardEveryRenewal int64             `json:"tokenReward" bson:"tokenReward"`
	ExpiresAfterDayCount    int               `json:"expiresAfterDayCount" bson:"expiresAfterDayCount"`
}

func ProductPlansGetByID(id string) (ProductPlans, error) {

	if !utils.IsValidObjectID(id) {
		return ProductPlans{}, errors.New("invalid id")
	}

	var productPlan ProductPlans = ProductPlans{}

	coll := mgm.Coll(&ProductPlans{})
	if coll == nil {
		return ProductPlans{}, errors.New("database connection not initialized")
	}

	err := coll.FindByID(id, &productPlan)

	if err != nil {
		return ProductPlans{}, err
	}

	if productPlan.ID.IsZero() {
		return productPlan, errors.New("product plan not found")
	}

	return productPlan, err
}

func GetPaymentGateway(paymentType string, userId string) (PaymentGateway, error) {
	var payment PaymentGateway
	switch paymentType {
	case PaymentGatewaysList.Razorpay:
		payment = &Razorpay{UserId: userId}
	case PaymentGatewaysList.Stripe:
		payment = &Stripe{UserId: userId}
	default:
		return nil, errors.New("invalid payment type")
	}
	return payment, nil
}
