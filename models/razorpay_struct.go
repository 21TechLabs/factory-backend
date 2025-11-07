package models

import "github.com/21TechLabs/factory-backend/utils"

// RazorpayBaseEvent represents the common outer structure for all Razorpay webhook events.
type RazorpayBaseEvent[T any] struct {
	Entity    string   `json:"entity"`
	AccountID string   `json:"account_id"`
	Event     string   `json:"event"`
	Contains  []string `json:"contains"`
	Payload   T        `json:"payload"`
	CreatedAt int64    `json:"created_at"`
}

type OrderStatus string

const (
	OrderStatusPaid   OrderStatus = "paid"
	OrderStatusFailed OrderStatus = "failed"
)

// RazorpayOrderPaidPayload contains the specific entities for the "order.paid" event.
type RazorpayOrderPaidPayload struct {
	Payment RazorpayPaymentWrapper `json:"payment"`
	Order   RazorpayOrderWrapper   `json:"order"`
}

// RazorpayPaymentWrapper wraps the main RazorpayPaymentEntity.
type RazorpayPaymentWrapper struct {
	Entity RazorpayPaymentEntity `json:"entity"`
}

// RazorpayOrderWrapper wraps the main RazorpayOrderEntity.
type RazorpayOrderWrapper struct {
	Entity RazorpayOrderEntity `json:"entity"`
}

// RazorpayCardDetails holds information about the card used for payment.
// This is present only when the payment method is "card".
type RazorpayCardDetails struct {
	ID            string `json:"id"`
	Entity        string `json:"entity"`
	Name          string `json:"name"`
	Last4         string `json:"last4"`
	Network       string `json:"network"`
	Type          string `json:"type"`
	Issuer        any    `json:"issuer"` // Can be null
	International bool   `json:"international"`
	EMI           bool   `json:"emi"`
	ExpiryMonth   any    `json:"expiry_month,omitempty"` // Added field
	ExpiryYear    any    `json:"expiry_year,omitempty"`  // Added field
}

// RazorpayPaymentEntity defines the structure of the payment object.
// Fields that vary by payment method (Bank, Wallet, VPA, Card) are pointers
// to handle cases where they are null.
type RazorpayPaymentEntity struct {
	ID               string               `json:"id"`
	Entity           string               `json:"entity"`
	Amount           int                  `json:"amount"`
	Currency         string               `json:"currency"`
	Status           string               `json:"status"`
	OrderID          string               `json:"order_id"`
	InvoiceID        any                  `json:"invoice_id"` // Can be null
	International    bool                 `json:"international"`
	Method           string               `json:"method"`
	AmountRefunded   int                  `json:"amount_refunded"`
	RefundStatus     any                  `json:"refund_status"` // Can be null
	Captured         bool                 `json:"captured"`
	Description      any                  `json:"description"` // Can be null
	CardID           *string              `json:"card_id"`     // Can be null
	Bank             *string              `json:"bank"`        // Can be null, e.g., for card payments
	Wallet           *string              `json:"wallet"`      // Can be null, e.g., for netbanking
	VPA              *string              `json:"vpa"`         // Can be null, e.g., for card payments
	Email            string               `json:"email"`
	Contact          string               `json:"contact"`
	Notes            []any                `json:"notes"` // Assuming it can contain various types
	Fee              int                  `json:"fee"`
	Tax              int                  `json:"tax"`
	ErrorCode        any                  `json:"error_code"`        // Can be null
	ErrorDescription any                  `json:"error_description"` // Can be null
	Card             *RazorpayCardDetails `json:"card,omitempty"`    // Present only for card method
	CreatedAt        int64                `json:"created_at"`
}

// RazorpayOrderEntity defines the structure of the order object.
type RazorpayOrderEntity struct {
	ID         string      `json:"id"`
	Entity     string      `json:"entity"`
	Amount     int         `json:"amount"`
	AmountPaid int         `json:"amount_paid"`
	AmountDue  int         `json:"amount_due"`
	Currency   string      `json:"currency"`
	Receipt    string      `json:"receipt"`
	OfferID    any         `json:"offer_id"` // Can be null
	Status     OrderStatus `json:"status"`
	Attempts   int         `json:"attempts"`
	Notes      []any       `json:"notes"`
	CreatedAt  int64       `json:"created_at"`
}

// RazorpayPaymentFailedPayload contains the payload for the "payment.failed" event.
// It uses a specific wrapper for the failed payment entity.
type RazorpayPaymentFailedPayload struct {
	Payment RazorpayPaymentFailedWrapper `json:"payment"`
}

// RazorpayPaymentFailedWrapper wraps the main RazorpayPaymentFailedEntity.
type RazorpayPaymentFailedWrapper struct {
	Entity RazorpayPaymentFailedEntity `json:"entity"`
}

// RazorpayAcquirerData holds transaction data from the bank or acquirer.
// Fields are pointers to handle cases where they are null.
type RazorpayAcquirerData struct {
	BankTransactionID *string `json:"bank_transaction_id,omitempty"`
	AuthCode          *string `json:"auth_code,omitempty"`
	RRN               *string `json:"rrn,omitempty"`
	TransactionID     *string `json:"transaction_id,omitempty"`
}

// RazorpayUPIDetails holds information specific to a UPI transaction.
// This is present only when the payment method is "upi".
type RazorpayUPIDetails struct {
	PayerAccountType string `json:"payer_account_type"`
	VPA              string `json:"vpa"`
	Flow             string `json:"flow"`
}

// RazorpayPaymentFailedEntity defines the structure of the payment object for a failed event.
// It includes detailed error fields and method-specific nested objects.
type RazorpayPaymentFailedEntity struct {
	ID               string               `json:"id"`
	Entity           string               `json:"entity"`
	Amount           int                  `json:"amount"`
	Currency         string               `json:"currency"`
	Status           string               `json:"status"`
	OrderID          string               `json:"order_id"`
	InvoiceID        any                  `json:"invoice_id"`
	International    bool                 `json:"international"`
	Method           string               `json:"method"`
	AmountRefunded   int                  `json:"amount_refunded"`
	RefundStatus     any                  `json:"refund_status"`
	Captured         bool                 `json:"captured"`
	Description      any                  `json:"description"`
	CardID           *string              `json:"card_id"`
	Bank             *string              `json:"bank"`
	Wallet           *string              `json:"wallet"`
	VPA              *string              `json:"vpa"`
	Email            string               `json:"email"`
	Contact          string               `json:"contact"`
	Notes            []any                `json:"notes"`
	Fee              any                  `json:"fee"` // Fee is null for failed payments
	Tax              any                  `json:"tax"` // Tax is null for failed payments
	ErrorCode        string               `json:"error_code"`
	ErrorDescription string               `json:"error_description"`
	ErrorSource      *string              `json:"error_source"`
	ErrorStep        *string              `json:"error_step"`
	ErrorReason      *string              `json:"error_reason"`
	AcquirerData     RazorpayAcquirerData `json:"acquirer_data"`
	CreatedAt        int64                `json:"created_at"`
	TokenID          *string              `json:"token_id,omitempty"` // Present in some card transactions

	// Method-specific Objects
	Card *RazorpayCardDetails `json:"card,omitempty"`
	UPI  *RazorpayUPIDetails  `json:"upi,omitempty"`
}

type RazorpaySubscriptionEventsPayload struct {
	Subscription RazorpaySubscriptionWrapper         `json:"subscription"`
	Payment      *RazorpaySubscriptionPaymentWrapper `json:"payment,omitempty"`
}

// RazorpaySubscriptionWrapper wraps the main RazorpaySubscriptionEntity.
type RazorpaySubscriptionWrapper struct {
	Entity RazorpaySubscriptionEntity `json:"entity"`
}

// RazorpaySubscriptionPaymentWrapper wraps the main RazorpaySubscriptionPaymentEntity.
type RazorpaySubscriptionPaymentWrapper struct {
	Entity RazorpaySubscriptionPaymentEntity `json:"entity"`
}

type RazorpaySubscriptionCreateEvent struct {
	ID                  string                `json:"id"`
	Entity              string                `json:"entity"`
	PlanID              string                `json:"plan_id"`
	Status              string                `json:"status"`
	CurrentStart        *int64                `json:"current_start"`
	CurrentEnd          *int64                `json:"current_end"`
	EndedAt             *int64                `json:"ended_at"`
	Quantity            int                   `json:"quantity"`
	Notes               utils.JSONMap[string] `json:"notes"`
	ChargeAt            int64                 `json:"charge_at"`
	StartAt             int64                 `json:"start_at"`
	EndAt               int64                 `json:"end_at"`
	AuthAttempts        int                   `json:"auth_attempts"`
	TotalCount          int                   `json:"total_count"`
	PaidCount           int                   `json:"paid_count"`
	CustomerNotify      bool                  `json:"customer_notify"`
	CreatedAt           int64                 `json:"created_at"`
	ExpireBy            int64                 `json:"expire_by"`
	ShortURL            string                `json:"short_url"`
	HasScheduledChanges bool                  `json:"has_scheduled_changes"`
	ChangeScheduledAt   *int64                `json:"change_scheduled_at"`
	Source              string                `json:"source"`
	OfferID             string                `json:"offer_id"`
	RemainingCount      int                   `json:"remaining_count"`
}

// RazorpaySubscriptionEntity defines the structure of the subscription object.
type RazorpaySubscriptionEntity struct {
	ID                  string                   `json:"id"`
	Entity              string                   `json:"entity"`
	PlanID              string                   `json:"plan_id"`
	CustomerID          string                   `json:"customer_id"`
	Status              utils.SubscriptionStatus `json:"status"`
	Type                *int                     `json:"type,omitempty"`
	CurrentStart        int64                    `json:"current_start"`
	CurrentEnd          int64                    `json:"current_end"`
	EndedAt             *int64                   `json:"ended_at"`
	Quantity            int                      `json:"quantity"`
	Notes               any                      `json:"notes"` // Can be an object or an empty array
	ChargeAt            *int64                   `json:"charge_at"`
	StartAt             int64                    `json:"start_at"`
	EndAt               int64                    `json:"end_at"`
	AuthAttempts        int                      `json:"auth_attempts"`
	TotalCount          int                      `json:"total_count"`
	PaidCount           int                      `json:"paid_count"`
	CustomerNotify      bool                     `json:"customer_notify"`
	CreatedAt           int64                    `json:"created_at"`
	ExpireBy            *int64                   `json:"expire_by"`
	ShortURL            any                      `json:"short_url"`
	HasScheduledChanges bool                     `json:"has_scheduled_changes"`
	ChangeScheduledAt   any                      `json:"change_scheduled_at"`
	Source              string                   `json:"source"`
	OfferID             *string                  `json:"offer_id"`
	RemainingCount      int                      `json:"remaining_count"`
	PaymentMethod       *string                  `json:"payment_method,omitempty"`
	PauseInitiatedBy    *string                  `json:"pause_initiated_by,omitempty"`
	CancelInitiatedBy   *string                  `json:"cancel_initiated_by,omitempty"`
}

// RazorpaySubscriptionPaymentEntity defines the structure of the payment object within a subscription event.
// It is created as a separate struct because some fields (like 'captured') differ from the main RazorpayPaymentEntity.
type RazorpaySubscriptionPaymentEntity struct {
	ID                string               `json:"id"`
	Entity            string               `json:"entity"`
	Amount            int                  `json:"amount"`
	Currency          string               `json:"currency"`
	Status            string               `json:"status"`
	OrderID           string               `json:"order_id"`
	InvoiceID         string               `json:"invoice_id"`
	International     bool                 `json:"international"`
	Method            string               `json:"method"`
	AmountRefunded    int                  `json:"amount_refunded"`
	AmountTransferred int                  `json:"amount_transferred"`
	RefundStatus      any                  `json:"refund_status"`
	Captured          string               `json:"captured"` // Note: This is a string ("1") in subscription payments
	Description       string               `json:"description"`
	CardID            *string              `json:"card_id"`
	Card              *RazorpayCardDetails `json:"card,omitempty"`
	Bank              *string              `json:"bank"`
	Wallet            *string              `json:"wallet"`
	VPA               *string              `json:"vpa"`
	Email             string               `json:"email"`
	Contact           string               `json:"contact"`
	CustomerID        *string              `json:"customer_id,omitempty"`
	TokenID           any                  `json:"token_id"`
	Notes             []any                `json:"notes"`
	Fee               int                  `json:"fee"`
	Tax               int                  `json:"tax"`
	ErrorCode         any                  `json:"error_code"`
	ErrorDescription  any                  `json:"error_description"`
	CreatedAt         int64                `json:"created_at"`
}
