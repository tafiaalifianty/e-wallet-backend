package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"assignment-golang-backend/internal/custom_error"
	"assignment-golang-backend/internal/entity"
	"assignment-golang-backend/internal/helper"
	"assignment-golang-backend/internal/usecase"
	"assignment-golang-backend/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_initUserRoutes(t *testing.T) {
	router := SetUpRouter()
	service := &usecase.Services{}
	handler := New(service)
	group := router.Group("/")

	handler.initUserRoutes(group)
}

func TestHandler_GetUserInfo(t *testing.T) {
	mockUser := &entity.User{
		Base: entity.Base{
			ID: MockTokenizedUser.ID,
		},
		Name:         MockTokenizedUser.Name,
		Email:        MockTokenizedUser.Email,
		WalletNumber: MockTokenizedUser.WalletNumber,
	}

	mockDataInInterface, err := StructToMap(mockUser)
	require.NoError(t, err)

	tests := []struct {
		name                   string
		userService            *mocks.IUserService
		mockUserFromMiddleware bool
		mock                   func(*mocks.IUserService)
		want                   helper.JsonResponse
	}{

		{
			name:                   "Error | Invalid or No User Key from middleware",
			userService:            mocks.NewIUserService(t),
			mockUserFromMiddleware: false,
			mock: func(us *mocks.IUserService) {
			},
			want: helper.JsonResponse{
				Code:    http.StatusBadRequest,
				Message: http.StatusText(http.StatusBadRequest),
				Data:    nil,
			},
		},
		{
			name:                   "Error | No User found from Service",
			userService:            mocks.NewIUserService(t),
			mockUserFromMiddleware: true,
			mock: func(us *mocks.IUserService) {
				us.On("FindByID", int(MockTokenizedUser.ID)).
					Return(nil, &custom_error.NoDataFound{DataType: "user"})
			},
			want: helper.JsonResponse{
				Code: http.StatusNotFound,
				Message: custom_error.NoDataFound{
					DataType: "user",
				}.Error(),
				Data: nil,
			},
		},
		{
			name:                   "Error | Other errors from Service",
			userService:            mocks.NewIUserService(t),
			mockUserFromMiddleware: true,
			mock: func(us *mocks.IUserService) {
				us.On("FindByID", int(MockTokenizedUser.ID)).
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
			userService:            mocks.NewIUserService(t),
			mockUserFromMiddleware: true,
			mock: func(us *mocks.IUserService) {
				us.On("FindByID", int(MockTokenizedUser.ID)).
					Return(mockUser, nil)
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
					User: tt.userService,
				},
			}

			tt.mock(tt.userService)

			r := SetUpRouter()

			endpoint := "/api/users/info"
			if tt.mockUserFromMiddleware {
				r.GET(endpoint, MiddlewareMockUser, h.GetUserInfo)
			} else {
				r.GET(endpoint, h.GetUserInfo)
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
