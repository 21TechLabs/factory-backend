package payments

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/21TechLabs/factory-be/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"github.com/razorpay/razorpay-go"
	"go.mongodb.org/mongo-driver/bson"
)

var razorpayClient *razorpay.Client

func init() {
	utils.LoadEnv()
	razorpayClient = razorpay.NewClient(utils.GetEnv("RAZORPAY_KEY_ID", false), utils.GetEnv("RAZORPAY_KEY_SECRET", false))
	razorpayClient.SetTimeout(3)
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
	userSubscription.Save(false)

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

func (s *Razorpay) UpdatePaymentStatus(subscriptionId string) (UserSubscription, error) {
	// fetch subscription by order id
	body, err := razorpayClient.Subscription.Fetch(subscriptionId, nil, nil)

	if err != nil {
		log.Println("Payment gateway update error razorpay.UpdatePaymentStatus:", err.Error())
		return UserSubscription{}, fmt.Errorf("failed to update Razorpay subscription: %w", err)
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		log.Println("Error converting map to JSON razorpay.UpdatePaymentStatus: ", err)
		return UserSubscription{}, fmt.Errorf("failed to update Razorpay subscription: %w", err)
	}

	var razorpaySubscriptionWebhook RazorpaySubscriptionFetch

	if err := json.Unmarshal(jsonData, &razorpaySubscriptionWebhook); err != nil {
		log.Println("Payment gateway update error razorpay.UpdatePaymentStatus: ", err)
		return UserSubscription{}, fmt.Errorf("failed to update Razorpay subscription: %w", err)
	}

	var userSubscription UserSubscription

	var userSubColl = mgm.Coll(&UserSubscription{})
	if err := userSubColl.First(bson.M{"paymentId." + PaymentGatewaysList.Razorpay: subscriptionId}, &userSubscription); err != nil {
		log.Println("Payment gateway update error razorpay.UpdatePaymentStatus: ", err)
		return UserSubscription{}, fmt.Errorf("failed to update Razorpay subscription: %w", err)
	}

	if userSubscription.ID.IsZero() {
		return UserSubscription{}, errors.New("subscription not found")
	}

	userSubscription.SubscriptionStartsAt = time.Unix(razorpaySubscriptionWebhook.StartAt, 0)
	userSubscription.ReSubscribeOn = time.Unix(razorpaySubscriptionWebhook.StartAt, 0)

	switch razorpaySubscriptionWebhook.Status {
	case SubscriptionStatusActive:
		userSubscription.Status = SubscriptionStatusActive
		userSubscription.SubscriptionStartsAt = time.Unix(razorpaySubscriptionWebhook.CurrentStart, 0)
		userSubscription.ReSubscribeOn = time.Unix(razorpaySubscriptionWebhook.CurrentStart, 0)
		userSubscription.SubscriptionEndsAt = time.Unix(razorpaySubscriptionWebhook.CurrentEnd, 0)
		userSubscription.HasExpired = false
	case SubscriptionStatusCompleted:
		userSubscription.Status = SubscriptionStatusCompleted
		userSubscription.HasExpired = false
		userSubscription.SubscriptionStartsAt = time.Unix(razorpaySubscriptionWebhook.CurrentStart, 0)
		userSubscription.ReSubscribeOn = time.Unix(razorpaySubscriptionWebhook.CurrentStart, 0)
		userSubscription.SubscriptionEndsAt = time.Unix(razorpaySubscriptionWebhook.CurrentEnd, 0)
	case SubscriptionStatusCharged:
		userSubscription.Status = SubscriptionStatusCharged
		userSubscription.HasExpired = false
		userSubscription.SubscriptionStartsAt = time.Unix(razorpaySubscriptionWebhook.CurrentStart, 0)
		userSubscription.ReSubscribeOn = time.Unix(razorpaySubscriptionWebhook.CurrentStart, 0)
		userSubscription.SubscriptionEndsAt = time.Unix(razorpaySubscriptionWebhook.CurrentEnd, 0)
	case SubscriptionStatusCancelled:
		userSubscription.Status = SubscriptionStatusCancelled
		userSubscription.HasExpired = true
	case SubscriptionStatusHalted:
		userSubscription.Status = SubscriptionStatusHalted
		userSubscription.HasExpired = true
	}

	userSubscription.Save(true)

	return UserSubscription{}, nil
}

func (s *Razorpay) VerifyWebhookSignature(c *fiber.Ctx) error {
	headerSignature := c.Get("x-razorpay-signature")

	var secret = utils.GetEnv("PAYMENTS_HMEC_SECRET", false)

	if !utils.ValidateHeaderHMACSha256(c.Body(), secret, headerSignature) {
		return errors.New("invalid signature")
	}

	return nil
}

func (s *Razorpay) GetOrderIdFromWebhookRequest(body []byte) (string, error) {
	var data RazorpaySubsctiptionWebhook
	fmt.Println(string(body))
	if err := json.Unmarshal(body, &data); err != nil {
		log.Printf("Payment gateway create error -- Get Order Id from Webhook Request controller.payments.GetOrderIdFromWebhookRequest.razorpay: %v", err)
		return "", err
	}
	return data.Payload.Subscription.Entity.ID, nil
}

func (s *Razorpay) SetUserId(userId string) {
	s.UserId = userId
}

func (s *Razorpay) SetUserViaOrderId(orderId string) error {
	var userSubColl = mgm.Coll(&UserSubscription{})

	var userSubscription UserSubscription

	if err := userSubColl.First(bson.M{"paymentId." + PaymentGatewaysList.Razorpay: orderId}, &userSubscription); err != nil {
		log.Printf("Payment gateway create error -- Get User Subscription controller.payments.SetUserViaOrderId.razorpay: %v", err)
		return err
	}

	s.SetUserId(userSubscription.UserId)
	return nil
}
