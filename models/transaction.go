package models

import (
	"github.com/21TechLabs/factory-backend/utils"
	"gorm.io/gorm"
)

type TransactionStore struct {
	DB *gorm.DB
}

func NewTransactionStore(db *gorm.DB) *TransactionStore {
	return &TransactionStore{DB: db}
}

type Transaction struct {
	gorm.Model
	UserID           uint                    `gorm:"column:user_id" json:"userId"`
	User             User                    `gorm:"foreignKey:UserID;references:ID" json:"user"`
	Token            int64                   `gorm:"column:token" json:"token"`
	Amount           int64                   `gorm:"column:amount" json:"amount"`
	Currency         utils.Currency          `gorm:"column:currency" json:"currency"`
	Status           utils.TransactionStatus `gorm:"column:status" json:"status"`
	PaymentGatewayID utils.JSONMap           `gorm:"type:json column:payment_gateway_id" json:"paymentGatewayId"`
	TransactionID    string                  `gorm:"column:transaction_id" json:"transactionId"`
	PaymentPlanID    *uint                   `gorm:"column:payment_plan_id" json:"-"`
	PaymentPlan      *PaymentPlan            `gorm:"foreignKey:PaymentPlanID;references:ID" json:"paymentPlan"`
}
