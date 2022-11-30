package dto

import (
	"time"

	"assignment-golang-backend/internal/entity"
)

type TopupRequestBody struct {
	Amount   int                    `json:"amount"    binding:"required"`
	SourceID entity.SourceOfFundsID `json:"source_id" binding:"required"`
}

type TransferRequestBody struct {
	Description string `json:"description"`
	To          int    `json:"To"          binding:"required"`
	Amount      int    `json:"amount"      binding:"required"`
}

type GetTransactionsByWalletNumberResponseBody struct {
	Pagination entity.Pagination       `json:"pagination"`
	Rows       []*FormattedTransaction `json:"rows"`
}

type FormattedTransaction struct {
	ID          int
	Amount      int                    `json:"amount"`
	Description string                 `json:"description,omitempty"`
	Type        entity.TransactionType `json:"type"`
	Datetime    time.Time              `json:"datetime"`
	Source      string                 `json:"source,omitempty"`
	From        int                    `json:"from,omitempty"`
	To          int                    `json:"to,omitempty"`
}

func FormatGetTransaction(
	transaction *entity.Transaction,
	sourceWalletNumber int,
) *FormattedTransaction {
	ResponseBody := &FormattedTransaction{
		ID:          transaction.ID,
		Description: transaction.Description,
		Type:        transaction.Type,
		Datetime:    transaction.Datetime,
		To:          transaction.To,
	}
	if transaction.Type == entity.TopUp {
		ResponseBody.Amount = transaction.Amount
		ResponseBody.Source = entity.SourceOfFundsID(*transaction.SourceID).
			String()
	} else if transaction.Type == entity.Transfer {
		ResponseBody.From = transaction.From
		if sourceWalletNumber == transaction.From {
			ResponseBody.Amount = -transaction.Amount
		} else {
			ResponseBody.Amount = transaction.Amount
		}
	} else {
		ResponseBody.Amount = transaction.Amount
	}

	return ResponseBody
}

func FormatMultipleGetTransaction(
	transactions []*entity.Transaction,
	sourceWalletNumber int,
) []*FormattedTransaction {
	formattedTransactions := []*FormattedTransaction{}
	for _, transaction := range transactions {
		formattedTransactions = append(
			formattedTransactions,
			FormatGetTransaction(transaction, sourceWalletNumber),
		)
	}

	return formattedTransactions
}

func FormatGetTransactionsByWalletNumberResponseBody(
	transactions []*entity.Transaction,
	pagination *entity.Pagination,
	sourceWalletNumber int,
) *GetTransactionsByWalletNumberResponseBody {
	transactionsFormatted := FormatMultipleGetTransaction(
		transactions,
		sourceWalletNumber,
	)

	return &GetTransactionsByWalletNumberResponseBody{
		Pagination: *pagination,
		Rows:       transactionsFormatted,
	}
}
