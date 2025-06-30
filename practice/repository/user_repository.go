package repository

import (
	"errors"
	"practice/model"

	"gorm.io/gorm"
)

type IUserRepository interface {
	GetUserByEmail(user *model.User, email string) error
	CreateUser(user *model.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db}
}

func (ur *userRepository) GetUserByEmail(user *model.User, email string) error {
	if ur.db == nil {
		return errors.New("database not connected")
	}
	if err := ur.db.Where("email = ?", email).First(user).Error; err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) CreateUser(user *model.User) error {
	if ur.db == nil {
		return errors.New("database not connected")
	}
	if err := ur.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}
