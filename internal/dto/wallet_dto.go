package dto

import "assignment-golang-backend/internal/entity"

func FormatWallet(wallet *entity.Wallet) *entity.Wallet {
	return &entity.Wallet{
		Base: entity.Base{
			ID: wallet.ID,
		},
		Number:  wallet.Number,
		Balance: wallet.Balance,
	}
}
