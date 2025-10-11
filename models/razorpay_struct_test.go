package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRazorpayBaseEvent_Deserialization tests JSON unmarshaling
func TestRazorpayBaseEvent_Deserialization(t *testing.T) {
	t.Run("OrderPaidEvent_ValidJSON", func(t *testing.T) {
		jsonData := `{
			"entity": "event",
			"account_id": "acc_ABC123",
			"event": "order.paid",
			"contains": ["payment", "order"],
			"created_at": 1609459200,
			"payload": {
				"payment": {
					"entity": {
						"id": "pay_XYZ789",
						"entity": "payment",
						"amount": 50000,
						"currency": "INR",
						"status": "captured",
						"order_id": "order_ABC123",
						"invoice_id": null,
						"international": false,
						"method": "card",
						"amount_refunded": 0,
						"refund_status": null,
						"captured": true,
						"description": null,
						"card_id": "card_123",
						"bank": null,
						"wallet": null,
						"vpa": null,
						"email": "test@example.com",
						"contact": "+919876543210",
						"notes": [],
						"fee": 1000,
						"tax": 180,
						"error_code": null,
						"error_description": null,
						"created_at": 1609459200
					}
				},
				"order": {
					"entity": {
						"id": "order_ABC123",
						"entity": "order",
						"amount": 50000,
						"amount_paid": 50000,
						"amount_due": 0,
						"currency": "INR",
						"receipt": "receipt_123",
						"offer_id": null,
						"status": "paid",
						"attempts": 1,
						"notes": [],
						"created_at": 1609459100
					}
				}
			}
		}`

		var event RazorpayBaseEvent[RazorpayOrderPaidPayload]
		err := json.Unmarshal([]byte(jsonData), &event)

		require.NoError(t, err)
		assert.Equal(t, "event", event.Entity)
		assert.Equal(t, "acc_ABC123", event.AccountID)
		assert.Equal(t, "order.paid", event.Event)
		assert.Equal(t, int64(1609459200), event.CreatedAt)
		assert.Equal(t, "pay_XYZ789", event.Payload.Payment.Entity.ID)
		assert.Equal(t, "order_ABC123", event.Payload.Order.Entity.ID)
		assert.Equal(t, OrderStatusPaid, event.Payload.Order.Entity.Status)
		assert.Equal(t, 50000, event.Payload.Payment.Entity.Amount)
	})

	t.Run("PaymentFailedEvent_ValidJSON", func(t *testing.T) {
		jsonData := `{
			"entity": "event",
			"account_id": "acc_ABC123",
			"event": "payment.failed",
			"contains": ["payment"],
			"created_at": 1609459300,
			"payload": {
				"payment": {
					"entity": {
						"id": "pay_FAILED123",
						"entity": "payment",
						"amount": 50000,
						"currency": "INR",
						"status": "failed",
						"order_id": "order_XYZ456",
						"invoice_id": null,
						"international": false,
						"method": "card",
						"amount_refunded": 0,
						"refund_status": null,
						"captured": false,
						"description": null,
						"card_id": "card_789",
						"bank": null,
						"wallet": null,
						"vpa": null,
						"email": "user@example.com",
						"contact": "+919876543210",
						"notes": [],
						"fee": null,
						"tax": null,
						"error_code": "BAD_REQUEST_ERROR",
						"error_description": "Payment processing failed",
						"error_source": "customer",
						"error_step": "payment_authorization",
						"error_reason": "payment_failed",
						"acquirer_data": {
							"bank_transaction_id": null,
							"auth_code": null,
							"rrn": null,
							"transaction_id": null
						},
						"created_at": 1609459300
					}
				}
			}
		}`

		var event RazorpayBaseEvent[RazorpayPaymentFailedPayload]
		err := json.Unmarshal([]byte(jsonData), &event)

		require.NoError(t, err)
		assert.Equal(t, "event", event.Entity)
		assert.Equal(t, "payment.failed", event.Event)
		assert.Equal(t, "pay_FAILED123", event.Payload.Payment.Entity.ID)
		assert.Equal(t, "order_XYZ456", event.Payload.Payment.Entity.OrderID)
		assert.Equal(t, "failed", event.Payload.Payment.Entity.Status)
		assert.Equal(t, "BAD_REQUEST_ERROR", event.Payload.Payment.Entity.ErrorCode)
		assert.Equal(t, "Payment processing failed", event.Payload.Payment.Entity.ErrorDescription)
		assert.NotNil(t, event.Payload.Payment.Entity.ErrorSource)
		assert.Equal(t, "customer", *event.Payload.Payment.Entity.ErrorSource)
	})

	t.Run("PaymentFailedEvent_WithCardDetails", func(t *testing.T) {
		jsonData := `{
			"entity": "event",
			"account_id": "acc_123",
			"event": "payment.failed",
			"contains": ["payment"],
			"created_at": 1609459400,
			"payload": {
				"payment": {
					"entity": {
						"id": "pay_CARD_FAILED",
						"entity": "payment",
						"amount": 25000,
						"currency": "INR",
						"status": "failed",
						"order_id": "order_CARD_123",
						"invoice_id": null,
						"international": false,
						"method": "card",
						"amount_refunded": 0,
						"refund_status": null,
						"captured": false,
						"description": "Premium plan",
						"card_id": "card_456",
						"bank": "HDFC",
						"wallet": null,
						"vpa": null,
						"email": "cardholder@example.com",
						"contact": "+919999999999",
						"notes": [],
						"fee": null,
						"tax": null,
						"error_code": "GATEWAY_ERROR",
						"error_description": "Card issuer declined",
						"error_source": "bank",
						"error_step": "payment_authorization",
						"error_reason": "card_declined",
						"acquirer_data": {
							"bank_transaction_id": "BNK123456",
							"auth_code": "AUTH123",
							"rrn": "RRN123456789",
							"transaction_id": "TXN987654"
						},
						"token_id": "token_123",
						"created_at": 1609459400,
						"card": {
							"id": "card_456",
							"entity": "card",
							"name": "Test User",
							"last4": "1234",
							"network": "Visa",
							"type": "credit",
							"issuer": "HDFC Bank",
							"international": false,
							"emi": true
						}
					}
				}
			}
		}`

		var event RazorpayBaseEvent[RazorpayPaymentFailedPayload]
		err := json.Unmarshal([]byte(jsonData), &event)

		require.NoError(t, err)
		assert.Equal(t, "pay_CARD_FAILED", event.Payload.Payment.Entity.ID)
		assert.Equal(t, "GATEWAY_ERROR", event.Payload.Payment.Entity.ErrorCode)
		
		// Check card details
		assert.NotNil(t, event.Payload.Payment.Entity.Card)
		assert.Equal(t, "card_456", event.Payload.Payment.Entity.Card.ID)
		assert.Equal(t, "1234", event.Payload.Payment.Entity.Card.Last4)
		assert.Equal(t, "Visa", event.Payload.Payment.Entity.Card.Network)
		
		// Check acquirer data
		assert.NotNil(t, event.Payload.Payment.Entity.AcquirerData.BankTransactionID)
		assert.Equal(t, "BNK123456", *event.Payload.Payment.Entity.AcquirerData.BankTransactionID)
	})

	t.Run("PaymentFailedEvent_WithUPIDetails", func(t *testing.T) {
		jsonData := `{
			"entity": "event",
			"account_id": "acc_UPI",
			"event": "payment.failed",
			"contains": ["payment"],
			"created_at": 1609459500,
			"payload": {
				"payment": {
					"entity": {
						"id": "pay_UPI_FAILED",
						"entity": "payment",
						"amount": 15000,
						"currency": "INR",
						"status": "failed",
						"order_id": "order_UPI_789",
						"invoice_id": null,
						"international": false,
						"method": "upi",
						"amount_refunded": 0,
						"refund_status": null,
						"captured": false,
						"description": null,
						"card_id": null,
						"bank": null,
						"wallet": null,
						"vpa": "user@paytm",
						"email": "upiuser@example.com",
						"contact": "+918888888888",
						"notes": [],
						"fee": null,
						"tax": null,
						"error_code": "BAD_REQUEST_ERROR",
						"error_description": "UPI transaction failed",
						"error_source": "customer",
						"error_step": "payment_initiation",
						"error_reason": "customer_cancelled",
						"acquirer_data": {},
						"created_at": 1609459500,
						"upi": {
							"payer_account_type": "savings",
							"vpa": "user@paytm",
							"flow": "collect"
						}
					}
				}
			}
		}`

		var event RazorpayBaseEvent[RazorpayPaymentFailedPayload]
		err := json.Unmarshal([]byte(jsonData), &event)

		require.NoError(t, err)
		assert.Equal(t, "pay_UPI_FAILED", event.Payload.Payment.Entity.ID)
		assert.Equal(t, "upi", event.Payload.Payment.Entity.Method)
		
		// Check UPI details
		assert.NotNil(t, event.Payload.Payment.Entity.UPI)
		assert.Equal(t, "user@paytm", event.Payload.Payment.Entity.UPI.VPA)
		assert.Equal(t, "savings", event.Payload.Payment.Entity.UPI.PayerAccountType)
		assert.Equal(t, "collect", event.Payload.Payment.Entity.UPI.Flow)
	})
}

// TestRazorpayPaymentFailedEntity_EdgeCases tests edge cases
func TestRazorpayPaymentFailedEntity_EdgeCases(t *testing.T) {
	t.Run("NullableFields_AllNull", func(t *testing.T) {
		jsonData := `{
			"id": "pay_NULL_TEST",
			"entity": "payment",
			"amount": 10000,
			"currency": "INR",
			"status": "failed",
			"order_id": "order_NULL",
			"invoice_id": null,
			"international": false,
			"method": "netbanking",
			"amount_refunded": 0,
			"refund_status": null,
			"captured": false,
			"description": null,
			"card_id": null,
			"bank": null,
			"wallet": null,
			"vpa": null,
			"email": "null@test.com",
			"contact": "+910000000000",
			"notes": [],
			"fee": null,
			"tax": null,
			"error_code": "SERVER_ERROR",
			"error_description": "Internal server error",
			"error_source": null,
			"error_step": null,
			"error_reason": null,
			"acquirer_data": {},
			"created_at": 1609459600
		}`

		var entity RazorpayPaymentFailedEntity
		err := json.Unmarshal([]byte(jsonData), &entity)

		require.NoError(t, err)
		assert.Equal(t, "pay_NULL_TEST", entity.ID)
		assert.Nil(t, entity.CardID)
		assert.Nil(t, entity.Bank)
		assert.Nil(t, entity.Wallet)
		assert.Nil(t, entity.VPA)
		assert.Nil(t, entity.ErrorSource)
		assert.Nil(t, entity.ErrorStep)
		assert.Nil(t, entity.ErrorReason)
		assert.Nil(t, entity.TokenID)
	})

	t.Run("EmptyAcquirerData", func(t *testing.T) {
		jsonData := `{
			"bank_transaction_id": null,
			"auth_code": null,
			"rrn": null,
			"transaction_id": null
		}`

		var acquirerData RazorpayAcquirerData
		err := json.Unmarshal([]byte(jsonData), &acquirerData)

		require.NoError(t, err)
		assert.Nil(t, acquirerData.BankTransactionID)
		assert.Nil(t, acquirerData.AuthCode)
		assert.Nil(t, acquirerData.RRN)
		assert.Nil(t, acquirerData.TransactionID)
	})

	t.Run("MinimalPaymentFailedEvent", func(t *testing.T) {
		jsonData := `{
			"entity": "event",
			"account_id": "acc_MIN",
			"event": "payment.failed",
			"contains": ["payment"],
			"created_at": 1609459700,
			"payload": {
				"payment": {
					"entity": {
						"id": "pay_MIN",
						"entity": "payment",
						"amount": 100,
						"currency": "INR",
						"status": "failed",
						"order_id": "order_MIN",
						"invoice_id": null,
						"international": false,
						"method": "netbanking",
						"amount_refunded": 0,
						"refund_status": null,
						"captured": false,
						"description": null,
						"card_id": null,
						"bank": null,
						"wallet": null,
						"vpa": null,
						"email": "",
						"contact": "",
						"notes": [],
						"fee": null,
						"tax": null,
						"error_code": "",
						"error_description": "",
						"error_source": null,
						"error_step": null,
						"error_reason": null,
						"acquirer_data": {},
						"created_at": 1609459700
					}
				}
			}
		}`

		var event RazorpayBaseEvent[RazorpayPaymentFailedPayload]
		err := json.Unmarshal([]byte(jsonData), &event)

		require.NoError(t, err)
		assert.Equal(t, "pay_MIN", event.Payload.Payment.Entity.ID)
		assert.Equal(t, "", event.Payload.Payment.Entity.Email)
		assert.Equal(t, "", event.Payload.Payment.Entity.ErrorCode)
	})
}

// TestOrderStatus_Constants tests order status constants
func TestOrderStatus_Constants(t *testing.T) {
	t.Run("OrderStatusPaid", func(t *testing.T) {
		assert.Equal(t, OrderStatus("paid"), OrderStatusPaid)
	})

	t.Run("OrderStatusFailed", func(t *testing.T) {
		assert.Equal(t, OrderStatus("failed"), OrderStatusFailed)
	})

	t.Run("StatusComparison", func(t *testing.T) {
		assert.NotEqual(t, OrderStatusPaid, OrderStatusFailed)
	})
}

// BenchmarkRazorpayEventDeserialization benchmarks JSON unmarshaling
func BenchmarkRazorpayEventDeserialization(b *testing.B) {
	jsonData := []byte(`{
		"entity": "event",
		"account_id": "acc_123",
		"event": "payment.failed",
		"contains": ["payment"],
		"created_at": 1609459200,
		"payload": {
			"payment": {
				"entity": {
					"id": "pay_123",
					"entity": "payment",
					"amount": 50000,
					"currency": "INR",
					"status": "failed",
					"order_id": "order_456",
					"invoice_id": null,
					"international": false,
					"method": "card",
					"amount_refunded": 0,
					"refund_status": null,
					"captured": false,
					"description": null,
					"card_id": null,
					"bank": null,
					"wallet": null,
					"vpa": null,
					"email": "test@example.com",
					"contact": "+919876543210",
					"notes": [],
					"fee": null,
					"tax": null,
					"error_code": "ERROR",
					"error_description": "Failed",
					"error_source": null,
					"error_step": null,
					"error_reason": null,
					"acquirer_data": {},
					"created_at": 1609459200
				}
			}
		}
	}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var event RazorpayBaseEvent[RazorpayPaymentFailedPayload]
		_ = json.Unmarshal(jsonData, &event)
	}
}