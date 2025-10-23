package models

import (
	"fmt"
	"log"
	"time"

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
	Logger                *log.Logger
	TransactionStore      *TransactionStore
	UserSubscriptionStore *UserSubscriptionStore
	UserStore             *UserStore
	Client                *razorpay.Client
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

/*

data := map[string]interface{}{
  "plan_id":"plan_JcwJfpjN6VHSGv",
  "total_count":3,
  "quantity": 1,
  "customer_notify":1,
  "addons":[]interface{}{
    map[string]interface{}{
      "item":map[string]interface{}{
        "name":"Delivery charges",
        "amount":3000,
        "currency":"INR",
      },
    },
  },
  "notes":map[string]interface{}{
    "notes_key_1":"Tea, Earl Grey, Hot",
    "notes_key_2":"Tea, Earl Grey… decaf.",
  },
}
*/

type RazorpayCreateSubscription struct {
	PlanID         string                `json:"plan_id"`
	TotalCount     int                   `json:"total_count"`
	Quantity       int                   `json:"quantity"`
	CustomerNotify int                   `json:"customer_notify"`
	Addons         utils.JSONMap[string] `json:"addons"`
	Notes          utils.JSONMap[string] `json:"notes"`
}

func (r RazorpayCreateSubscription) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"plan_id":         r.PlanID,
		"total_count":     r.TotalCount,
		"quantity":        r.Quantity,
		"customer_notify": r.CustomerNotify,
		"addons":          r.Addons,
		"notes":           r.Notes,
	}
}

func NewRazorpayPG(log *log.Logger, transactionStore *TransactionStore, userStore *UserStore, uss *UserSubscriptionStore) *RazorpayPG {
	return &RazorpayPG{
		Logger:                log,
		TransactionStore:      transactionStore,
		Client:                RazorpayClient,
		UserStore:             userStore,
		UserSubscriptionStore: uss,
	}
}

func (rpg *RazorpayPG) InitiatePayment(productPlan *ProductPlan, user *User, count int) (*Transaction, error) {
	switch productPlan.PlanType {
	case utils.PlanTypeOneTime:
		return rpg.initiatePaymentOneTime(productPlan, user, count)
	case utils.PlanTypeSubscription:
		return rpg.initiatePaymentSubscription(productPlan, user)
	default:
		return nil, utils.ErrInvalidPlanType
	}

}

func (rpg *RazorpayPG) initiatePaymentOneTime(productPlan *ProductPlan, user *User, count int) (*Transaction, error) {
	// create a transaction
	txn := &dto.TransactionCreateDto{
		Token:                       productPlan.Tokens * int64(count),
		Amount:                      float64(count) * (productPlan.PlanPrice) * 100, // Convert to smallest currency unit
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
		Amount:         txn.Amount, // Convert to smallest currency unit
		Currency:       txn.Currency,
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
	subscriptionPlanID := productPlan.PaymentGatewayID[PaymentGatewayRazorpay]

	var subscriptionBody = RazorpayCreateSubscription{
		PlanID:         subscriptionPlanID,
		TotalCount:     1,
		Quantity:       1,
		CustomerNotify: 1,
		Addons:         nil,
		Notes:          map[string]string{"user_id": fmt.Sprintf("%d", user.ID)},
	}
	_sub, err := rpg.Client.Subscription.Create(subscriptionBody.ToMap(), nil)
	if err != nil {
		return nil, err
	}

	sub, err := utils.MapToStruct[RazorpaySubscriptionCreateEvent](_sub)
	if err != nil {
		return nil, err
	}

	txn := &dto.TransactionCreateDto{
		Token:                       productPlan.Tokens,
		Amount:                      (productPlan.PlanPrice) * 100, // Convert to the smallest currency unit
		Currency:                    productPlan.PlanCurrency,
		Status:                      utils.TransactionStatusPending,
		ProductPlanID:               &productPlan.ID,
		PaymentGatewayName:          PaymentGatewayRazorpay,
		PaymentGatewayRedirectURL:   sub.ShortURL,
		PaymentGatewayTransactionID: sub.ID,
	}

	startAt := time.Unix(sub.StartAt, 0)
	endAt := time.Unix(sub.EndAt, 0)

	userSub := UserSubscription{
		UserID:             user.ID,
		StartDate:          startAt,
		EndDate:            endAt,
		IsActive:           false,
		SubscriptionStatus: utils.SubscriptionStatusPending,
		PaymentGatewayName: PaymentGatewayRazorpay,
		SubscriptionID:     sub.ID,
		ProductPlanID:      productPlan.ID,
		ChargedCount:       sub.PaidCount,
		TotalChargedCount:  sub.TotalCount,
		Suspended:          false,
	}

	transaction, err := rpg.TransactionStore.CreateTransaction(txn, user)
	if err != nil {
		return nil, err
	}

	if _, err = rpg.UserSubscriptionStore.Create(userSub); err != nil {
		return nil, err
	}

	return &transaction, nil
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

	if txn.Status == utils.TransactionStatusCompleted {
		rpg.Logger.Printf("Order %s already captured; ignoring duplicate event", orderId)
		return txn, nil
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

func (rpg *RazorpayPG) ProcessFailedPayments(event RazorpayBaseEvent[RazorpayPaymentFailedPayload]) (*Transaction, error) {
	orderId := event.Payload.Payment.Entity.OrderID

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

	txn, err := rpg.TransactionStore.GetByPaymentGatewayTransactionID(event.Payload.Payment.Entity.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by order ID: %w", err)
	}

	if txn == nil {
		return nil, utils.ErrTransactionNotFound
	}

	if txn.Status == utils.TransactionStatusCompleted {
		rpg.Logger.Printf("Order %s already processed; ignoring event", orderId)
		return txn, nil
	}

	txn.Status = utils.TransactionStatusFailed
	txn.PaymentGatewayRedirectURL = ""

	if err := rpg.TransactionStore.Update(txn); err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	rpg.Logger.Printf("%s failed for transaction ID %d", orderId, txn.ID)
	return txn, nil
}

func (rpg *RazorpayPG) ProcessSubscriptions(subscription RazorpayBaseEvent[RazorpaySubscriptionEventsPayload]) (*Transaction, error) {
	subId := subscription.Payload.Subscription.Entity.ID

	if subId == "" {
		return nil, utils.ErrInvalidOrderID
	}

	//fetch subscription
	sub, err := rpg.Client.Subscription.Fetch(subId, nil, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch subscription: %w", err)
	}

	if sub == nil {
		return nil, utils.ErrSubscriptionNotFound
	}

	subEvent := utils.SubscriptionStatus(subscription.Event)

	userSub, err := rpg.UserSubscriptionStore.FindBySubscriptionID(subId)
	if err != nil {
		return nil, err
	}

	subEntity := subscription.Payload.Subscription.Entity

	userSub.SubscriptionStatus = subEvent
	userSub.ChargedCount = subEntity.PaidCount
	userSub.TotalChargedCount = subEntity.TotalCount

	userSub.StartDate = time.Unix(subEntity.StartAt, 0)
	userSub.EndDate = time.Unix(subEntity.EndAt, 0)

	switch subEvent {
	case utils.SubscriptionStatusActive:
		userSub.IsActive = true
	case utils.SubscriptionStatusPaused:
		userSub.IsActive = false
	case utils.SubscriptionStatusResumed:
		userSub.IsActive = true
	case utils.SubscriptionStatusCancelled:
		userSub.IsActive = false
	case utils.SubscriptionStatusCompleted:
		userSub.IsActive = false
	case utils.SubscriptionStatusPending:
		userSub.IsActive = false
	case utils.SubscriptionStatusHalted:
		userSub.IsActive = false
	case utils.SubscriptionStatusCharged:
		userSub.IsActive = true
		// everytime a sub is charged, we need to create a new txn
		productPlan := userSub.ProductPlan
		paymentEntity := subscription.Payload.Payment.Entity
		txn := &dto.TransactionCreateDto{
			Token:                       productPlan.Tokens,
			Amount:                      float64(paymentEntity.Amount),
			Currency:                    productPlan.PlanCurrency,
			Status:                      utils.TransactionStatusCompleted,
			ProductPlanID:               &(productPlan).ID,
			PaymentGatewayName:          PaymentGatewayRazorpay,
			PaymentGatewayRedirectURL:   "",
			PaymentGatewayTransactionID: paymentEntity.ID,
		}
		_, err = rpg.TransactionStore.CreateTransaction(txn, &userSub.User)

		if err != nil {
			return nil, err
		}
	default:
		return nil, utils.ErrInvalidSubscription
	}

	if err = rpg.UserSubscriptionStore.Save(userSub); err != nil {
		return nil, err
	}

	return nil, nil
}
