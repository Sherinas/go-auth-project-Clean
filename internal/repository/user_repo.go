package repository

import (
	"github.com/Sherinas/go-auth-project-Clean/internal/domain"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (u *userRepo) Create(user *domain.User) error {
	return u.db.Create(user).Error
}

func (u *userRepo) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := u.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
