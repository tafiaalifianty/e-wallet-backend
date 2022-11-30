package entity

type Wallet struct {
	Base
	Number  int `json:"wallet_number" gorm:"unique"`
	Balance int `json:"balance"`
}
