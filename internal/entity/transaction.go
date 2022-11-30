package entity

import (
	"fmt"
	"time"
)

type Transaction struct {
	Base
	Amount      int              `json:"amount"`
	Description string           `json:"description,omitempty"`
	Type        TransactionType  `json:"type"`
	Datetime    time.Time        `json:"datetime"`
	SourceID    *SourceOfFundsID `json:"source_id,omitempty"`
	From        int              `json:"from_number"           gorm:"column:from_number"`
	FromWallet  Wallet           `json:"from_wallet"           gorm:"references:Number;foreignKey:From;constraint:OnUpdate:CASCADE"`
	To          int              `json:"to_number"             gorm:"column:to_number"`
	ToWallet    Wallet           `json:"to_wallet"             gorm:"references:Number;foreignKey:To;constraint:OnUpdate:CASCADE"`
}

type SourceOfFundsID int

const (
	BankTransfer SourceOfFundsID = iota + 1
	CreditCard
	Cash
)

func (e SourceOfFundsID) String() string {
	switch e {
	case BankTransfer:
		return "Bank Transfer"
	case CreditCard:
		return "Credit Card"
	case Cash:
		return "Cash"
	default:
		return fmt.Sprintf("%d", int(e))
	}
}

type TransactionType string

const (
	Transfer TransactionType = "TRANSFER"
	TopUp    TransactionType = "TOP_UP"
)
