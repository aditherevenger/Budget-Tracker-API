package repository

import (
	"github.com/aditherevenger/Budget-Tracker-API/database"
	"github.com/aditherevenger/Budget-Tracker-API/models"
)

type TransactionRepository interface {
	Create(transaction *models.Transaction) error
	GetByID(id uint, userID uint) (*models.Transaction, error)
	GetByUserID(userID uint, filter *models.TransactionFilter) ([]models.Transaction, error)
	Update(transaction *models.Transaction) error
	Delete(id uint, userID uint) error
	GetSummary(userID uint, startDate, endDate string) (map[string]interface{}, error)
}

type transactionRepository struct{}

func NewTransactionRepository() TransactionRepository {
	return &transactionRepository{}
}

func (r *transactionRepository) Create(transaction *models.Transaction) error {
	return database.DB.Create(transaction).Error
}

func (r *transactionRepository) GetByID(id uint, userID uint) (*models.Transaction, error) {
	var transaction models.Transaction
	err := database.DB.Preload("Category").Where("id = ? AND user_id = ?", id, userID).First(&transaction).Error
	return &transaction, err
}

func (r *transactionRepository) GetByUserID(userID uint, filter *models.TransactionFilter) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := database.DB.Preload("Category").Where("user_id = ?", userID)

	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}

	if filter.CategoryID != 0 {
		query = query.Where("category_id = ?", filter.CategoryID)
	}

	if !filter.StartDate.IsZero() {
		query = query.Where("date >= ?", filter.StartDate)
	}

	if !filter.EndDate.IsZero() {
		query = query.Where("date <= ?", filter.EndDate)
	}

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	err := query.Order("date DESC").Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) Update(transaction *models.Transaction) error {
	return database.DB.Save(transaction).Error
}

func (r *transactionRepository) Delete(id uint, userID uint) error {
	return database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Transaction{}).Error
}

func (r *transactionRepository) GetSummary(userID uint, startDate, endDate string) (map[string]interface{}, error) {
	var result struct {
		TotalIncome  float64
		TotalExpense float64
		NetBalance   float64
	}
	query := database.DB.Model(&models.Transaction{}).Where("user_id = ?", userID)
	if startDate != "" {
		query = query.Where("date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("date <= ?", endDate)
	}

	err := query.Select(
		"SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END) as total_income",
		"SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END) as total_expense",
		"SUM(CASE WHEN type = 'income' THEN amount ELSE -amount END) as net_balance",
	).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_income":  result.TotalIncome,
		"total_expense": result.TotalExpense,
		"net_balance":   result.NetBalance,
	}, nil
}
