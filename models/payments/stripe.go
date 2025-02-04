package payments

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

type Stripe struct {
	UserId string
}

func (s *Stripe) CreatePayment(productPlan ProductPlans, planIdx int) (UserSubscription, error) {
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

func (s *Stripe) UpdatePaymentStatus(orderId string) (UserSubscription, error) {
	return UserSubscription{}, nil
}

func (s *Stripe) VerifyWebhookSignature(*fiber.Ctx) error {
	return nil
}

func (s *Stripe) GetOrderIdFromWebhookRequest([]byte) (string, error) {
	return "", nil
}

func (s *Stripe) SetUserId(userId string) {
	s.UserId = userId
}

func (s *Stripe) SetUserViaOrderId(orderId string) error {
	var userSubColl = mgm.Coll(&UserSubscription{})

	var userSubscription UserSubscription

	if err := userSubColl.First(bson.M{"paymentId." + PaymentGatewaysList.Stripe: orderId}, &userSubscription); err != nil {
		log.Printf("Payment gateway create error -- Get User Subscription controller.payments.SetUserViaOrderId.stripe: %v", err)
		return err
	}

	s.SetUserId(userSubscription.UserId)
	return nil
}
