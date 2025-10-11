package utils

import "errors"

type PaymentGatewayError struct {
	Message string `json:"message"`
}

var (
	ErrInvalidPlanType        = errors.New("invalid plan type")
	ErrSubscriptionNotFound   = errors.New("subscription not found")
	ErrInvalidSubscription    = errors.New("subscription not found")
	ErrPaymentGatewayNotFound = errors.New("payment gateway not found")
	ErrTransactionNotFound    = errors.New("transaction not found")
	ErrInvalidOrderID         = errors.New("invalid order ID")
	ErrInvalidLimit           = errors.New("invalid limit")
)

func (e *PaymentGatewayError) Error() string {
	return e.Message
}
