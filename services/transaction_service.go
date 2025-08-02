package services

import (
	"errors"
	"github.com/aditherevenger/Budget-Tracker-API/models"
	"github.com/aditherevenger/Budget-Tracker-API/repository"
)

type TransactionService interface {
	CreateTransaction(userID uint, req *models.CreateTransactionRequest) (*models.Transaction, error)
	GetTransactions(userID uint, filter *models.TransactionFilter) ([]models.Transaction, error)
	GetTransactionByID(id uint, userID uint) (*models.Transaction, error)
	UpdateTransaction(id uint, userID uint, req *models.UpdateTransactionRequest) (*models.Transaction, error)
	DeleteTransaction(id uint, userID uint) error
	GetSummary(userID uint, startDate, endDate string) (map[string]interface{}, error)
}

type transactionService struct {
	transactionRepo repository.TransactionRepository
	categoryRepo    repository.CategoryRepository
}

func NewTransactionService(transactionRepo repository.TransactionRepository, categoryRepo repository.CategoryRepository) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
		categoryRepo:    categoryRepo,
	}
}

func (s *transactionService) CreateTransaction(userID uint, req *models.CreateTransactionRequest) (*models.Transaction, error) {
	// Verify that the category belongs to the user
	_, err := s.categoryRepo.GetByID(req.CategoryID, userID)
	if err != nil {
		return nil, errors.New("category not found or does not belong to user")
	}

	transaction := &models.Transaction{
		UserID:      userID,
		CategoryID:  req.CategoryID,
		Amount:      req.Amount,
		Type:        req.Type,
		Description: req.Description,
		Date:        req.Date,
	}

	err = s.transactionRepo.Create(transaction)
	if err != nil {
		return nil, err
	}

	// Fetch the transaction with category details
	return s.transactionRepo.GetByID(transaction.ID, userID)
}

func (s *transactionService) GetTransactions(userID uint, filter *models.TransactionFilter) ([]models.Transaction, error) {
	return s.transactionRepo.GetByUserID(userID, filter)
}

func (s *transactionService) GetTransactionByID(id uint, userID uint) (*models.Transaction, error) {
	return s.transactionRepo.GetByID(id, userID)
}

func (s *transactionService) UpdateTransaction(id uint, userID uint, req *models.UpdateTransactionRequest) (*models.Transaction, error) {
	transaction, err := s.transactionRepo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}

	if req.CategoryID != nil {
		// Verify that the category belongs to the user
		_, err := s.categoryRepo.GetByID(*req.CategoryID, userID)
		if err != nil {
			return nil, errors.New("category not found or does not belong to user")
		}
		transaction.CategoryID = *req.CategoryID
	}

	if req.Amount != nil {
		transaction.Amount = *req.Amount
	}

	if req.Type != nil {
		transaction.Type = *req.Type
	}

	if req.Description != nil {
		transaction.Description = *req.Description
	}

	if req.Date != nil {
		transaction.Date = *req.Date
	}

	err = s.transactionRepo.Update(transaction)
	if err != nil {
		return nil, err
	}

	// Fetch the updated transaction with category details
	return s.transactionRepo.GetByID(transaction.ID, userID)
}

func (s *transactionService) DeleteTransaction(id uint, userID uint) error {
	return s.transactionRepo.Delete(id, userID)
}

func (s *transactionService) GetSummary(userID uint, startDate, endDate string) (map[string]interface{}, error) {
	return s.transactionRepo.GetSummary(userID, startDate, endDate)
}
