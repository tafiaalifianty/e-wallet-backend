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
	"github.com/stretchr/testify/require"
)

func TestHandler_initAuthRoutes(t *testing.T) {
	router := SetUpRouter()
	service := &usecase.Services{}
	handler := New(service)
	group := router.Group("/")

	handler.initAuthRoutes(group)
}

func TestHandler_Register(t *testing.T) {
	mockDataInInterface, err := StructToMap(&entity.Token{})
	require.NoError(t, err)

	invalidBody := &dto.RegisterRequestBody{
		Name: "user",
	}
	validBody := &dto.RegisterRequestBody{
		Name:     "user",
		Email:    "user@email.com",
		Password: "password",
	}
	user := &entity.User{
		Name:     validBody.Name,
		Email:    validBody.Email,
		Password: validBody.Password,
	}
	tests := []struct {
		name        string
		body        io.Reader
		input       dto.RegisterRequestBody
		authService *mocks.IAuthService
		mock        func(*mocks.IAuthService)
		want        helper.JsonResponse
	}{

		{
			name:        "Error | Invalid Request Body",
			authService: mocks.NewIAuthService(t),
			body:        MakeRequestBody(invalidBody),
			mock: func(us *mocks.IAuthService) {
			},
			want: helper.JsonResponse{
				Code:    http.StatusBadRequest,
				Message: custom_error.InvalidRequestBody{}.Error(),
				Data:    nil,
			},
		},
		{
			name:        "Error | Error from Auth Service",
			authService: mocks.NewIAuthService(t),
			body:        MakeRequestBody(validBody),
			mock: func(us *mocks.IAuthService) {
				us.On("Register", user).
					Return(nil, &custom_error.FailedToCreateData{DataType: "User"})
			},
			want: helper.JsonResponse{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
				Data:    nil,
			},
		},
		{
			name:        "Success",
			authService: mocks.NewIAuthService(t),
			body:        MakeRequestBody(validBody),
			mock: func(us *mocks.IAuthService) {
				us.On("Register", user).
					Return(&entity.Token{}, nil)
			},
			want: helper.JsonResponse{
				Code:    http.StatusCreated,
				Message: http.StatusText(http.StatusCreated),
				Data:    mockDataInInterface,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				services: &usecase.Services{
					Auth: tt.authService,
				},
			}

			tt.mock(tt.authService)

			r := SetUpRouter()
			endpoint := "/api/auth/register"
			r.POST(endpoint, h.Register)
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

func TestHandler_Login(t *testing.T) {
	mockDataInInterface, err := StructToMap(&entity.Token{})
	require.NoError(t, err)

	invalidBody := &dto.LoginRequestBody{
		Email: "user@email.com",
	}
	validBody := &dto.LoginRequestBody{
		Email:    "user@email.com",
		Password: "password",
	}
	tests := []struct {
		name        string
		body        io.Reader
		input       dto.LoginRequestBody
		authService *mocks.IAuthService
		mock        func(*mocks.IAuthService)
		want        helper.JsonResponse
	}{

		{
			name:        "Error | Invalid Request Body",
			authService: mocks.NewIAuthService(t),
			body:        MakeRequestBody(invalidBody),
			mock: func(us *mocks.IAuthService) {
			},
			want: helper.JsonResponse{
				Code:    http.StatusBadRequest,
				Message: custom_error.InvalidRequestBody{}.Error(),
				Data:    nil,
			},
		},
		{
			name:        "Error | Error from AuthService",
			authService: mocks.NewIAuthService(t),
			body:        MakeRequestBody(validBody),
			mock: func(us *mocks.IAuthService) {
				us.On("Login", validBody.Email, validBody.Password).
					Return(nil, fmt.Errorf("error"))
			},
			want: helper.JsonResponse{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
				Data:    nil,
			},
		},
		{
			name:        "Success",
			authService: mocks.NewIAuthService(t),
			body:        MakeRequestBody(validBody),
			mock: func(us *mocks.IAuthService) {
				us.On("Login", validBody.Email, validBody.Password).
					Return(&entity.Token{}, nil)
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
					Auth: tt.authService,
				},
			}

			tt.mock(tt.authService)

			r := SetUpRouter()
			endpoint := "/api/auth/login"
			r.POST(endpoint, h.Login)
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
