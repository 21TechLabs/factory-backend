package utils

import "slices"

type Number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

type Currency string

const (
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
	CurrencyINR Currency = "INR"
)

type SubscriptionStatus string

// Constants for various subscription statuses.
const (
	SubscriptionStatusActive    SubscriptionStatus = "subscription.active"
	SubscriptionStatusPending   SubscriptionStatus = "subscription.pending"
	SubscriptionStatusHalted    SubscriptionStatus = "subscription.halted"
	SubscriptionStatusCancelled SubscriptionStatus = "subscription.cancelled"
	SubscriptionStatusCompleted SubscriptionStatus = "subscription.completed"
	SubscriptionStatusPaused    SubscriptionStatus = "subscription.paused"
	SubscriptionStatusResumed   SubscriptionStatus = "subscription.resumed"
	SubscriptionStatusCharged   SubscriptionStatus = "subscription.charged"
)

var SubscriptionHooks []SubscriptionStatus = []SubscriptionStatus{
	SubscriptionStatusActive,
	SubscriptionStatusPending,
	SubscriptionStatusHalted,
	SubscriptionStatusCancelled,
	SubscriptionStatusCompleted,
	SubscriptionStatusPaused,
	SubscriptionStatusResumed,
	SubscriptionStatusCharged,
}

func SubscriptionWebhookIsValid(webhook string) bool {
	return webhook == "order.paid" || webhook == "subscription.created" || slices.Contains[[]SubscriptionStatus, SubscriptionStatus](SubscriptionHooks, SubscriptionStatus(webhook))
}

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
)

type PlanType string

const (
	PlanTypeSubscription PlanType = "subscription"
	PlanTypeOneTime      PlanType = "one_time"
)

var PlanTypes map[string]PlanType = map[string]PlanType{
	string(PlanTypeSubscription): PlanTypeSubscription,
	string(PlanTypeOneTime):      PlanTypeOneTime,
}

func (pt PlanType) IsValid() bool {
	_, exists := PlanTypes[string(pt)]
	return exists
}
