package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"assignment-golang-backend/internal/custom_error"
	"assignment-golang-backend/internal/dto"
	"assignment-golang-backend/internal/entity"
	"assignment-golang-backend/internal/helper"
	"assignment-golang-backend/internal/usecase"
	"assignment-golang-backend/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_initTransactionRoutes(t *testing.T) {
	router := SetUpRouter()
	service := &usecase.Services{}
	handler := New(service)
	group := router.Group("/")

	handler.initTransactionRoutes(group)
}

func TestHandler_Topup(t *testing.T) {
	mockDataInInterface, err := StructToMap(&dto.FormattedTransaction{})
	require.NoError(t, err)

	tests := []struct {
		name                   string
		transactionService     *mocks.ITransactionService
		body                   io.Reader
		mockUserFromMiddleware bool
		mock                   func(*mocks.ITransactionService)
		want                   helper.JsonResponse
	}{

		{
			name:               "Error | Invalid Request Body",
			transactionService: mocks.NewITransactionService(t),
			body: MakeRequestBody(dto.TopupRequestBody{
				Amount: 100000,
			}),
			mockUserFromMiddleware: true,
			mock: func(ts *mocks.ITransactionService) {
			},
			want: helper.JsonResponse{
				Code:    http.StatusBadRequest,
				Message: custom_error.InvalidRequestBody{}.Error(),
				Data:    nil,
			},
		},
		{
			name:               "Error | Invalid Source ID Fund",
			transactionService: mocks.NewITransactionService(t),
			body: MakeRequestBody(dto.TopupRequestBody{
				Amount:   10,
				SourceID: 4,
			}),
			mockUserFromMiddleware: true,
			mock: func(ts *mocks.ITransactionService) {
			},
			want: helper.JsonResponse{
				Code:    http.StatusBadRequest,
				Message: custom_error.InvalidRequestBody{}.Error(),
				Data:    nil,
			},
		},
		{
			name:               "Error | Amount not between Min and Max amount of topup",
			transactionService: mocks.NewITransactionService(t),
			body: MakeRequestBody(dto.TopupRequestBody{
				Amount:   MAX_TOPUP_AMOUNT + 1,
				SourceID: 1,
			}),
			mockUserFromMiddleware: true,
			mock: func(ts *mocks.ITransactionService) {
			},
			want: helper.JsonResponse{
				Code: http.StatusBadRequest,
				Message: custom_error.AmountNotInRange{
					Minimum: MIN_TOPUP_AMOUNT,
					Maximum: MAX_TOPUP_AMOUNT,
				}.Error(),
				Data: nil,
			},
		},
		{
			name:               "Error | Failed to get user key from middleware",
			transactionService: mocks.NewITransactionService(t),
			body: MakeRequestBody(dto.TopupRequestBody{
				Amount:   MAX_TOPUP_AMOUNT - 1,
				SourceID: 1,
			}),
			mockUserFromMiddleware: false,
			mock: func(ts *mocks.ITransactionService) {
			},
			want: helper.JsonResponse{
				Code:    http.StatusInternalServerError,
				Message: custom_error.FailedToGetInfoFromToken{}.Error(),
				Data:    nil,
			},
		},
		{
			name:               "Error | Error from services",
			transactionService: mocks.NewITransactionService(t),
			body: MakeRequestBody(dto.TopupRequestBody{
				Amount:   MAX_TOPUP_AMOUNT - 1,
				SourceID: 1,
			}),
			mockUserFromMiddleware: true,
			mock: func(ts *mocks.ITransactionService) {
				ts.On("CreateTopup", mock.MatchedBy(func(i interface{}) bool {
					topup := i.(*entity.Transaction)
					return topup.Amount == MAX_TOPUP_AMOUNT-1
				})).Return(nil, fmt.Errorf("error"))
			},
			want: helper.JsonResponse{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
				Data:    nil,
			},
		},
		{
			name:               "Success",
			transactionService: mocks.NewITransactionService(t),
			body: MakeRequestBody(dto.TopupRequestBody{
				Amount:   MAX_TOPUP_AMOUNT - 1,
				SourceID: 1,
			}),
			mockUserFromMiddleware: true,
			mock: func(ts *mocks.ITransactionService) {
				ts.On("CreateTopup", mock.MatchedBy(func(i interface{}) bool {
					topup := i.(*entity.Transaction)
					return topup.Amount == MAX_TOPUP_AMOUNT-1
				})).Return(&entity.Transaction{}, nil)
			},
			want: helper.JsonResponse{
				Code:    http.StatusOK,
				Message: http.StatusText(http.StatusOK),
				Data:    mockDataInInterface,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				services: &usecase.Services{
					Transaction: tt.transactionService,
				},
			}

			tt.mock(tt.transactionService)

			r := SetUpRouter()

			endpoint := "/api/transaction/topup"
			if tt.mockUserFromMiddleware {
				r.POST(endpoint, MiddlewareMockUser, h.Topup)
			} else {
				r.POST(endpoint, h.Topup)
			}

			req, _ := http.NewRequest(
				http.MethodPost,
				endpoint,
				tt.body,
			)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			var response helper.JsonResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, tt.want.Code, w.Code)
			assert.Equal(t, tt.want, response)
		})
	}
}

func TestHandler_Transfer(t *testing.T) {
	validBody := dto.TransferRequestBody{
		Amount:      MIN_TRANSFER_AMOUNT,
		To:          2,
		Description: "description",
	}
	mockTransfer := &entity.Transaction{
		Amount:      validBody.Amount,
		To:          validBody.To,
		Description: validBody.Description,
	}

	mockDataInInterface, err := StructToMap(&dto.FormattedTransaction{
		Amount:      mockTransfer.Amount,
		To:          mockTransfer.To,
		Description: mockTransfer.Description,
	})
	require.NoError(t, err)

	tests := []struct {
		name                   string
		transactionService     *mocks.ITransactionService
		body                   io.Reader
		mockUserFromMiddleware bool
		mock                   func(*mocks.ITransactionService)
		want                   helper.JsonResponse
	}{

		{
			name:                   "Error | Invalid Request Body",
			transactionService:     mocks.NewITransactionService(t),
			body:                   MakeRequestBody(dto.TransferRequestBody{}),
			mockUserFromMiddleware: true,
			mock: func(ts *mocks.ITransactionService) {
			},
			want: helper.JsonResponse{
				Code:    http.StatusBadRequest,
				Message: custom_error.InvalidRequestBody{}.Error(),
				Data:    nil,
			},
		},
		{
			name:               "Error | Amount not between Min and Max amount of topup",
			transactionService: mocks.NewITransactionService(t),
			body: MakeRequestBody(dto.TransferRequestBody{
				Amount:      MIN_TRANSFER_AMOUNT - 1,
				To:          2,
				Description: "description",
			}),
			mockUserFromMiddleware: true,
			mock: func(ts *mocks.ITransactionService) {
			},
			want: helper.JsonResponse{
				Code: http.StatusBadRequest,
				Message: custom_error.AmountNotInRange{
					Minimum: MIN_TRANSFER_AMOUNT,
					Maximum: MAX_TRANSFER_AMOUNT,
				}.Error(),
				Data: nil,
			},
		},
		{
			name:                   "Error | Failed to get user key from middleware",
			transactionService:     mocks.NewITransactionService(t),
			body:                   MakeRequestBody(validBody),
			mockUserFromMiddleware: false,
			mock: func(ts *mocks.ITransactionService) {
			},
			want: helper.JsonResponse{
				Code:    http.StatusInternalServerError,
				Message: custom_error.FailedToGetInfoFromToken{}.Error(),
				Data:    nil,
			},
		},
		{
			name:               "Error | Destination wallet is same as user's wallet",
			transactionService: mocks.NewITransactionService(t),
			body: MakeRequestBody(dto.TransferRequestBody{
				Amount:      MIN_TRANSFER_AMOUNT,
				To:          MockTokenizedUser.WalletNumber,
				Description: "description",
			}),
			mockUserFromMiddleware: true,
			mock: func(ts *mocks.ITransactionService) {
			},
			want: helper.JsonResponse{
				Code:    http.StatusBadRequest,
				Message: custom_error.CannotTransferToOwnWallet{}.Error(),
				Data:    nil,
			},
		},
		{
			name:                   "Error | No data found from service",
			transactionService:     mocks.NewITransactionService(t),
			body:                   MakeRequestBody(validBody),
			mockUserFromMiddleware: true,
			mock: func(ts *mocks.ITransactionService) {
				ts.On("CreateTransaction", mock.MatchedBy(func(i interface{}) bool {
					transfer := i.(*entity.Transaction)
					return transfer.To == validBody.To &&
						transfer.Amount == validBody.Amount &&
						transfer.Description == validBody.Description
				})).
					Return(nil, &custom_error.NoDataFound{})
			},
			want: helper.JsonResponse{
				Code:    http.StatusNotFound,
				Message: custom_error.NoDataFound{}.Error(),
				Data:    nil,
			},
		},
		{
			name:                   "Error | Insufficient balance to transfer from services",
			transactionService:     mocks.NewITransactionService(t),
			body:                   MakeRequestBody(validBody),
			mockUserFromMiddleware: true,
			mock: func(ts *mocks.ITransactionService) {
				ts.On("CreateTransaction", mock.MatchedBy(func(i interface{}) bool {
					transfer := i.(*entity.Transaction)
					return transfer.To == validBody.To &&
						transfer.Amount == validBody.Amount &&
						transfer.Description == validBody.Description
				})).
					Return(nil, &custom_error.InsufficientBalance{})
			},
			want: helper.JsonResponse{
				Code:    http.StatusBadRequest,
				Message: custom_error.InsufficientBalance{}.Error(),
				Data:    nil,
			},
		},
		{
			name:                   "Error | Other error from transfer",
			transactionService:     mocks.NewITransactionService(t),
			body:                   MakeRequestBody(validBody),
			mockUserFromMiddleware: true,
			mock: func(ts *mocks.ITransactionService) {
				ts.On("CreateTransaction", mock.MatchedBy(func(i interface{}) bool {
					transfer := i.(*entity.Transaction)
					return transfer.To == validBody.To &&
						transfer.Amount == validBody.Amount &&
						transfer.Description == validBody.Description
				})).
					Return(nil, fmt.Errorf("error"))
			},
			want: helper.JsonResponse{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
				Data:    nil,
			},
		},
		{
			name:                   "Success",
			transactionService:     mocks.NewITransactionService(t),
			body:                   MakeRequestBody(validBody),
			mockUserFromMiddleware: true,
			mock: func(ts *mocks.ITransactionService) {
				ts.On("CreateTransaction", mock.MatchedBy(func(i interface{}) bool {
					transfer := i.(*entity.Transaction)
					return transfer.To == validBody.To &&
						transfer.Amount == validBody.Amount &&
						transfer.Description == validBody.Description
				})).
					Return(mockTransfer, nil)
			},
			want: helper.JsonResponse{
				Code:    http.StatusOK,
				Message: http.StatusText(http.StatusOK),
				Data:    mockDataInInterface,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				services: &usecase.Services{
					Transaction: tt.transactionService,
				},
			}

			tt.mock(tt.transactionService)

			r := SetUpRouter()

			endpoint := "/api/transaction/transfer"
			if tt.mockUserFromMiddleware {
				r.POST(endpoint, MiddlewareMockUser, h.Transfer)
			} else {
				r.POST(endpoint, h.Transfer)
			}

			req, _ := http.NewRequest(
				http.MethodPost,
				endpoint,
				tt.body,
			)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			var response helper.JsonResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, tt.want.Code, w.Code)
			assert.Equal(t, tt.want, response)
		})
	}
}

func TestHandler_GetTransactionsByWalletNumber(t *testing.T) {
	mockDefaultPagination := &entity.Pagination{
		Limit:      10,
		Page:       1,
		Search:     "",
		Sort:       "desc",
		SortBy:     "datetime",
		TotalRows:  0,
		TotalPages: 0,
	}

	mockDataInInterface, err := StructToMap(
		&dto.GetTransactionsByWalletNumberResponseBody{
			Pagination: *mockDefaultPagination,
			Rows:       []*dto.FormattedTransaction{},
		},
	)
	require.NoError(t, err)

	tests := []struct {
		name                   string
		transactionService     *mocks.ITransactionService
		mockUserFromMiddleware bool
		endpoint               string
		mock                   func(*mocks.ITransactionService)
		want                   helper.JsonResponse
	}{

		{
			name:                   "Error | Failed to get user key from middleware",
			transactionService:     mocks.NewITransactionService(t),
			mockUserFromMiddleware: false,
			mock: func(ts *mocks.ITransactionService) {
			},
			want: helper.JsonResponse{
				Code:    http.StatusInternalServerError,
				Message: custom_error.FailedToGetInfoFromToken{}.Error(),
				Data:    nil,
			},
		},
		{
			name:                   "Error | No Data Found from service",
			transactionService:     mocks.NewITransactionService(t),
			mockUserFromMiddleware: true,
			mock: func(ts *mocks.ITransactionService) {
				ts.On("FindByWalletNumber", MockTokenizedUser.WalletNumber, mockDefaultPagination).
					Return(nil, nil, &custom_error.NoDataFound{DataType: "transaction"})
			},
			want: helper.JsonResponse{
				Code: http.StatusNotFound,
				Message: custom_error.NoDataFound{
					DataType: "transaction",
				}.Error(),
				Data: nil,
			},
		},
		{
			name:                   "Error | Other error from service",
			transactionService:     mocks.NewITransactionService(t),
			mockUserFromMiddleware: true,
			mock: func(ts *mocks.ITransactionService) {
				ts.On("FindByWalletNumber", MockTokenizedUser.WalletNumber, mockDefaultPagination).
					Return(nil, nil, fmt.Errorf("error"))
			},
			want: helper.JsonResponse{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
				Data:    nil,
			},
		},
		{
			name:                   "Success",
			transactionService:     mocks.NewITransactionService(t),
			mockUserFromMiddleware: true,
			mock: func(ts *mocks.ITransactionService) {
				ts.On("FindByWalletNumber", MockTokenizedUser.WalletNumber, mockDefaultPagination).
					Return([]*entity.Transaction{}, mockDefaultPagination, nil)
			},
			want: helper.JsonResponse{
				Code:    http.StatusOK,
				Message: http.StatusText(http.StatusOK),
				Data:    mockDataInInterface,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				services: &usecase.Services{
					Transaction: tt.transactionService,
				},
			}

			tt.mock(tt.transactionService)

			r := SetUpRouter()

			endpoint := "/api/transaction"
			if tt.mockUserFromMiddleware {
				r.GET(
					endpoint,
					MiddlewareMockUser,
					h.GetTransactionsByWalletNumber,
				)
			} else {
				r.GET(endpoint, h.GetTransactionsByWalletNumber)
			}

			req, _ := http.NewRequest(
				http.MethodGet,
				endpoint,
				nil,
			)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			var response helper.JsonResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, tt.want.Code, w.Code)
			assert.Equal(t, tt.want, response)
		})
	}
}
