package utils

type Currency string

const (
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
	CurrencyINR Currency = "INR"
)

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
