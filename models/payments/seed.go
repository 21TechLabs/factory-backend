package payments

import (
	"fmt"

	"github.com/kamva/mgm/v3"
)

func SeedProductData() {
	// Seed the payment gateway
	var product ProductPlans = ProductPlans{
		AppCode:       "app1",
		Name:          "Product 1",
		Description:   "Product 1 Description",
		Type:          "subscription",
		Amount:        0,
		TokenReward:   0,
		Subscriptions: []Subscription{},
	}

	// Seed the subscription
	var subscription Subscription = Subscription{
		Name:                    "Subscription 1",
		Description:             "Subscription 1 Description",
		PlanID:                  map[string]string{"razorpay": "plan_1", "stripe": "plan_1"},
		Amount:                  100,
		Features:                []string{"Feature 1", "Feature 2"},
		EveryNDays:              30,
		Level:                   1,
		TokenRewardEveryRenewal: 0,
		BillingCycles:           36,
		ExpiresAfterDayCount:    365 * 3,
	}
	product.Subscriptions = append(product.Subscriptions, subscription)

	err := mgm.Coll(&ProductPlans{}).Create(&product)

	if err != nil {
		panic(err)
	}
	fmt.Println("Product data seeded successfully")
}
