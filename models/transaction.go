package models

import (
	"time"

	"github.com/21TechLabs/factory-backend/dto"
	"github.com/21TechLabs/factory-backend/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionStore struct {
	DB        *gorm.DB
	UserStore *UserStore
}

// NewTransactionStore creates a TransactionStore using the provided gorm.DB as its database connection.
func NewTransactionStore(db *gorm.DB) *TransactionStore {
	// var fileStore = NewFileStore(db)
	// var userStore = NewUserStore(db, fileStore)
	return &TransactionStore{DB: db}
}

type Transaction struct {
	ID                          uint                    `gorm:"primaryKey;autoIncrement" json:"id"`
	ReceiptId                   string                  `gorm:"column:receipt_id;unique" json:"receipt_id"`
	UserID                      uint                    `gorm:"column:user_id" json:"userId"`
	User                        User                    `gorm:"foreignKey:UserID;references:ID" json:"-"`
	Token                       int64                   `gorm:"column:token" json:"token"`
	Amount                      float64                 `gorm:"column:amount" json:"amount"`
	Currency                    utils.Currency          `gorm:"column:currency" json:"currency"`
	Status                      utils.TransactionStatus `gorm:"column:status" json:"status"`
	PaymentGatewayName          string                  `gorm:"column:payment_gateway_name" json:"paymentGatewayName"`
	PaymentGatewayRedirectURL   string                  `gorm:"column:payment_gateway_redirect_url" json:"paymentGatewayRedirectUrl"`
	PaymentGatewayTransactionID string                  `gorm:"column:transaction_id" json:"transactionId"`
	ProductPlanID               *uint                   `gorm:"column:product_plan_id" json:"-"`
	PaymentPlan                 *ProductPlan            `gorm:"foreignKey:ProductPlanID;references:ID" json:"productPlan"`
	CreatedAt                   time.Time               `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt                   time.Time               `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

// transaction := &Transaction{
// 		UserID:             rpg.User.ID,
// 		Token:              rpg.ProductPlan.Tokens,
// 		Amount:             rpg.ProductPlan.PlanPrice,
// 		Currency:           rpg.ProductPlan.PlanCurrency,
// 		Status:             utils.TransactionStatusPending,
// 		ProductPlanID:      &rpg.ProductPlan.ID,
// 		PaymentGatewayName: PaymentGatewayRazorpay,
// 	}

func (ts *TransactionStore) CreateTransaction(transaction *dto.TransactionCreateDto, user *User) (Transaction, error) {

	txn := Transaction{
		UserID:                      user.ID,
		Token:                       transaction.Token,
		Amount:                      transaction.Amount,
		Currency:                    transaction.Currency,
		Status:                      transaction.Status,
		ProductPlanID:               transaction.ProductPlanID,
		PaymentGatewayName:          transaction.PaymentGatewayName,
		PaymentGatewayRedirectURL:   transaction.PaymentGatewayRedirectURL,
		PaymentGatewayTransactionID: transaction.PaymentGatewayTransactionID,
		ReceiptId:                   uuid.NewString(),
	}

	tx := ts.DB.Create(&txn)

	if tx.Error != nil {
		return Transaction{}, tx.Error
	}

	return txn, nil
}

func (ts *TransactionStore) Update(txn *Transaction) error {
	if txn.ID == 0 {
		return utils.ErrTransactionNotFound
	}
	result := ts.DB.Save(txn)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (ts *TransactionStore) GetByPaymentGatewayTransactionID(transactionID string) (*Transaction, error) {
	var transaction Transaction
	query := ts.DB.Model(&Transaction{}).Where("payment_gateway_transaction_id = ?", transactionID)
	query = query.Preload("User")
	result := query.First(&transaction)
	if result.Error != nil {
		return nil, result.Error
	}
	return &transaction, nil
}
