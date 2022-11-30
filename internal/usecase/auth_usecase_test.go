package usecase

import (
	"fmt"
	"testing"

	"assignment-golang-backend/internal/custom_error"
	"assignment-golang-backend/internal/entity"
	"assignment-golang-backend/internal/helper"
	"assignment-golang-backend/mocks"

	"github.com/stretchr/testify/assert"
)

func TestNewAuthService(t *testing.T) {
	NewAuthService(mocks.NewIUserRepository(t), mocks.NewIWalletRepository(t))
}

func Test_authService_Register(t *testing.T) {
	mockUser := &entity.User{
		Base: entity.Base{
			ID: 1,
		},
		Name:     "name",
		Email:    "email@email.com",
		Password: "password",
	}

	mockWallet := &entity.Wallet{}

	mockTokenString, _ := helper.GenerateJWT(mockUser)

	tests := []struct {
		name             string
		user             *entity.User
		userRepository   *mocks.IUserRepository
		walletRepository *mocks.IWalletRepository
		mock             func(*mocks.IUserRepository, *mocks.IWalletRepository)
		want             *entity.Token
		wantErr          bool
		expectedErr      error
	}{
		{
			name:             "ERROR | EmaiAlreadyUsedError",
			user:             mockUser,
			userRepository:   mocks.NewIUserRepository(t),
			walletRepository: mocks.NewIWalletRepository(t),
			mock: func(ir *mocks.IUserRepository, wr *mocks.IWalletRepository) {
				ir.On("FindByEmail", mockUser.Email).
					Return(mockUser, 1, nil)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: &custom_error.EmailAlreadyUsed{},
		},
		{
			name:             "ERROR | Other error from finding if email already exists",
			user:             mockUser,
			userRepository:   mocks.NewIUserRepository(t),
			walletRepository: mocks.NewIWalletRepository(t),
			mock: func(ir *mocks.IUserRepository, wr *mocks.IWalletRepository) {
				ir.On("FindByEmail", mockUser.Email).
					Return(mockUser, 0, fmt.Errorf("error"))
			},
			want:        nil,
			wantErr:     true,
			expectedErr: fmt.Errorf("error"),
		},
		{
			name:             "ERROR | NoDataCreated error if received error from wallet repository",
			user:             mockUser,
			userRepository:   mocks.NewIUserRepository(t),
			walletRepository: mocks.NewIWalletRepository(t),
			mock: func(ir *mocks.IUserRepository, wr *mocks.IWalletRepository) {
				ir.On("FindByEmail", mockUser.Email).
					Return(mockUser, 0, nil)
				wr.On("CreateWallet", &entity.Wallet{}).
					Return(&entity.Wallet{}, 0, &custom_error.FailedToCreateData{DataType: "Wallet"})
			},
			want:        nil,
			wantErr:     true,
			expectedErr: &custom_error.FailedToCreateData{DataType: "Wallet"},
		},
		{
			name:             "ERROR | NoDataCreated error if received error from user repository",
			user:             mockUser,
			userRepository:   mocks.NewIUserRepository(t),
			walletRepository: mocks.NewIWalletRepository(t),
			mock: func(ir *mocks.IUserRepository, wr *mocks.IWalletRepository) {
				ir.On("FindByEmail", mockUser.Email).
					Return(mockUser, 0, nil)
				wr.On("CreateWallet", mockWallet).
					Return(mockWallet, 1, nil)
				ir.On("CreateUser", mockUser).
					Return(mockUser, 1, fmt.Errorf("error"))
			},
			want:        nil,
			wantErr:     true,
			expectedErr: &custom_error.FailedToCreateData{DataType: "User"},
		},
		{
			name:             "SUCCESS",
			user:             mockUser,
			userRepository:   mocks.NewIUserRepository(t),
			walletRepository: mocks.NewIWalletRepository(t),
			mock: func(ir *mocks.IUserRepository, wr *mocks.IWalletRepository) {
				ir.On("FindByEmail", mockUser.Email).
					Return(mockUser, 0, nil)
				wr.On("CreateWallet", mockWallet).
					Return(mockWallet, 1, nil)
				ir.On("CreateUser", mockUser).
					Return(mockUser, 1, nil)
			},
			want:        &entity.Token{IDToken: mockTokenString},
			wantErr:     false,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &authService{
				userRepository:   tt.userRepository,
				walletRepository: tt.walletRepository,
			}

			tt.mock(tt.userRepository, tt.walletRepository)

			got, err := s.Register(tt.user)

			if !tt.wantErr {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.IDToken, got.IDToken)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_authService_Login(t *testing.T) {
	mockHashed, _ := helper.HashAndSalt("password")
	mockUser := &entity.User{
		Base: entity.Base{
			ID: 1,
		},
		Name:     "name",
		Email:    "email@email.com",
		Password: mockHashed,
	}

	mockTokenString, _ := helper.GenerateJWT(mockUser)
	type args struct {
		email    string
		password string
	}
	tests := []struct {
		name           string
		args           args
		userRepository *mocks.IUserRepository
		mock           func(*mocks.IUserRepository)
		want           *entity.Token
		wantErr        bool
		expectedErr    error
	}{
		{
			name:           "Error | NoDataFound error when no user found",
			args:           args{},
			userRepository: mocks.NewIUserRepository(t),
			mock: func(ir *mocks.IUserRepository) {
				ir.On("FindByEmail", "").Return(nil, 0, nil)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: &custom_error.FailedToCreateData{DataType: "user"},
		},
		{
			name:           "Error | Error other than no data found from service",
			args:           args{},
			userRepository: mocks.NewIUserRepository(t),
			mock: func(ir *mocks.IUserRepository) {
				ir.On("FindByEmail", "").Return(nil, 1, fmt.Errorf("error"))
			},
			want:        nil,
			wantErr:     true,
			expectedErr: fmt.Errorf("error"),
		},
		{
			name: "Error | Wrong Password Error",
			args: args{
				email:    mockUser.Email,
				password: "wrong_password",
			},
			userRepository: mocks.NewIUserRepository(t),
			mock: func(ir *mocks.IUserRepository) {
				ir.On("FindByEmail", mockUser.Email).Return(mockUser, 1, nil)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: &custom_error.WrongPassword{},
		},
		{
			name: "Success",
			args: args{
				email:    mockUser.Email,
				password: "password",
			},
			userRepository: mocks.NewIUserRepository(t),
			mock: func(ir *mocks.IUserRepository) {
				ir.On("FindByEmail", mockUser.Email).Return(mockUser, 1, nil)
			},
			want:        &entity.Token{IDToken: mockTokenString},
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &authService{
				userRepository: tt.userRepository,
			}

			tt.mock(tt.userRepository)

			got, err := s.Login(tt.args.email, tt.args.password)

			if !tt.wantErr {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.IDToken, got.IDToken)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
