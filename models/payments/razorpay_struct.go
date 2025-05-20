package payments

import "time"

type Razorpay struct {
	UserId string
}

type RazorpaySubscriptionCreate struct {
	AuthAttempts        int                `json:"auth_attempts"`
	ChangeScheduledAt   *time.Time         `json:"change_scheduled_at"`
	ChargeAt            *time.Time         `json:"charge_at"`
	CreatedAt           int64              `json:"created_at"`
	CurrentEnd          *time.Time         `json:"current_end"`
	CurrentStart        *time.Time         `json:"current_start"`
	CustomerNotify      bool               `json:"customer_notify"`
	EndAt               *time.Time         `json:"end_at"`
	EndedAt             *time.Time         `json:"ended_at"`
	Entity              string             `json:"entity"`
	ExpireBy            *time.Time         `json:"expire_by"`
	HasScheduledChanges bool               `json:"has_scheduled_changes"`
	ID                  string             `json:"id"`
	Notes               []string           `json:"notes"`
	PaidCount           int                `json:"paid_count"`
	PlanID              string             `json:"plan_id"`
	Quantity            int                `json:"quantity"`
	RemainingCount      int                `json:"remaining_count"`
	ShortURL            string             `json:"short_url"`
	Source              string             `json:"source"`
	StartAt             *time.Time         `json:"start_at"`
	Status              SubscriptionStatus `json:"status"`
	TotalCount          int                `json:"total_count"`
}

type RazorpaySubscriptionFetch struct {
	ID                  string             `json:"id"`
	Entity              string             `json:"entity"`
	PlanID              string             `json:"plan_id"`
	CustomerID          string             `json:"customer_id"`
	Status              SubscriptionStatus `json:"status"`
	CurrentStart        int64              `json:"current_start"`
	CurrentEnd          int64              `json:"current_end"`
	EndedAt             *int64             `json:"ended_at"`
	Quantity            int                `json:"quantity"`
	Notes               []string           `json:"notes"`
	ChargeAt            int64              `json:"charge_at"`
	StartAt             int64              `json:"start_at"`
	EndAt               int64              `json:"end_at"`
	AuthAttempts        int                `json:"auth_attempts"`
	TotalCount          int                `json:"total_count"`
	PaidCount           int                `json:"paid_count"`
	CustomerNotify      bool               `json:"customer_notify"`
	CreatedAt           int64              `json:"created_at"`
	ExpireBy            int64              `json:"expire_by"`
	ShortURL            string             `json:"short_url"`
	HasScheduledChanges bool               `json:"has_scheduled_changes"`
	ChangeScheduledAt   *int64             `json:"change_scheduled_at"`
	Source              string             `json:"source"`
	OfferID             string             `json:"offer_id"`
	RemainingCount      int                `json:"remaining_count"`
}

type RazorpayPayment struct {
	ID                string              `json:"id"`
	Entity            string              `json:"entity"`
	Amount            int                 `json:"amount"`
	Currency          string              `json:"currency"`
	Status            SubscriptionStatus  `json:"status"`
	OrderID           *string             `json:"order_id"`
	InvoiceID         *string             `json:"invoice_id"`
	International     bool                `json:"international"`
	Method            string              `json:"method"`
	AmountRefunded    int                 `json:"amount_refunded"`
	AmountTransferred int                 `json:"amount_transferred"`
	RefundStatus      *string             `json:"refund_status"`
	Captured          string              `json:"captured"`
	Description       string              `json:"description"`
	CardID            string              `json:"card_id"`
	Card              RazorpayPaymentCard `json:"card"`
	Bank              *string             `json:"bank"`
	Wallet            *string             `json:"wallet"`
	VPA               *string             `json:"vpa"`
	Email             string              `json:"email"`
	Contact           string              `json:"contact"`
	CustomerID        string              `json:"customer_id"`
	TokenID           *string             `json:"token_id"`
	Notes             []string            `json:"notes"`
	Fee               int                 `json:"fee"`
	Tax               int                 `json:"tax"`
	ErrorCode         *string             `json:"error_code"`
	ErrorDescription  *string             `json:"error_description"`
	CreatedAt         int64               `json:"created_at"`
}

type RazorpayPaymentCard struct {
	ID            string `json:"id"`
	Entity        string `json:"entity"`
	Name          string `json:"name"`
	Last4         string `json:"last4"`
	Network       string `json:"network"`
	Type          string `json:"type"`
	Issuer        string `json:"issuer"`
	International bool   `json:"international"`
	EMI           bool   `json:"emi"`
	ExpiryMonth   any    `json:"expiry_month"`
	ExpiryYear    any    `json:"expiry_year"`
}

type RazorpayWebhookPayload struct {
	Subscription struct {
		Entity RazorpaySubscriptionFetch `json:"entity"`
	} `json:"subscription"`
	Payment struct {
		Entity RazorpayPayment `json:"entity"`
	} `json:"payment"`
}

type RazorpaySubsctiptionWebhook struct {
	Entity    string                 `json:"entity"`
	AccountID string                 `json:"account_id"`
	Event     string                 `json:"event"`
	Contains  []string               `json:"contains"`
	Payload   RazorpayWebhookPayload `json:"payload"`
	CreatedAt int64                  `json:"created_at"`
}
