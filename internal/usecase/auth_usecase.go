package usecase

import (
	"assignment-golang-backend/internal/custom_error"
	"assignment-golang-backend/internal/entity"
	"assignment-golang-backend/internal/helper"
	"assignment-golang-backend/internal/repository"
)

type IAuthService interface {
	Login(string, string) (*entity.Token, error)
	Register(*entity.User) (*entity.Token, error)
}

type authService struct {
	userRepository   repository.IUserRepository
	walletRepository repository.IWalletRepository
}

func NewAuthService(
	ur repository.IUserRepository,
	wr repository.IWalletRepository,
) IAuthService {
	return &authService{
		userRepository:   ur,
		walletRepository: wr,
	}
}

func (s *authService) Login(
	email, password string,
) (*entity.Token, error) {
	user, rowsAffected, err := s.userRepository.FindByEmail(email)

	if rowsAffected == 0 {
		return nil, &custom_error.FailedToCreateData{DataType: "user"}
	}

	if err != nil {
		return nil, err
	}

	if !helper.ComparePasswords(user.Password, []byte(password)) {
		return nil, &custom_error.WrongPassword{}
	}

	if err != nil {
		return nil, err
	}

	tokenString, err := helper.GenerateJWT(user)
	if err != nil {
		return nil, err
	}

	return &entity.Token{IDToken: tokenString, User: *user}, nil
}

func (s *authService) Register(user *entity.User) (*entity.Token, error) {
	var err error

	_, rowsAffected, err := s.userRepository.FindByEmail(user.Email)

	if rowsAffected != 0 {
		return nil, &custom_error.EmailAlreadyUsed{}
	}

	if err != nil {
		return nil, err
	}

	user.Password, err = helper.HashAndSalt(user.Password)

	if err != nil {
		return nil, err
	}

	wallet := &entity.Wallet{}
	wallet, rowsAffected, err = s.walletRepository.CreateWallet(wallet)

	if rowsAffected == 0 || err != nil {
		return nil, &custom_error.FailedToCreateData{DataType: "Wallet"}
	}

	user.Wallet = *wallet
	user.WalletNumber = wallet.Number

	user, rowsAffected, err = s.userRepository.CreateUser(user)

	if rowsAffected == 0 || err != nil {
		return nil, &custom_error.FailedToCreateData{DataType: "User"}
	}

	tokenString, err := helper.GenerateJWT(user)

	if err != nil {
		return nil, err
	}

	return &entity.Token{IDToken: tokenString, User: *user}, nil
}
