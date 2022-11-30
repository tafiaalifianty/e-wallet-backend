package repository

import "gorm.io/gorm"

type Repositories struct {
	Users        IUserRepository
	Wallets      IWalletRepository
	Transactions ITransactionRepository
}

func New(db *gorm.DB) *Repositories {
	return &Repositories{
		Users:        NewUserRepository(db),
		Wallets:      NewWalletRepository(db),
		Transactions: NewTransactionRepository(db),
	}
}
