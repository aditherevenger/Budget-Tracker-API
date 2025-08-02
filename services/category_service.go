package services

import (
	"github.com/aditherevenger/Budget-Tracker-API/models"
	"github.com/aditherevenger/Budget-Tracker-API/repository"
)

type CategoryService interface {
	CreateCategory(userID uint, req *models.CreateCategoryRequest) (*models.Category, error)
	GetCategories(userID uint) ([]models.Category, error)
	GetCategoryByID(id uint, userID uint) (*models.Category, error)
	UpdateCategory(id uint, userID uint, req *models.UpdateCategoryRequest) (*models.Category, error)
	DeleteCategory(id uint, userID uint) error
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *categoryService) CreateCategory(userID uint, req *models.CreateCategoryRequest) (*models.Category, error) {
	category := &models.Category{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
	}

	if category.Color == "" {
		category.Color = "#007bff"
	}

	err := s.categoryRepo.Create(category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) GetCategories(userID uint) ([]models.Category, error) {
	return s.categoryRepo.GetByUserID(userID, nil)
}

func (s *categoryService) GetCategoryByID(id uint, userID uint) (*models.Category, error) {
	return s.categoryRepo.GetByID(id, userID)
}

func (s *categoryService) UpdateCategory(id uint, userID uint, req *models.UpdateCategoryRequest) (*models.Category, error) {
	category, err := s.categoryRepo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		category.Name = *req.Name
	}

	if req.Description != nil {
		category.Description = *req.Description
	}

	if req.Color != nil {
		category.Color = *req.Color
	}

	err = s.categoryRepo.Update(category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) DeleteCategory(id uint, userID uint) error {
	return s.categoryRepo.Delete(id, userID)
}
