package models

import (
	"log"

	"github.com/21TechLabs/factory-backend/utils"
)

type PaymentGateway string

const (
	PaymentGatewayRazorpay = "razorpay"
)

type PaymentGatewayInterface interface {
	InitiatePayment(*ProductPlan, *User, int) (*Transaction, error)
	CaptureOrderPaid(RazorpayBaseEvent[RazorpayOrderPaidPayload]) (*Transaction, error)
	ProcessFailedPayments(RazorpayBaseEvent[RazorpayPaymentFailedPayload]) (*Transaction, error)
}

func GetPaymentGateway(gateway string, logger *log.Logger, transactionStore *TransactionStore, userStore *UserStore) (PaymentGatewayInterface, error) {
	switch gateway {
	case PaymentGatewayRazorpay:
		return NewRazorpayPG(logger, transactionStore, userStore), nil
	default:
		return nil, utils.ErrPaymentGatewayNotFound
	}
}
