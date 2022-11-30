package dto

import "assignment-golang-backend/internal/entity"

func FormatUser(user *entity.User) *entity.User {
	return &entity.User{
		Base: entity.Base{
			ID: user.ID,
		},
		Name:         user.Name,
		Email:        user.Email,
		WalletNumber: user.WalletNumber,
		Wallet:       *FormatWallet(&user.Wallet),
	}
}
