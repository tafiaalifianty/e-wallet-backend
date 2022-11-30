package entity

type User struct {
	Base
	Name         string `json:"name"`
	Email        string `json:"email"              gorm:"unique"`
	Password     string `json:"password,omitempty"`
	WalletNumber int    `json:"wallet_number"`
	Wallet       Wallet `json:"wallet"             gorm:"references:Number;foreignKey:WalletNumber;constraint:OnUpdate:CASCADE"`
}
type TokenizedUser struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	WalletNumber int    `json:"wallet_number"`
}
