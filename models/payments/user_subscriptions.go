package payments

import (
	"time"

	"github.com/kamva/mgm/v3"
)

type UserSubscription struct {
	mgm.DefaultModel         `bson:",inline"`
	UserId                   string             `json:"userId" bson:"userId"`
	Plan                     string             `json:"plan" bson:"plan"`
	SubscriptionStartsAt     time.Time          `json:"subscriptionStartsAt" bson:"subscriptionStartsAt"`
	SubscriptionEndsAt       time.Time          `json:"subscriptionEndsAt" bson:"subscriptionEndsAt"`
	ReSubscribeOn            time.Time          `json:"reSubscribeOn" bson:"reSubscribeOn"`
	SubscriptionTenureInDays int                `json:"subscriptionTenureInDays" bson:"subscriptionTenureInDays"`
	Status                   SubscriptionStatus `json:"status" bson:"status"`
	PlanLevel                int                `json:"planLevel" bson:"planLevel"`
	TokenReward              int64              `json:"tokenReward" bson:"tokenReward"`
	PaymentId                map[string]string  `json:"paymentId" bson:"paymentId"`
	AppCode                  string             `json:"appCode" bson:"appCode"`
	HasExpired               bool               `json:"hasExpired" bson:"hasExpired"`
	PaymentMethod            string             `json:"paymentMethod" bson:"paymentMethod"`
	ExtraOptions             map[string]string  `json:"extraOptions" bson:"extraOptions"`
}

func (us *UserSubscription) Save(update bool) error {
	if update {
		return mgm.Coll(us).Update(us)
	}
	return mgm.Coll(us).Create(us)
}
