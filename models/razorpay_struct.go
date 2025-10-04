package models

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
