package repository

import (
	"fmt"

	"assignment-golang-backend/internal/entity"

	"gorm.io/gorm"
)

type ITransactionRepository interface {
	CreateTransaction(
		*entity.Transaction,
	) (*entity.Transaction, int, error)
	FindByWalletNumberWithQuery(
		int,
		*entity.Pagination,
	) ([]*entity.Transaction, int, error)
	CountTransactionByWalletNumber(int, string) int
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) ITransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) CreateTransaction(
	transaction *entity.Transaction,
) (*entity.Transaction, int, error) {
	result := r.db.Create(&transaction)
	return transaction, int(result.RowsAffected), result.Error
}

func (r *transactionRepository) FindByWalletNumberWithQuery(
	walletNumber int,
	pagination *entity.Pagination,
) ([]*entity.Transaction, int, error) {
	var transactions []*entity.Transaction
	result := r.db.Where("transactions.from_number = ? OR transactions.to_number = ?", walletNumber, walletNumber).
		Where("description ILIKE ?", "%"+pagination.Search+"%").
		Order(fmt.Sprintf("transactions.%s %s", pagination.SortBy, pagination.Sort)).
		Offset((pagination.Page - 1) * pagination.Limit).
		Limit(pagination.Limit).
		Find(&transactions)
	return transactions, int(result.RowsAffected), result.Error
}

func (r *transactionRepository) CountTransactionByWalletNumber(
	walletNumber int,
	descriptionQuery string,
) int {
	var totalRows int64
	r.db.Model(&entity.Transaction{}).
		Where("transactions.from_number = ? OR transactions.to_number = ?", walletNumber, walletNumber).
		Where("description ILIKE ?", "%"+descriptionQuery+"%").
		Count(&totalRows)

	return int(totalRows)
}
