package usecase

import (
	"fmt"
	"testing"

	"assignment-golang-backend/internal/custom_error"
	"assignment-golang-backend/internal/entity"
	"assignment-golang-backend/mocks"

	"github.com/stretchr/testify/assert"
)

func TestNewTransactionService(t *testing.T) {
	NewTransactionService(
		mocks.NewITransactionRepository(t),
		mocks.NewIWalletRepository(t),
	)
}

func Test_transactionService_CreateTopup(t *testing.T) {
	mockTopup := &entity.Transaction{
		To:     1,
		Amount: 1,
	}
	mockWallet := &entity.Wallet{}

	type repositories struct {
		transactionRepository *mocks.ITransactionRepository
		walletRepository      *mocks.IWalletRepository
	}
	tests := []struct {
		name         string
		repositories repositories
		mock         func(tr *mocks.ITransactionRepository, wr *mocks.IWalletRepository)
		topup        *entity.Transaction
		want         *entity.Transaction
		wantErr      bool
		expectedErr  error
	}{
		{
			name: "Error | Failed to create transaction data from transaction repository",
			repositories: repositories{
				transactionRepository: mocks.NewITransactionRepository(t),
				walletRepository:      mocks.NewIWalletRepository(t),
			},
			mock: func(tr *mocks.ITransactionRepository, wr *mocks.IWalletRepository) {
				tr.On("CreateTransaction", &entity.Transaction{}).
					Return(nil, 0, nil)
			},
			topup:   &entity.Transaction{},
			want:    nil,
			wantErr: true,
			expectedErr: &custom_error.FailedToCreateData{
				DataType: "transaction",
			},
		},
		{
			name: "Error | Other errors from transaction repository",
			repositories: repositories{
				transactionRepository: mocks.NewITransactionRepository(t),
				walletRepository:      mocks.NewIWalletRepository(t),
			},
			mock: func(tr *mocks.ITransactionRepository, wr *mocks.IWalletRepository) {
				tr.On("CreateTransaction", &entity.Transaction{}).
					Return(nil, 1, fmt.Errorf("error"))
			},
			topup:       &entity.Transaction{},
			want:        nil,
			wantErr:     true,
			expectedErr: fmt.Errorf("error"),
		},
		{
			name: "Error | Failed to create increment wallet balance from wallet repository",
			repositories: repositories{
				transactionRepository: mocks.NewITransactionRepository(t),
				walletRepository:      mocks.NewIWalletRepository(t),
			},
			mock: func(tr *mocks.ITransactionRepository, wr *mocks.IWalletRepository) {
				tr.On("CreateTransaction", mockTopup).
					Return(mockTopup, 1, nil)
				wr.On("IncrementBalanceByValue", mockTopup.To, mockTopup.Amount).
					Return(nil, 0, nil)
			},
			topup:   mockTopup,
			want:    nil,
			wantErr: true,
			expectedErr: &custom_error.FailedToUpdateData{
				DataType: "wallet balance",
			},
		},
		{
			name: "Error | Other errors from wallet repository",
			repositories: repositories{
				transactionRepository: mocks.NewITransactionRepository(t),
				walletRepository:      mocks.NewIWalletRepository(t),
			},
			mock: func(tr *mocks.ITransactionRepository, wr *mocks.IWalletRepository) {
				tr.On("CreateTransaction", mockTopup).
					Return(mockTopup, 1, nil)
				wr.On("IncrementBalanceByValue", mockTopup.To, mockTopup.Amount).
					Return(nil, 1, fmt.Errorf("error"))
			},
			topup:       mockTopup,
			want:        nil,
			wantErr:     true,
			expectedErr: fmt.Errorf("error"),
		},
		{
			name: "Success",
			repositories: repositories{
				transactionRepository: mocks.NewITransactionRepository(t),
				walletRepository:      mocks.NewIWalletRepository(t),
			},
			mock: func(tr *mocks.ITransactionRepository, wr *mocks.IWalletRepository) {
				tr.On("CreateTransaction", mockTopup).
					Return(mockTopup, 1, nil)
				wr.On("IncrementBalanceByValue", mockTopup.To, mockTopup.Amount).
					Return(mockWallet, 1, nil)
			},
			topup: mockTopup,
			want: &entity.Transaction{
				To:         mockTopup.To,
				Amount:     mockTopup.Amount,
				ToWallet:   *mockWallet,
				FromWallet: *mockWallet,
			},
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &transactionService{
				transactionRepository: tt.repositories.transactionRepository,
				walletRepository:      tt.repositories.walletRepository,
			}

			tt.mock(
				tt.repositories.transactionRepository,
				tt.repositories.walletRepository,
			)

			got, err := s.CreateTopup(tt.topup)

			if !tt.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_transactionService_CreateTransaction(t *testing.T) {
	mockTransfer := &entity.Transaction{
		Amount:      1000,
		Description: "",
		Type:        entity.Transfer,
		From:        1,
		To:          2,
	}
	mockFromWallet := &entity.Wallet{Number: 1, Balance: mockTransfer.Amount}
	mockToWallet := &entity.Wallet{Number: 2, Balance: mockTransfer.Amount}
	mockOtherError := fmt.Errorf("error")
	type repositories struct {
		transactionRepository *mocks.ITransactionRepository
		walletRepository      *mocks.IWalletRepository
	}
	tests := []struct {
		name         string
		repositories repositories
		mock         func(tr *mocks.ITransactionRepository, wr *mocks.IWalletRepository)
		transfer     *entity.Transaction
		want         *entity.Transaction
		wantErr      bool
		expectedErr  error
	}{
		{
			name: "Error | No source wallet found from wallet repository",
			repositories: repositories{
				transactionRepository: mocks.NewITransactionRepository(t),
				walletRepository:      mocks.NewIWalletRepository(t),
			},
			mock: func(tr *mocks.ITransactionRepository, wr *mocks.IWalletRepository) {
				wr.On("FindByNumber", mockTransfer.From).
					Return(nil, 0, nil)
			},
			transfer: mockTransfer,
			want:     nil,
			wantErr:  true,
			expectedErr: &custom_error.NoDataFound{
				DataType: "source wallet",
			},
		},
		{
			name: "Error | Other error from repository when finding source wallet",
			repositories: repositories{
				transactionRepository: mocks.NewITransactionRepository(t),
				walletRepository:      mocks.NewIWalletRepository(t),
			},
			mock: func(tr *mocks.ITransactionRepository, wr *mocks.IWalletRepository) {
				wr.On("FindByNumber", mockTransfer.From).
					Return(nil, 1, mockOtherError)
			},
			transfer:    mockTransfer,
			want:        nil,
			wantErr:     true,
			expectedErr: mockOtherError,
		},
		{
			name: "Error | Source wallet's balance is insufficient",
			repositories: repositories{
				transactionRepository: mocks.NewITransactionRepository(t),
				walletRepository:      mocks.NewIWalletRepository(t),
			},
			mock: func(tr *mocks.ITransactionRepository, wr *mocks.IWalletRepository) {
				wr.On("FindByNumber", mockTransfer.From).
					Return(&entity.Wallet{Balance: 0}, 1, nil)
			},
			transfer:    mockTransfer,
			want:        nil,
			wantErr:     true,
			expectedErr: &custom_error.InsufficientBalance{},
		},
		{
			name: "Error | Destination wallet not found from repository",
			repositories: repositories{
				transactionRepository: mocks.NewITransactionRepository(t),
				walletRepository:      mocks.NewIWalletRepository(t),
			},
			mock: func(tr *mocks.ITransactionRepository, wr *mocks.IWalletRepository) {
				wr.On("FindByNumber", mockTransfer.From).
					Return(mockFromWallet, 1, nil)
				wr.On("FindByNumber", mockTransfer.To).
					Return(nil, 0, nil)
			},
			transfer: mockTransfer,
			want:     nil,
			wantErr:  true,
			expectedErr: &custom_error.NoDataFound{
				DataType: "destination wallet",
			},
		},
		{
			name: "Error | Other error from repository when trying to get destination wallet",
			repositories: repositories{
				transactionRepository: mocks.NewITransactionRepository(t),
				walletRepository:      mocks.NewIWalletRepository(t),
			},
			mock: func(tr *mocks.ITransactionRepository, wr *mocks.IWalletRepository) {
				wr.On("FindByNumber", mockTransfer.From).
					Return(mockFromWallet, 1, nil)
				wr.On("FindByNumber", mockTransfer.To).
					Return(nil, 1, mockOtherError)
			},
			transfer:    mockTransfer,
			want:        nil,
			wantErr:     true,
			expectedErr: mockOtherError,
		},
		{
			name: "Error | Failed to decrement source wallet balance error",
			repositories: repositories{
				transactionRepository: mocks.NewITransactionRepository(t),
				walletRepository:      mocks.NewIWalletRepository(t),
			},
			mock: func(tr *mocks.ITransactionRepository, wr *mocks.IWalletRepository) {
				wr.On("FindByNumber", mockTransfer.From).
					Return(mockFromWallet, 1, nil)
				wr.On("FindByNumber", mockTransfer.To).
					Return(mockToWallet, 1, nil)
				wr.On("DecrementBalanceByValue", mockTransfer.From, mockTransfer.Amount).
					Return(nil, 0, nil)
			},
			transfer: mockTransfer,
			want:     nil,
			wantErr:  true,
			expectedErr: &custom_error.FailedToUpdateData{
				DataType: "source wallet balance",
			},
		},
		{
			name: "Error | Other error from repository when decrement source wallet balance",
			repositories: repositories{
				transactionRepository: mocks.NewITransactionRepository(t),
				walletRepository:      mocks.NewIWalletRepository(t),
			},
			mock: func(tr *mocks.ITransactionRepository, wr *mocks.IWalletRepository) {
				wr.On("FindByNumber", mockTransfer.From).
					Return(mockFromWallet, 1, nil)
				wr.On("FindByNumber", mockTransfer.To).
					Return(mockToWallet, 1, nil)
				wr.On("DecrementBalanceByValue", mockTransfer.From, mockTransfer.Amount).
					Return(nil, 1, mockOtherError)
			},
			transfer:    mockTransfer,
			want:        nil,
			wantErr:     true,
			expectedErr: mockOtherError,
		},
		{
			name: "Error | Failed to increment destination wallet balance error",
			repositories: repositories{
				transactionRepository: mocks.NewITransactionRepository(t),
				walletRepository:      mocks.NewIWalletRepository(t),
			},
			mock: func(tr *mocks.ITransactionRepository, wr *mocks.IWalletRepository) {
				wr.On("FindByNumber", mockTransfer.From).
					Return(mockFromWallet, 1, nil)
				wr.On("FindByNumber", mockTransfer.To).
					Return(mockToWallet, 1, nil)
				wr.On("DecrementBalanceByValue", mockTransfer.From, mockTransfer.Amount).
					Return(mockFromWallet, 1, nil)
				wr.On("IncrementBalanceByValue", mockTransfer.To, mockTransfer.Amount).
					Return(nil, 0, nil)
			},
			transfer: mockTransfer,
			want:     nil,
			wantErr:  true,
			expectedErr: &custom_error.FailedToUpdateData{
				DataType: "destination wallet balance",
			},
		},
		{
			name: "Error | Other error from repository when decrementing destination wallet balance",
			repositories: repositories{
				transactionRepository: mocks.NewITransactionRepository(t),
				walletRepository:      mocks.NewIWalletRepository(t),
			},
			mock: func(tr *mocks.ITransactionRepository, wr *mocks.IWalletRepository) {
				wr.On("FindByNumber", mockTransfer.From).
					Return(mockFromWallet, 1, nil)
				wr.On("FindByNumber", mockTransfer.To).
					Return(mockToWallet, 1, nil)
				wr.On("DecrementBalanceByValue", mockTransfer.From, mockTransfer.Amount).
					Return(mockFromWallet, 1, nil)
				wr.On("IncrementBalanceByValue", mockTransfer.To, mockTransfer.Amount).
					Return(nil, 1, mockOtherError)
			},
			transfer:    mockTransfer,
			want:        nil,
			wantErr:     true,
			expectedErr: mockOtherError,
		},
		{
			name: "Error | Failed to create transaction data",
			repositories: repositories{
				transactionRepository: mocks.NewITransactionRepository(t),
				walletRepository:      mocks.NewIWalletRepository(t),
			},
			mock: func(tr *mocks.ITransactionRepository, wr *mocks.IWalletRepository) {
				wr.On("FindByNumber", mockTransfer.From).
					Return(mockFromWallet, 1, nil)
				wr.On("FindByNumber", mockTransfer.To).
					Return(mockToWallet, 1, nil)
				wr.On("DecrementBalanceByValue", mockTransfer.From, mockTransfer.Amount).
					Return(mockFromWallet, 1, nil)
				wr.On("IncrementBalanceByValue", mockTransfer.To, mockTransfer.Amount).
					Return(mockToWallet, 1, nil)
				tr.On("CreateTransaction", mockTransfer).
					Return(nil, 0, nil)
			},
			transfer: mockTransfer,
			want:     nil,
			wantErr:  true,
			expectedErr: &custom_error.FailedToCreateData{
				DataType: "transaction",
			},
		},
		{
			name: "Error | Other error when creating repository data",
			repositories: repositories{
				transactionRepository: mocks.NewITransactionRepository(t),
				walletRepository:      mocks.NewIWalletRepository(t),
			},
			mock: func(tr *mocks.ITransactionRepository, wr *mocks.IWalletRepository) {
				wr.On("FindByNumber", mockTransfer.From).
					Return(mockFromWallet, 1, nil)
				wr.On("FindByNumber", mockTransfer.To).
					Return(mockToWallet, 1, nil)
				wr.On("DecrementBalanceByValue", mockTransfer.From, mockTransfer.Amount).
					Return(mockFromWallet, 1, nil)
				wr.On("IncrementBalanceByValue", mockTransfer.To, mockTransfer.Amount).
					Return(mockToWallet, 1, nil)
				tr.On("CreateTransaction", mockTransfer).
					Return(nil, 1, mockOtherError)
			},
			transfer:    mockTransfer,
			want:        nil,
			wantErr:     true,
			expectedErr: mockOtherError,
		},
		{
			name: "Success",
			repositories: repositories{
				transactionRepository: mocks.NewITransactionRepository(t),
				walletRepository:      mocks.NewIWalletRepository(t),
			},
			mock: func(tr *mocks.ITransactionRepository, wr *mocks.IWalletRepository) {
				wr.On("FindByNumber", mockTransfer.From).
					Return(mockFromWallet, 1, nil)
				wr.On("FindByNumber", mockTransfer.To).
					Return(mockToWallet, 1, nil)
				wr.On("DecrementBalanceByValue", mockTransfer.From, mockTransfer.Amount).
					Return(mockFromWallet, 1, nil)
				wr.On("IncrementBalanceByValue", mockTransfer.To, mockTransfer.Amount).
					Return(mockToWallet, 1, nil)
				tr.On("CreateTransaction", mockTransfer).
					Return(mockTransfer, 1, nil)
			},
			transfer: mockTransfer,
			want: &entity.Transaction{
				Amount:      mockTransfer.Amount,
				Description: mockTransfer.Description,
				Type:        mockTransfer.Type,
				From:        mockTransfer.From,
				To:          mockTransfer.To,
				FromWallet:  *mockFromWallet,
				ToWallet:    *mockToWallet,
			},
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &transactionService{
				transactionRepository: tt.repositories.transactionRepository,
				walletRepository:      tt.repositories.walletRepository,
			}

			tt.mock(
				tt.repositories.transactionRepository,
				tt.repositories.walletRepository,
			)

			got, err := s.CreateTransaction(tt.transfer)

			if !tt.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_transactionService_FindByWalletNumber(t *testing.T) {

	tests := []struct {
		name                  string
		transactionRepository *mocks.ITransactionRepository
		mock                  func(tr *mocks.ITransactionRepository)
		walletNumber          int
		pagination            *entity.Pagination
		want                  []*entity.Transaction
		want1                 *entity.Pagination
		wantErr               bool
		expectedErr           error
	}{
		{
			name:                  "Error | No Data Found",
			transactionRepository: mocks.NewITransactionRepository(t),
			mock: func(tr *mocks.ITransactionRepository) {
				tr.On("FindByWalletNumberWithQuery", 1, &entity.Pagination{}).
					Return(nil, 0, nil)
			},
			walletNumber: 1,
			pagination:   &entity.Pagination{},
			want:         []*entity.Transaction(nil),
			wantErr:      true,
			expectedErr:  &custom_error.NoDataFound{DataType: "transaction"},
		},
		{
			name:                  "Error | Other error from repository",
			transactionRepository: mocks.NewITransactionRepository(t),
			mock: func(tr *mocks.ITransactionRepository) {
				tr.On("FindByWalletNumberWithQuery", 1, &entity.Pagination{}).
					Return(nil, 1, fmt.Errorf("error"))
			},
			walletNumber: 1,
			pagination:   &entity.Pagination{},
			want:         []*entity.Transaction(nil),
			wantErr:      true,
			expectedErr:  fmt.Errorf("error"),
		},
		{
			name:                  "Success",
			transactionRepository: mocks.NewITransactionRepository(t),
			mock: func(tr *mocks.ITransactionRepository) {
				tr.On("FindByWalletNumberWithQuery", 1, &entity.Pagination{Limit: 1}).
					Return([]*entity.Transaction{}, 1, nil)
				tr.On("CountTransactionByWalletNumber", 1, "").
					Return(10)
			},
			walletNumber: 1,
			pagination:   &entity.Pagination{Limit: 1},
			want:         []*entity.Transaction{},
			want1: &entity.Pagination{
				TotalRows:  10,
				TotalPages: 10,
				Limit:      1,
			},
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &transactionService{
				transactionRepository: tt.transactionRepository,
			}

			tt.mock(
				tt.transactionRepository,
			)

			got, got1, err := s.FindByWalletNumber(
				tt.walletNumber,
				tt.pagination,
			)

			if !tt.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			}

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}
