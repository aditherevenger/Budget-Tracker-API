package repository

import (
	"github.com/aditherevenger/Budget-Tracker-API/database"
	"github.com/aditherevenger/Budget-Tracker-API/models"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	Update(user *models.User) error
	Delete(user *models.User) error
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) Create(user *models.User) error {
	return database.DB.Create(user).Error
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := database.DB.First(&user, id).Error
	return &user, err
}

func (r *userRepository) Update(user *models.User) error {
	return database.DB.Save(user).Error
}

func (r *userRepository) Delete(user *models.User) error {
	return database.DB.Delete(user).Error
}
