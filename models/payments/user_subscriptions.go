package payments

import (
	"time"

	"github.com/kamva/mgm/v3"
)

type UserSubscription struct {
	mgm.DefaultModel         `bson:",inline"`
	UserId                   string            `json:"userId" bson:"userId"`
	Plan                     string            `json:"plan" bson:"plan"`
	SubscriptionStartsAt     time.Time         `json:"subscriptionStartsAt" bson:"subscriptionStartsAt"`
	SubscriptionEndsAt       time.Time         `json:"subscriptionEndsAt" bson:"subscriptionEndsAt"`
	ReSubscribeOn            time.Time         `json:"reSubscribeOn" bson:"reSubscribeOn"`
	SubscriptionTenureInDays int               `json:"subscriptionTenureInDays" bson:"subscriptionTenureInDays"`
	Status                   string            `json:"status" bson:"status"`
	PlanLevel                int               `json:"planLevel" bson:"planLevel"`
	TokenReward              int64             `json:"tokenReward" bson:"tokenReward"`
	PaymentId                map[string]string `json:"paymentId" bson:"paymentId"`
	AppCode                  string            `json:"appCode" bson:"appCode"`
	HasExpired               bool              `json:"hasExpired" bson:"hasExpired"`
	PaymentMethod            string            `json:"paymentMethod" bson:"paymentMethod"`
	ExtraOptions             map[string]string `json:"extraOptions" bson:"extraOptions"`
}

type PaymentGateway interface {
	CreatePayment(productPlan ProductPlans, planIdx int) (UserSubscription, error)
	VerifyPayment(orderId string) (UserSubscription, error)
	CancelSubscription(string) (UserSubscription, error)
	GetUserSubscription() (interface{}, error)
}

func (us *UserSubscription) GetUserSubscription(AppCode string) (UserSubscription, error) {
	return UserSubscription{}, nil
}

func (us *UserSubscription) Save() error {
	return mgm.Coll(us).Create(us)
}
