package repository

import (
	"github.com/aditherevenger/Budget-Tracker-API/database"
	"github.com/aditherevenger/Budget-Tracker-API/models"
)

type CategoryRepository interface {
	Create(category *models.Category) error
	GetByUserID(userID uint, filter *models.User) ([]models.Category, error)
	GetByID(id uint, userID uint) (*models.Category, error)
	Update(category *models.Category) error
	Delete(id uint, userID uint) error
}

type categoryRepository struct{}

func NewCategoryRepository() CategoryRepository {
	return &categoryRepository{}
}

func (r *categoryRepository) Create(category *models.Category) error {
	return database.DB.Create(category).Error
}

func (r *categoryRepository) GetByUserID(userID uint, filter *models.User) ([]models.Category, error) {
	var categories []models.Category
	err := database.DB.Where("user_id = ?", userID).Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) GetByID(id uint, userID uint) (*models.Category, error) {
	var category models.Category
	err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&category).Error
	return &category, err
}

func (r *categoryRepository) Update(category *models.Category) error {
	return database.DB.Save(category).Error
}

func (r *categoryRepository) Delete(id uint, userID uint) error {
	return database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Category{}).Error
}
