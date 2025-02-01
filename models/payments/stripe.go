package payments

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
