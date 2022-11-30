package usecase

import "assignment-golang-backend/internal/repository"

type Services struct {
	Auth        IAuthService
	User        IUserService
	Transaction ITransactionService
}

func New(r *repository.Repositories) *Services {
	return &Services{
		Auth:        NewAuthService(r.Users, r.Wallets),
		User:        NewUserService(r.Users),
		Transaction: NewTransactionService(r.Transactions, r.Wallets),
	}
}
