package models

import (
	"fmt"
	"log"

	"github.com/21TechLabs/factory-backend/dto"
	"github.com/21TechLabs/factory-backend/utils"
	"github.com/razorpay/razorpay-go"
)

var RazorpayClient *razorpay.Client

func init() {
	var razorpayClientApiKey = utils.GetEnv("RAZORPAY_KEY_ID", false)
	var razorpayClientApiSecret = utils.GetEnv("RAZORPAY_KEY_SECRET", false)
	RazorpayClient = razorpay.NewClient(razorpayClientApiKey, razorpayClientApiSecret)
}

type RazorpayPG struct {
	Logger           *log.Logger
	TransactionStore *TransactionStore
	UserStore        *UserStore
	Client           *razorpay.Client
}

type RazorpayCreateOrder struct {
	Amount         float64           `json:"amount"`
	Currency       utils.Currency    `json:"currency"`
	Receipt        string            `json:"receipt,omitempty"`
	PartialPayment bool              `json:"partial_payment,omitempty"`
	Notes          map[string]string `json:"notes,omitempty"`
}

func (rco *RazorpayCreateOrder) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"amount":          rco.Amount,
		"currency":        rco.Currency,
		"receipt":         rco.Receipt,
		"partial_payment": rco.PartialPayment,
		"notes":           rco.Notes,
	}
}

func NewRazorpayPG(log *log.Logger, transactionStore *TransactionStore, userStore *UserStore) *RazorpayPG {
	return &RazorpayPG{
		Logger:           log,
		TransactionStore: transactionStore,
		Client:           RazorpayClient,
		UserStore:        userStore,
	}
}

func (rpg *RazorpayPG) InitiatePayment(productPlan *ProductPlan, user *User) (*Transaction, error) {
	switch productPlan.PlanType {
	case utils.PlanTypeOneTime:
		return rpg.initiatePaymentOneTime(productPlan, user)
	case utils.PlanTypeSubscription:
		return rpg.initiatePaymentSubscription(productPlan, user)
	default:
		return nil, utils.ErrInvalidPlanType
	}

}

func (rpg *RazorpayPG) initiatePaymentOneTime(productPlan *ProductPlan, user *User) (*Transaction, error) {
	// create a transaction
	txn := &dto.TransactionCreateDto{
		Token:                       productPlan.Tokens,
		Amount:                      productPlan.PlanPrice * 100, // Convert to smallest currency unit
		Currency:                    productPlan.PlanCurrency,
		Status:                      utils.TransactionStatusPending,
		ProductPlanID:               &(productPlan).ID,
		PaymentGatewayName:          PaymentGatewayRazorpay,
		PaymentGatewayRedirectURL:   "",
		PaymentGatewayTransactionID: "",
	}

	transaction, err := rpg.TransactionStore.CreateTransaction(txn, user)
	if err != nil {
		return nil, err
	}

	order := RazorpayCreateOrder{
		Amount:         productPlan.PlanPrice, // Convert to smallest currency unit
		Currency:       productPlan.PlanCurrency,
		Receipt:        transaction.ReceiptId,
		PartialPayment: false,
		Notes:          map[string]string{"user_id": fmt.Sprintf("%d", user.ID)},
	}

	orderData, err := rpg.Client.Order.Create(order.ToMap(), nil)

	if err != nil {
		return nil, err
	}

	id, ok := orderData["id"].(string)
	if !ok {
		return nil, utils.ErrInvalidOrderID
	}
	transaction.PaymentGatewayTransactionID = id

	// Update the transaction with the Razorpay order details
	if err := rpg.TransactionStore.Update(&transaction); err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (rpg *RazorpayPG) initiatePaymentSubscription(productPlan *ProductPlan, user *User) (*Transaction, error) {
	return nil, nil
}

func (rpg *RazorpayPG) CaptureOrderPaid(event RazorpayBaseEvent[RazorpayOrderPaidPayload]) (*Transaction, error) {
	orderId := event.Payload.Order.Entity.ID

	if orderId == "" {
		return nil, utils.ErrInvalidOrderID
	}

	order, err := rpg.Client.Order.Fetch(orderId, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order: %w", err)
	}

	if order["status"] != OrderStatusPaid {
		return nil, fmt.Errorf("order is not paid, status: %s", order["status"])
	}

	txn, err := rpg.TransactionStore.GetByPaymentGatewayTransactionID(event.Payload.Order.Entity.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by order ID: %w", err)
	}

	if txn == nil {
		return nil, utils.ErrTransactionNotFound
	}

	txn.Status = utils.TransactionStatusCompleted
	txn.PaymentGatewayRedirectURL = ""
	txn.PaymentGatewayTransactionID = orderId

	if err := rpg.TransactionStore.Update(txn); err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	// add tokens to user
	user, err := rpg.UserStore.UserGetById(txn.UserID)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	user.Tokens += txn.Token

	if err := rpg.UserStore.Update(&user); err != nil {
		return nil, fmt.Errorf("failed to update user tokens: %w", err)
	}

	rpg.Logger.Printf("Order %s captured successfully for transaction ID %d", orderId, txn.ID)
	return txn, nil
}
