package usecase

import (
	"assignment-golang-backend/internal/custom_error"
	"assignment-golang-backend/internal/entity"
	"assignment-golang-backend/internal/repository"
)

type IUserService interface {
	FindByID(int) (*entity.User, error)
}

type userService struct {
	userRepository repository.IUserRepository
}

func NewUserService(ur repository.IUserRepository) IUserService {
	return &userService{
		userRepository: ur,
	}
}

func (s *userService) FindByID(id int) (*entity.User, error) {
	user, rowsAffected, err := s.userRepository.FindByID(id)

	if rowsAffected == 0 {
		return nil, &custom_error.NoDataFound{DataType: "user"}
	}

	if err != nil {
		return nil, err
	}

	return user, nil

}
