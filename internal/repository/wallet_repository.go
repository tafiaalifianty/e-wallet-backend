package repository

import (
	"assignment-golang-backend/internal/entity"

	"gorm.io/gorm"
)

const (
	WALLET_STARTING_NUMBER = 100000
)

type IWalletRepository interface {
	CreateWallet(*entity.Wallet) (*entity.Wallet, int, error)
	FindByNumber(int) (*entity.Wallet, int, error)
	IncrementBalanceByValue(
		int, int,
	) (*entity.Wallet, int, error)
	DecrementBalanceByValue(
		int, int,
	) (*entity.Wallet, int, error)
}

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) IWalletRepository {
	return &walletRepository{
		db: db,
	}
}

func (r *walletRepository) CreateWallet(
	wallet *entity.Wallet,
) (*entity.Wallet, int, error) {
	result := r.db.Create(&wallet)
	if int(result.RowsAffected) == 0 && result.Error != nil {
		return nil, int(result.RowsAffected), result.Error
	}
	result = r.db.Model(&wallet).
		Update("number", WALLET_STARTING_NUMBER+wallet.ID)

	return wallet, int(result.RowsAffected), result.Error
}

func (r *walletRepository) FindByNumber(
	number int,
) (*entity.Wallet, int, error) {
	var wallet *entity.Wallet
	result := r.db.Where("number = ?", number).First(&wallet)
	return wallet, int(result.RowsAffected), result.Error
}

func (r *walletRepository) IncrementBalanceByValue(
	id, value int,
) (*entity.Wallet, int, error) {
	var wallet entity.Wallet
	r.db.Where("number = ?", id).First(&wallet)

	result := r.db.Model(&wallet).Update("balance", wallet.Balance+value)

	return &wallet, int(result.RowsAffected), result.Error
}

func (r *walletRepository) DecrementBalanceByValue(
	id, value int,
) (*entity.Wallet, int, error) {
	var wallet entity.Wallet
	r.db.Where("number = ?", id).First(&wallet)

	result := r.db.Model(&wallet).Update("balance", wallet.Balance-value)

	return &wallet, int(result.RowsAffected), result.Error
}
