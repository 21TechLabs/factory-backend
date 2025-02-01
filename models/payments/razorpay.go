package payments

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/21TechLabs/factory-be/utils"
	"github.com/razorpay/razorpay-go"
)

type Razorpay struct {
	UserId string
}

type RazorpaySubscriptionCreate struct {
	AuthAttempts        int        `json:"auth_attempts"`
	ChangeScheduledAt   *time.Time `json:"change_scheduled_at"`
	ChargeAt            *time.Time `json:"charge_at"`
	CreatedAt           int64      `json:"created_at"`
	CurrentEnd          *time.Time `json:"current_end"`
	CurrentStart        *time.Time `json:"current_start"`
	CustomerNotify      bool       `json:"customer_notify"`
	EndAt               *time.Time `json:"end_at"`
	EndedAt             *time.Time `json:"ended_at"`
	Entity              string     `json:"entity"`
	ExpireBy            *time.Time `json:"expire_by"`
	HasScheduledChanges bool       `json:"has_scheduled_changes"`
	ID                  string     `json:"id"`
	Notes               []string   `json:"notes"`
	PaidCount           int        `json:"paid_count"`
	PlanID              string     `json:"plan_id"`
	Quantity            int        `json:"quantity"`
	RemainingCount      int        `json:"remaining_count"`
	ShortURL            string     `json:"short_url"`
	Source              string     `json:"source"`
	StartAt             *time.Time `json:"start_at"`
	Status              string     `json:"status"`
	TotalCount          int        `json:"total_count"`
}

var razorpayClient *razorpay.Client

func init() {
	utils.LoadEnv()
	razorpayClient = razorpay.NewClient(utils.GetEnv("RAZORPAY_KEY_ID", false), utils.GetEnv("RAZORPAY_KEY_SECRET", false))
}

func (sub *Razorpay) CreatePayment(productPlan ProductPlans, planIdx int) (UserSubscription, error) {
	if planIdx >= len(productPlan.Subscriptions) || planIdx < 0 {
		log.Printf("Payment gateway razorpay.CreatePayment error: %v", "Plan index out of range")
		return UserSubscription{}, errors.New("plan index out of range")
	}

	plan := productPlan.Subscriptions[planIdx]

	data := map[string]interface{}{
		"plan_id":         plan.PlanID[PaymentGatewaysList.Razorpay],
		"total_count":     plan.BillingCycles,
		"quantity":        1,
		"customer_notify": 1,
		"addons": []interface{}{
			map[string]interface{}{},
		},
		"notes": map[string]interface{}{},
	}

	body, err := razorpayClient.Subscription.Create(data, nil)

	if err != nil {
		log.Printf("Payment gateway create error: %v", err)
		return UserSubscription{}, fmt.Errorf("failed to create Razorpay subscription: %w", err)
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error converting map to JSON:", err)
		return UserSubscription{}, fmt.Errorf("failed to create Razorpay subscription: %w", err)
	}

	var razorpaySubscription RazorpaySubscriptionCreate

	if err = json.Unmarshal(jsonData, &razorpaySubscription); err != nil {
		log.Printf("Payment gateway create error: %v", err)
		return UserSubscription{}, fmt.Errorf("failed to create Razorpay subscription: %w", err)
	}

	var userSubscription = UserSubscription{
		UserId:                   sub.UserId,
		Plan:                     plan.Name,
		Status:                   razorpaySubscription.Status,
		PlanLevel:                plan.Level,
		SubscriptionTenureInDays: plan.ExpiresAfterDayCount,
		PaymentId:                map[string]string{PaymentGatewaysList.Razorpay: razorpaySubscription.ID},
		AppCode:                  productPlan.AppCode,
		TokenReward:              plan.TokenRewardEveryRenewal,
		HasExpired:               false,
		ExtraOptions: map[string]string{
			"rpay_redirect_url": razorpaySubscription.ShortURL,
		},
		PaymentMethod: PaymentGatewaysList.Razorpay,
	}
	userSubscription.Save()

	return userSubscription, nil
}

func (s *Razorpay) VerifyPayment(string) (UserSubscription, error) {
	return UserSubscription{}, nil
}

func (s *Razorpay) CancelSubscription(string) (UserSubscription, error) {
	return UserSubscription{}, nil
}

func (s *Razorpay) GetUserSubscription() (interface{}, error) {
	return UserSubscription{}, nil
}
