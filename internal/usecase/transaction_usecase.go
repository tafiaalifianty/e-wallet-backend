package usecase

import (
	"math"

	"assignment-golang-backend/internal/custom_error"
	"assignment-golang-backend/internal/entity"
	"assignment-golang-backend/internal/repository"
)

type ITransactionService interface {
	CreateTopup(*entity.Transaction) (*entity.Transaction, error)
	CreateTransaction(*entity.Transaction) (*entity.Transaction, error)
	FindByWalletNumber(
		int,
		*entity.Pagination,
	) ([]*entity.Transaction, *entity.Pagination, error)
}

type transactionService struct {
	transactionRepository repository.ITransactionRepository
	walletRepository      repository.IWalletRepository
}

func NewTransactionService(
	tr repository.ITransactionRepository,
	wr repository.IWalletRepository,
) ITransactionService {
	return &transactionService{
		transactionRepository: tr,
		walletRepository:      wr,
	}
}

func (s *transactionService) CreateTransaction(
	transferRecord *entity.Transaction,
) (*entity.Transaction, error) {
	fromWallet, rowsAffected, err := s.walletRepository.FindByNumber(
		transferRecord.From,
	)

	if rowsAffected == 0 {
		return nil, &custom_error.NoDataFound{DataType: "source wallet"}
	}

	if err != nil {
		return nil, err
	}

	if fromWallet.Balance < transferRecord.Amount {
		return nil, &custom_error.InsufficientBalance{}
	}

	_, rowsAffected, err = s.walletRepository.FindByNumber(
		transferRecord.To,
	)

	if rowsAffected == 0 {
		return nil, &custom_error.NoDataFound{DataType: "destination wallet"}
	}

	if err != nil {
		return nil, err
	}

	fromWallet, rowsAffected, err = s.walletRepository.DecrementBalanceByValue(
		transferRecord.From,
		transferRecord.Amount,
	)

	if rowsAffected == 0 {
		return nil, &custom_error.FailedToUpdateData{
			DataType: "source wallet balance",
		}
	}

	if err != nil {
		return nil, err
	}

	toWallet, rowsAffected, err := s.walletRepository.IncrementBalanceByValue(
		transferRecord.To,
		transferRecord.Amount,
	)

	if rowsAffected == 0 {
		return nil, &custom_error.FailedToUpdateData{
			DataType: "destination wallet balance",
		}
	}

	if err != nil {
		return nil, err
	}

	transferRecord, rowsAffected, err = s.transactionRepository.CreateTransaction(
		transferRecord,
	)

	if rowsAffected == 0 {
		return nil, &custom_error.FailedToCreateData{DataType: "transaction"}
	}

	if err != nil {
		return nil, err
	}

	transferRecord.FromWallet = *fromWallet
	transferRecord.ToWallet = *toWallet

	return transferRecord, nil
}

func (s *transactionService) CreateTopup(
	topup *entity.Transaction,
) (*entity.Transaction, error) {
	topup, rowsAffected, err := s.transactionRepository.CreateTransaction(topup)

	if rowsAffected == 0 {
		return nil, &custom_error.FailedToCreateData{DataType: "transaction"}
	}

	if err != nil {
		return nil, err
	}

	wallet, rowsAffected, err := s.walletRepository.IncrementBalanceByValue(
		topup.To,
		topup.Amount,
	)

	if rowsAffected == 0 {
		return nil, &custom_error.FailedToUpdateData{
			DataType: "wallet balance",
		}
	}

	if err != nil {
		return nil, err
	}

	topup.ToWallet = *wallet
	topup.FromWallet = *wallet

	return topup, nil
}

func (s *transactionService) FindByWalletNumber(
	walletNumber int,
	pagination *entity.Pagination,
) ([]*entity.Transaction, *entity.Pagination, error) {
	transactions, rowsAffected, err := s.transactionRepository.FindByWalletNumberWithQuery(
		walletNumber,
		pagination,
	)

	if rowsAffected == 0 {
		return nil, nil, &custom_error.NoDataFound{DataType: "transaction"}
	}

	if err != nil {
		return nil, nil, err
	}

	totalRows := s.transactionRepository.CountTransactionByWalletNumber(
		walletNumber,
		pagination.Search,
	)

	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))

	pagination.TotalPages = totalPages
	pagination.TotalRows = totalRows

	return transactions, pagination, nil
}
