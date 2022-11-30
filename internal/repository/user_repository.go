package repository

import (
	"assignment-golang-backend/internal/entity"

	"gorm.io/gorm"
)

type IUserRepository interface {
	CreateUser(*entity.User) (*entity.User, int, error)
	FindByID(int) (*entity.User, int, error)
	FindByEmail(string) (*entity.User, int, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUser(
	user *entity.User,
) (*entity.User, int, error) {
	result := r.db.Create(&user)
	return user, int(result.RowsAffected), result.Error
}

func (r *userRepository) FindByID(id int) (*entity.User, int, error) {
	var user *entity.User
	result := r.db.Joins("Wallet").First(&user, id)
	return user, int(result.RowsAffected), result.Error
}

func (r *userRepository) FindByEmail(email string) (*entity.User, int, error) {
	var user *entity.User
	result := r.db.Where("email = ?", email).Find(&user)
	return user, int(result.RowsAffected), result.Error
}
