package models

import (
	"net/http"
)

type Stripe struct {
	UserId uint
}

func (s *Stripe) CreatePayment(uss *UserSubscriptionStore, productPlan ProductPlan, planIdx int) (UserSubscription, error) {
	return UserSubscription{}, nil
}

func (s *Stripe) VerifyPayment(string) (UserSubscription, error) {
	return UserSubscription{}, nil
}

func (s *Stripe) CancelSubscription(string) (UserSubscription, error) {
	return UserSubscription{}, nil
}

func (s *Stripe) GetUserSubscription() (interface{}, error) {
	return UserSubscription{}, nil
}

func (s *Stripe) UpdatePaymentStatus(uss *UserSubscriptionStore, orderId string) (UserSubscription, error) {
	return UserSubscription{}, nil
}

func (s *Stripe) VerifyWebhookSignature(*http.Request) error {
	return nil
}

func (s *Stripe) GetOrderIdFromWebhookRequest([]byte) (string, error) {
	return "", nil
}

func (s *Stripe) SetUserId(userId uint) {
	s.UserId = userId
}

func (s *Stripe) SetUserViaOrderId(uss *UserSubscriptionStore, orderId string) error {
	var userSubColl = uss.DB.Model(&UserSubscription{})

	var userSubscription UserSubscription

	result := userSubColl.Where("payment_id->? = ?", PaymentGatewaysList.Stripe, orderId).First(&userSubscription)

	if result.Error != nil {
		return result.Error
	}

	s.SetUserId(userSubscription.UserID)
	return nil
}
