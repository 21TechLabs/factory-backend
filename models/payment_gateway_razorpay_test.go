package models

import (
	"errors"
	"log"
	"os"
	"testing"

	"github.com/21TechLabs/factory-backend/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockRazorpayClient mocks the Razorpay client for testing
type MockRazorpayClient struct {
	mock.Mock
}

type MockOrderAPI struct {
	mock.Mock
}

func (m *MockOrderAPI) Fetch(orderID string, extraParams map[string]interface{}, headers map[string]string) (map[string]interface{}, error) {
	args := m.Called(orderID, extraParams, headers)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockOrderAPI) Create(data map[string]interface{}, extraParams map[string]interface{}) (map[string]interface{}, error) {
	args := m.Called(data, extraParams)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// MockTransactionStore mocks the transaction store
type MockTransactionStore struct {
	mock.Mock
}

func (m *MockTransactionStore) GetByPaymentGatewayTransactionID(txnID string) (*Transaction, error) {
	args := m.Called(txnID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Transaction), args.Error(1)
}

func (m *MockTransactionStore) Update(txn *Transaction) error {
	args := m.Called(txn)
	return args.Error(0)
}

// MockUserStore mocks the user store
type MockUserStore struct {
	mock.Mock
}

func (m *MockUserStore) UserGetById(userID uint) (User, error) {
	args := m.Called(userID)
	return args.Get(0).(User), args.Error(1)
}

func (m *MockUserStore) Update(user *User) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestProcessFailedPayments(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.LstdFlags)

	t.Run("Success_UpdatesTransactionToFailed", func(t *testing.T) {
		// Setup
		mockTxnStore := new(MockTransactionStore)
		mockUserStore := new(MockUserStore)
		mockOrderAPI := new(MockOrderAPI)

		rpg := &RazorpayPG{
			Logger:           logger,
			TransactionStore: mockTxnStore,
			UserStore:        mockUserStore,
		}

		orderID := "order_123"
		paymentID := "pay_456"
		
		event := RazorpayBaseEvent[RazorpayPaymentFailedPayload]{
			Entity:    "event",
			AccountID: "acc_123",
			Event:     "payment.failed",
			Payload: RazorpayPaymentFailedPayload{
				Payment: RazorpayPaymentFailedWrapper{
					Entity: RazorpayPaymentFailedEntity{
						ID:               paymentID,
						OrderID:          orderID,
						Status:           "failed",
						Amount:           10000,
						Currency:         "INR",
						ErrorCode:        "BAD_REQUEST_ERROR",
						ErrorDescription: "Payment failed",
					},
				},
			},
		}

		// Mock order fetch - return paid status
		// Note: The actual implementation has a bug - it checks order status is "paid"
		// for a failed payment, which doesn't make logical sense
		mockOrderAPI.On("Fetch", orderID, mock.Anything, mock.Anything).Return(
			map[string]interface{}{
				"id":     orderID,
				"status": "paid",
			},
			nil,
		)

		// Mock transaction retrieval
		existingTxn := &Transaction{
			ID:                          1,
			UserID:                      100,
			Status:                      utils.TransactionStatusPending,
			Token:                       1000,
			Amount:                      10000,
			PaymentGatewayTransactionID: orderID,
		}
		mockTxnStore.On("GetByPaymentGatewayTransactionID", paymentID).Return(existingTxn, nil)

		// Mock transaction update
		mockTxnStore.On("Update", mock.MatchedBy(func(txn *Transaction) bool {
			return txn.Status == utils.TransactionStatusFailed &&
				txn.PaymentGatewayRedirectURL == "" &&
				txn.ID == existingTxn.ID
		})).Return(nil)

		// Create a custom client structure to use our mock
		type mockClient struct {
			Order *MockOrderAPI
		}
		rpg.Client = &mockClient{Order: mockOrderAPI}

		// Execute
		resultTxn, err := rpg.ProcessFailedPayments(event)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, resultTxn)
		assert.Equal(t, utils.TransactionStatusFailed, resultTxn.Status)
		assert.Equal(t, "", resultTxn.PaymentGatewayRedirectURL)
		assert.Equal(t, uint(1), resultTxn.ID)

		mockOrderAPI.AssertExpectations(t)
		mockTxnStore.AssertExpectations(t)
	})

	t.Run("Error_EmptyOrderID", func(t *testing.T) {
		mockTxnStore := new(MockTransactionStore)
		mockUserStore := new(MockUserStore)

		rpg := &RazorpayPG{
			Logger:           logger,
			TransactionStore: mockTxnStore,
			UserStore:        mockUserStore,
		}

		event := RazorpayBaseEvent[RazorpayPaymentFailedPayload]{
			Payload: RazorpayPaymentFailedPayload{
				Payment: RazorpayPaymentFailedWrapper{
					Entity: RazorpayPaymentFailedEntity{
						ID:      "pay_123",
						OrderID: "", // Empty order ID
					},
				},
			},
		}

		// Execute
		resultTxn, err := rpg.ProcessFailedPayments(event)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resultTxn)
		assert.Equal(t, utils.ErrInvalidOrderID, err)
	})

	t.Run("Error_OrderFetchFails", func(t *testing.T) {
		mockTxnStore := new(MockTransactionStore)
		mockUserStore := new(MockUserStore)
		mockOrderAPI := new(MockOrderAPI)

		rpg := &RazorpayPG{
			Logger:           logger,
			TransactionStore: mockTxnStore,
			UserStore:        mockUserStore,
		}

		orderID := "order_123"
		event := RazorpayBaseEvent[RazorpayPaymentFailedPayload]{
			Payload: RazorpayPaymentFailedPayload{
				Payment: RazorpayPaymentFailedWrapper{
					Entity: RazorpayPaymentFailedEntity{
						ID:      "pay_456",
						OrderID: orderID,
					},
				},
			},
		}

		fetchError := errors.New("network error")
		mockOrderAPI.On("Fetch", orderID, mock.Anything, mock.Anything).Return(
			nil,
			fetchError,
		)

		type mockClient struct {
			Order *MockOrderAPI
		}
		rpg.Client = &mockClient{Order: mockOrderAPI}

		// Execute
		resultTxn, err := rpg.ProcessFailedPayments(event)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resultTxn)
		assert.Contains(t, err.Error(), "failed to fetch order")

		mockOrderAPI.AssertExpectations(t)
	})

	t.Run("Error_OrderNotPaid", func(t *testing.T) {
		mockTxnStore := new(MockTransactionStore)
		mockUserStore := new(MockUserStore)
		mockOrderAPI := new(MockOrderAPI)

		rpg := &RazorpayPG{
			Logger:           logger,
			TransactionStore: mockTxnStore,
			UserStore:        mockUserStore,
		}

		orderID := "order_123"
		event := RazorpayBaseEvent[RazorpayPaymentFailedPayload]{
			Payload: RazorpayPaymentFailedPayload{
				Payment: RazorpayPaymentFailedWrapper{
					Entity: RazorpayPaymentFailedEntity{
						ID:      "pay_456",
						OrderID: orderID,
					},
				},
			},
		}

		mockOrderAPI.On("Fetch", orderID, mock.Anything, mock.Anything).Return(
			map[string]interface{}{
				"id":     orderID,
				"status": "created", // Not paid
			},
			nil,
		)

		type mockClient struct {
			Order *MockOrderAPI
		}
		rpg.Client = &mockClient{Order: mockOrderAPI}

		// Execute
		resultTxn, err := rpg.ProcessFailedPayments(event)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resultTxn)
		assert.Contains(t, err.Error(), "order is not paid")

		mockOrderAPI.AssertExpectations(t)
	})

	t.Run("Error_TransactionNotFound", func(t *testing.T) {
		mockTxnStore := new(MockTransactionStore)
		mockUserStore := new(MockUserStore)
		mockOrderAPI := new(MockOrderAPI)

		rpg := &RazorpayPG{
			Logger:           logger,
			TransactionStore: mockTxnStore,
			UserStore:        mockUserStore,
		}

		orderID := "order_123"
		paymentID := "pay_456"
		event := RazorpayBaseEvent[RazorpayPaymentFailedPayload]{
			Payload: RazorpayPaymentFailedPayload{
				Payment: RazorpayPaymentFailedWrapper{
					Entity: RazorpayPaymentFailedEntity{
						ID:      paymentID,
						OrderID: orderID,
					},
				},
			},
		}

		mockOrderAPI.On("Fetch", orderID, mock.Anything, mock.Anything).Return(
			map[string]interface{}{
				"id":     orderID,
				"status": "paid",
			},
			nil,
		)

		mockTxnStore.On("GetByPaymentGatewayTransactionID", paymentID).Return(
			nil,
			nil, // No error, but nil transaction
		)

		type mockClient struct {
			Order *MockOrderAPI
		}
		rpg.Client = &mockClient{Order: mockOrderAPI}

		// Execute
		resultTxn, err := rpg.ProcessFailedPayments(event)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resultTxn)
		assert.Equal(t, utils.ErrTransactionNotFound, err)

		mockOrderAPI.AssertExpectations(t)
		mockTxnStore.AssertExpectations(t)
	})

	t.Run("Error_TransactionRetrievalFails", func(t *testing.T) {
		mockTxnStore := new(MockTransactionStore)
		mockUserStore := new(MockUserStore)
		mockOrderAPI := new(MockOrderAPI)

		rpg := &RazorpayPG{
			Logger:           logger,
			TransactionStore: mockTxnStore,
			UserStore:        mockUserStore,
		}

		orderID := "order_123"
		paymentID := "pay_456"
		event := RazorpayBaseEvent[RazorpayPaymentFailedPayload]{
			Payload: RazorpayPaymentFailedPayload{
				Payment: RazorpayPaymentFailedWrapper{
					Entity: RazorpayPaymentFailedEntity{
						ID:      paymentID,
						OrderID: orderID,
					},
				},
			},
		}

		mockOrderAPI.On("Fetch", orderID, mock.Anything, mock.Anything).Return(
			map[string]interface{}{
				"id":     orderID,
				"status": "paid",
			},
			nil,
		)

		dbError := errors.New("database connection error")
		mockTxnStore.On("GetByPaymentGatewayTransactionID", paymentID).Return(
			nil,
			dbError,
		)

		type mockClient struct {
			Order *MockOrderAPI
		}
		rpg.Client = &mockClient{Order: mockOrderAPI}

		// Execute
		resultTxn, err := rpg.ProcessFailedPayments(event)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resultTxn)
		assert.Contains(t, err.Error(), "failed to get transaction by order ID")

		mockOrderAPI.AssertExpectations(t)
		mockTxnStore.AssertExpectations(t)
	})

	t.Run("Success_AlreadyCompleted_IgnoresEvent", func(t *testing.T) {
		mockTxnStore := new(MockTransactionStore)
		mockUserStore := new(MockUserStore)
		mockOrderAPI := new(MockOrderAPI)

		rpg := &RazorpayPG{
			Logger:           logger,
			TransactionStore: mockTxnStore,
			UserStore:        mockUserStore,
		}

		orderID := "order_123"
		paymentID := "pay_456"
		event := RazorpayBaseEvent[RazorpayPaymentFailedPayload]{
			Payload: RazorpayPaymentFailedPayload{
				Payment: RazorpayPaymentFailedWrapper{
					Entity: RazorpayPaymentFailedEntity{
						ID:      paymentID,
						OrderID: orderID,
					},
				},
			},
		}

		mockOrderAPI.On("Fetch", orderID, mock.Anything, mock.Anything).Return(
			map[string]interface{}{
				"id":     orderID,
				"status": "paid",
			},
			nil,
		)

		// Transaction already completed
		existingTxn := &Transaction{
			ID:                          1,
			UserID:                      100,
			Status:                      utils.TransactionStatusCompleted, // Already completed
			Token:                       1000,
			Amount:                      10000,
			PaymentGatewayTransactionID: orderID,
		}
		mockTxnStore.On("GetByPaymentGatewayTransactionID", paymentID).Return(existingTxn, nil)

		// Update should NOT be called
		// mockTxnStore should not expect Update call

		type mockClient struct {
			Order *MockOrderAPI
		}
		rpg.Client = &mockClient{Order: mockOrderAPI}

		// Execute
		resultTxn, err := rpg.ProcessFailedPayments(event)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, resultTxn)
		assert.Equal(t, utils.TransactionStatusCompleted, resultTxn.Status) // Still completed
		assert.Equal(t, uint(1), resultTxn.ID)

		mockOrderAPI.AssertExpectations(t)
		mockTxnStore.AssertExpectations(t)
		// Verify Update was NOT called
		mockTxnStore.AssertNotCalled(t, "Update")
	})

	t.Run("Error_TransactionUpdateFails", func(t *testing.T) {
		mockTxnStore := new(MockTransactionStore)
		mockUserStore := new(MockUserStore)
		mockOrderAPI := new(MockOrderAPI)

		rpg := &RazorpayPG{
			Logger:           logger,
			TransactionStore: mockTxnStore,
			UserStore:        mockUserStore,
		}

		orderID := "order_123"
		paymentID := "pay_456"
		event := RazorpayBaseEvent[RazorpayPaymentFailedPayload]{
			Payload: RazorpayPaymentFailedPayload{
				Payment: RazorpayPaymentFailedWrapper{
					Entity: RazorpayPaymentFailedEntity{
						ID:      paymentID,
						OrderID: orderID,
					},
				},
			},
		}

		mockOrderAPI.On("Fetch", orderID, mock.Anything, mock.Anything).Return(
			map[string]interface{}{
				"id":     orderID,
				"status": "paid",
			},
			nil,
		)

		existingTxn := &Transaction{
			ID:                          1,
			UserID:                      100,
			Status:                      utils.TransactionStatusPending,
			Token:                       1000,
			Amount:                      10000,
			PaymentGatewayTransactionID: orderID,
		}
		mockTxnStore.On("GetByPaymentGatewayTransactionID", paymentID).Return(existingTxn, nil)

		updateError := errors.New("database write error")
		mockTxnStore.On("Update", mock.Anything).Return(updateError)

		type mockClient struct {
			Order *MockOrderAPI
		}
		rpg.Client = &mockClient{Order: mockOrderAPI}

		// Execute
		resultTxn, err := rpg.ProcessFailedPayments(event)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resultTxn)
		assert.Contains(t, err.Error(), "failed to update transaction")

		mockOrderAPI.AssertExpectations(t)
		mockTxnStore.AssertExpectations(t)
	})

	t.Run("Success_WithDifferentErrorCodes", func(t *testing.T) {
		errorCodes := []string{
			"BAD_REQUEST_ERROR",
			"GATEWAY_ERROR",
			"SERVER_ERROR",
			"INSUFFICIENT_FUNDS",
		}

		for _, errorCode := range errorCodes {
			t.Run("ErrorCode_"+errorCode, func(t *testing.T) {
				mockTxnStore := new(MockTransactionStore)
				mockUserStore := new(MockUserStore)
				mockOrderAPI := new(MockOrderAPI)

				rpg := &RazorpayPG{
					Logger:           logger,
					TransactionStore: mockTxnStore,
					UserStore:        mockUserStore,
				}

				orderID := "order_" + errorCode
				paymentID := "pay_" + errorCode
				
				event := RazorpayBaseEvent[RazorpayPaymentFailedPayload]{
					Payload: RazorpayPaymentFailedPayload{
						Payment: RazorpayPaymentFailedWrapper{
							Entity: RazorpayPaymentFailedEntity{
								ID:               paymentID,
								OrderID:          orderID,
								ErrorCode:        errorCode,
								ErrorDescription: "Payment failed with " + errorCode,
							},
						},
					},
				}

				mockOrderAPI.On("Fetch", orderID, mock.Anything, mock.Anything).Return(
					map[string]interface{}{
						"id":     orderID,
						"status": "paid",
					},
					nil,
				)

				existingTxn := &Transaction{
					ID:                          1,
					Status:                      utils.TransactionStatusPending,
					PaymentGatewayTransactionID: orderID,
				}
				mockTxnStore.On("GetByPaymentGatewayTransactionID", paymentID).Return(existingTxn, nil)
				mockTxnStore.On("Update", mock.Anything).Return(nil)

				type mockClient struct {
					Order *MockOrderAPI
				}
				rpg.Client = &mockClient{Order: mockOrderAPI}

				// Execute
				resultTxn, err := rpg.ProcessFailedPayments(event)

				// Assert
				assert.NoError(t, err)
				assert.NotNil(t, resultTxn)
				assert.Equal(t, utils.TransactionStatusFailed, resultTxn.Status)

				mockOrderAPI.AssertExpectations(t)
				mockTxnStore.AssertExpectations(t)
			})
		}
	})
}

// TestRazorpayCreateOrder_ToMap tests the ToMap conversion method
func TestRazorpayCreateOrder_ToMap(t *testing.T) {
	t.Run("CompleteOrder", func(t *testing.T) {
		order := RazorpayCreateOrder{
			Amount:         10000.50,
			Currency:       utils.CurrencyINR,
			Receipt:        "receipt_123",
			PartialPayment: false,
			Notes: map[string]string{
				"user_id": "100",
				"plan":    "premium",
			},
		}

		result := order.ToMap()

		assert.NotNil(t, result)
		assert.Equal(t, 10000.50, result["amount"])
		assert.Equal(t, utils.CurrencyINR, result["currency"])
		assert.Equal(t, "receipt_123", result["receipt"])
		assert.Equal(t, false, result["partial_payment"])
		assert.NotNil(t, result["notes"])
		notes := result["notes"].(map[string]string)
		assert.Equal(t, "100", notes["user_id"])
		assert.Equal(t, "premium", notes["plan"])
	})

	t.Run("MinimalOrder", func(t *testing.T) {
		order := RazorpayCreateOrder{
			Amount:   5000.0,
			Currency: utils.CurrencyUSD,
		}

		result := order.ToMap()

		assert.NotNil(t, result)
		assert.Equal(t, 5000.0, result["amount"])
		assert.Equal(t, utils.CurrencyUSD, result["currency"])
		assert.Equal(t, "", result["receipt"])
		assert.Equal(t, false, result["partial_payment"])
		assert.Nil(t, result["notes"])
	})

	t.Run("WithPartialPayment", func(t *testing.T) {
		order := RazorpayCreateOrder{
			Amount:         15000.0,
			Currency:       utils.CurrencyEUR,
			PartialPayment: true,
		}

		result := order.ToMap()

		assert.NotNil(t, result)
		assert.Equal(t, true, result["partial_payment"])
	})
}

// TestNewRazorpayPG tests the constructor
func TestNewRazorpayPG(t *testing.T) {
	t.Run("CreatesValidInstance", func(t *testing.T) {
		logger := log.New(os.Stdout, "TEST: ", log.LstdFlags)
		mockTxnStore := new(MockTransactionStore)
		mockUserStore := new(MockUserStore)

		rpg := NewRazorpayPG(logger, mockTxnStore, mockUserStore)

		assert.NotNil(t, rpg)
		assert.Equal(t, logger, rpg.Logger)
		assert.Equal(t, mockTxnStore, rpg.TransactionStore)
		assert.Equal(t, mockUserStore, rpg.UserStore)
		assert.NotNil(t, rpg.Client) // RazorpayClient is set globally
	})
}