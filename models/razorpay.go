package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/21TechLabs/factory-backend/utils"
	"github.com/razorpay/razorpay-go"
)

var razorpayClient *razorpay.Client

func init() {
	utils.LoadEnv()
	razorpayClient = razorpay.NewClient(utils.GetEnv("RAZORPAY_KEY_ID", false), utils.GetEnv("RAZORPAY_KEY_SECRET", false))
	razorpayClient.SetTimeout(3)
}

func (sub *Razorpay) CreatePayment(uss *UserSubscriptionStore, productPlan ProductPlan, planIdx int) (UserSubscription, error) {
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
		UserID:                   sub.UserId,
		Plan:                     plan.Name,
		Status:                   razorpaySubscription.Status,
		PlanLevel:                plan.Level,
		SubscriptionTenureInDays: plan.ExpiresAfterDayCount,
		PaymentId:                utils.JSONMap{PaymentGatewaysList.Razorpay: razorpaySubscription.ID},
		AppCode:                  productPlan.AppCode,
		TokenReward:              plan.TokenRewardEveryRenewal,
		HasExpired:               false,
		ExtraOptions: utils.JSONMap{
			"rpay_redirect_url": razorpaySubscription.ShortURL,
		},
		PaymentMethod: PaymentGatewaysList.Razorpay,
	}
	uss.Save(&userSubscription, false)

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

func (s *Razorpay) UpdatePaymentStatus(uss *UserSubscriptionStore, subscriptionId string) (UserSubscription, error) {
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

	var userSubColl = uss.DB.Model(&UserSubscription{})

	result := userSubColl.Where("payment_id @> ?", fmt.Sprintf("{\"%s\": \"%s\"}", PaymentGatewaysList.Razorpay, subscriptionId)).First(&userSubscription)
	if result.Error != nil {
		log.Println("Payment gateway update error razorpay.UpdatePaymentStatus: ", result.Error)
		return UserSubscription{}, fmt.Errorf("failed to update Razorpay subscription: %w", result.Error)
	}

	userSubscription.SubscriptionStartsAt = time.Unix(razorpaySubscriptionWebhook.StartAt, 0)
	userSubscription.ReSubscribeOn = time.Unix(razorpaySubscriptionWebhook.StartAt, 0)

	switch razorpaySubscriptionWebhook.Status {
	case SubscriptionStatusActive:
		userSubscription.Status = SubscriptionStatusActive
		userSubscription.HasExpired = false
		userSubscription.SubscriptionStartsAt = time.Unix(razorpaySubscriptionWebhook.CurrentStart, 0)
		userSubscription.ReSubscribeOn = time.Unix(razorpaySubscriptionWebhook.CurrentStart, 0)
		userSubscription.SubscriptionEndsAt = time.Unix(razorpaySubscriptionWebhook.CurrentEnd, 0)
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

	uss.Save(&userSubscription, true)

	return UserSubscription{}, nil
}

func (s *Razorpay) VerifyWebhookSignature(r *http.Request) error {
	headerSignature := r.Header.Get("x-razorpay-signature")
	if headerSignature == "" {
		log.Printf("Payment gateway create error -- VerifyWebhookSignature controller.payments.VerifyWebhookSignature.razorpay: %v", "Header signature not found")
		return errors.New("header signature not found")
	}

	var secret = utils.GetEnv("PAYMENTS_HMEC_SECRET", false)

	if secret == "" {
		log.Printf("Payment gateway create error -- VerifyWebhookSignature controller.payments.VerifyWebhookSignature.razorpay: %v", "HMAC secret not found")
		return errors.New("HMAC secret not found")
	}

	if r.Body == nil {
		log.Printf("Payment gateway create error -- VerifyWebhookSignature controller.payments.VerifyWebhookSignature.razorpay: %v", "Request body is empty")
		return errors.New("request body is empty")
	}

	c, err := r.GetBody()
	if err != nil {
		log.Printf("Payment gateway create error -- VerifyWebhookSignature controller.payments.VerifyWebhookSignature.razorpay: %v", err)
		return err
	}

	bodyBytes := make([]byte, r.ContentLength)
	_, err = c.Read(bodyBytes)
	if err != nil {
		log.Printf("Payment gateway create error -- VerifyWebhookSignature controller.payments.VerifyWebhookSignature.razorpay: %v", err)
		return err
	}

	if len(bodyBytes) == 0 {
		log.Printf("Payment gateway create error -- VerifyWebhookSignature controller.payments.VerifyWebhookSignature.razorpay: %v", "Request body is empty")
		return errors.New("request body is empty")
	}

	if !utils.ValidateHeaderHMACSha256(bodyBytes, secret, headerSignature) {
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

func (s *Razorpay) SetUserId(userId uint) {
	s.UserId = userId
}

func (s *Razorpay) SetUserViaOrderId(uss *UserSubscriptionStore, orderId string) error {
	var userSubColl = uss.DB.Model(&UserSubscription{})

	var userSubscription UserSubscription

	result := userSubColl.Where("payment_id @> ?", fmt.Sprintf("{\"%s\": \"%s\"}", PaymentGatewaysList.Razorpay, orderId)).First(&userSubscription)
	if result.Error != nil {
		log.Printf("Payment gateway create error -- Get User Subscription controller.payments.SetUserViaOrderId.razorpay: %v", result.Error)
		return result.Error
	}

	s.SetUserId(userSubscription.UserID)
	return nil
}
