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
    // Income
    incomeQuery := database.DB.Model(&models.Transaction{}).Where("user_id = ? AND type = ?", userID, models.Income)
    if startDate != "" {
        incomeQuery = incomeQuery.Where("date >= ?", startDate)
    }
    if endDate != "" {
        incomeQuery = incomeQuery.Where("date <= ?", endDate)
    }
    var totalIncome float64
    if err := incomeQuery.Select("COALESCE(SUM(amount), 0)").Scan(&totalIncome).Error; err != nil {
        return nil, err
    }

    // Expense
    expenseQuery := database.DB.Model(&models.Transaction{}).Where("user_id = ? AND type = ?", userID, models.Expense)
    if startDate != "" {
        expenseQuery = expenseQuery.Where("date >= ?", startDate)
    }
    if endDate != "" {
        expenseQuery = expenseQuery.Where("date <= ?", endDate)
    }
    var totalExpense float64
    if err := expenseQuery.Select("COALESCE(SUM(amount), 0)").Scan(&totalExpense).Error; err != nil {
        return nil, err
    }

    net := totalIncome - totalExpense

    return map[string]interface{}{
        "total_income":  totalIncome,
        "total_expense": totalExpense,
        "net_balance":   net,
    }, nil
}
