package services

import (
	"errors"
	"testing"
	"time"

	"github.com/aditherevenger/Budget-Tracker-API/models"
	"github.com/aditherevenger/Budget-Tracker-API/repository"
)

type mockTxnRepo struct {
	CreateFn   func(transaction *models.Transaction) error
	GetByIDFn  func(id uint, userID uint) (*models.Transaction, error)
	ListFn     func(userID uint, filter *models.TransactionFilter) ([]models.Transaction, error)
	UpdateFn   func(transaction *models.Transaction) error
	DeleteFn   func(id uint, userID uint) error
	SummaryFn  func(userID uint, startDate, endDate string) (map[string]interface{}, error)
}

func (m *mockTxnRepo) Create(transaction *models.Transaction) error                                    { return m.CreateFn(transaction) }
func (m *mockTxnRepo) GetByID(id uint, userID uint) (*models.Transaction, error)                       { return m.GetByIDFn(id, userID) }
func (m *mockTxnRepo) GetByUserID(userID uint, filter *models.TransactionFilter) ([]models.Transaction, error) {
	return m.ListFn(userID, filter)
}
func (m *mockTxnRepo) Update(transaction *models.Transaction) error                                    { return m.UpdateFn(transaction) }
func (m *mockTxnRepo) Delete(id uint, userID uint) error                                               { return m.DeleteFn(id, userID) }
func (m *mockTxnRepo) GetSummary(userID uint, startDate, endDate string) (map[string]interface{}, error) {
	return m.SummaryFn(userID, startDate, endDate)
}

var _ repository.TransactionRepository = (*mockTxnRepo)(nil)

type mockCatRepo struct {
	GetByIDFn func(id uint, userID uint) (*models.Category, error)
}

func (m *mockCatRepo) Create(category *models.Category) error { return nil }
func (m *mockCatRepo) GetByUserID(userID uint, filter *models.User) ([]models.Category, error) { return nil, nil }
func (m *mockCatRepo) GetByID(id uint, userID uint) (*models.Category, error) { return m.GetByIDFn(id, userID) }
func (m *mockCatRepo) Update(category *models.Category) error { return nil }
func (m *mockCatRepo) Delete(id uint, userID uint) error { return nil }

var _ repository.CategoryRepository = (*mockCatRepo)(nil)

func TestTransactionService_Create_Success(t *testing.T) {
	mTxn := &mockTxnRepo{ CreateFn: func(transaction *models.Transaction) error { transaction.ID = 1; return nil }, GetByIDFn: func(id uint, userID uint) (*models.Transaction, error) { return &models.Transaction{ID: id, UserID: userID, CategoryID: 2, Amount: 10, Type: models.Expense}, nil } }
	mCat := &mockCatRepo{ GetByIDFn: func(id uint, userID uint) (*models.Category, error) { return &models.Category{ID: id, UserID: userID, Name: "Food"}, nil } }
	svc := NewTransactionService(mTxn, mCat)
	req := &models.CreateTransactionRequest{CategoryID: 2, Amount: 10, Type: models.Expense, Description: "Coffee", Date: time.Now().UTC()}
	tx, err := svc.CreateTransaction(5, req)
	if err != nil { t.Fatalf("create: %v", err) }
	if tx.ID == 0 || tx.CategoryID != 2 { t.Fatalf("unexpected: %+v", tx) }
}

func TestTransactionService_Create_CategoryNotOwned(t *testing.T) {
	mTxn := &mockTxnRepo{}
	mCat := &mockCatRepo{ GetByIDFn: func(id uint, userID uint) (*models.Category, error) { return nil, errors.New("not found") } }
	svc := NewTransactionService(mTxn, mCat)
	_, err := svc.CreateTransaction(5, &models.CreateTransactionRequest{CategoryID: 2, Amount: 10, Type: models.Expense, Date: time.Now().UTC()})
	if err == nil { t.Fatalf("expected error when category not owned") }
}

func TestTransactionService_List_Get_GetByID(t *testing.T) {
	now := time.Now().UTC()
	mTxn := &mockTxnRepo{ ListFn: func(userID uint, filter *models.TransactionFilter) ([]models.Transaction, error) { return []models.Transaction{{ID: 1, UserID: userID}}, nil }, GetByIDFn: func(id uint, userID uint) (*models.Transaction, error) { return &models.Transaction{ID: id, UserID: userID, Date: now}, nil } }
	mCat := &mockCatRepo{ GetByIDFn: func(id uint, userID uint) (*models.Category, error) { return &models.Category{ID: id, UserID: userID}, nil } }
	svc := NewTransactionService(mTxn, mCat)
	items, err := svc.GetTransactions(7, &models.TransactionFilter{})
	if err != nil || len(items) != 1 { t.Fatalf("list: %v len=%d", err, len(items)) }
	got, err := svc.GetTransactionByID(1, 7)
	if err != nil || got.ID != 1 { t.Fatalf("get: %v got=%+v", err, got) }
}

func TestTransactionService_Update_Success(t *testing.T) {
    // stateful mock: keeps current transaction and reflects updates
    current := &models.Transaction{ID: 1, UserID: 7, CategoryID: 2, Amount: 10, Type: models.Expense}
    mTxn := &mockTxnRepo{
        GetByIDFn: func(id uint, userID uint) (*models.Transaction, error) { return current, nil },
        UpdateFn: func(transaction *models.Transaction) error {
            current = transaction
            return nil
        },
    }
    mCat := &mockCatRepo{ GetByIDFn: func(id uint, userID uint) (*models.Category, error) { return &models.Category{ID: id, UserID: userID}, nil } }
    svc := NewTransactionService(mTxn, mCat)
    newAmt := 20.0
    newCat := uint(3)
    req := &models.UpdateTransactionRequest{Amount: &newAmt, CategoryID: &newCat}
    tx, err := svc.UpdateTransaction(1, 7, req)
    if err != nil { t.Fatalf("update: %v", err) }
    if tx.Amount != 20 || tx.CategoryID != 3 { t.Fatalf("unexpected: %+v", tx) }
    got, err := svc.GetTransactionByID(1, 7)
    if err != nil || got.ID != 1 || got.Amount != 20 || got.CategoryID != 3 { t.Fatalf("get: %v got=%+v", err, got) }
}

func TestTransactionService_Update_InvalidCategory(t *testing.T) {
	mTxn := &mockTxnRepo{ GetByIDFn: func(id uint, userID uint) (*models.Transaction, error) { return &models.Transaction{ID: id, UserID: userID, CategoryID: 2, Amount: 10, Type: models.Expense}, nil } }
	mCat := &mockCatRepo{ GetByIDFn: func(id uint, userID uint) (*models.Category, error) { return nil, errors.New("not found") } }
	svc := NewTransactionService(mTxn, mCat)
	newCat := uint(99)
	_, err := svc.UpdateTransaction(1, 7, &models.UpdateTransactionRequest{CategoryID: &newCat})
	if err == nil { t.Fatalf("expected error when category not found/owned") }
}

func TestTransactionService_Delete_And_Summary(t *testing.T) {
	mTxn := &mockTxnRepo{ DeleteFn: func(id uint, userID uint) error { return nil }, SummaryFn: func(userID uint, startDate, endDate string) (map[string]interface{}, error) { return map[string]interface{}{"total_income": 100.0, "total_expense": 50.0, "net_balance": 50.0}, nil } }
	mCat := &mockCatRepo{ GetByIDFn: func(id uint, userID uint) (*models.Category, error) { return &models.Category{ID: id, UserID: userID}, nil } }
	svc := NewTransactionService(mTxn, mCat)
	if err := svc.DeleteTransaction(2, 7); err != nil { t.Fatalf("delete: %v", err) }
	sum, err := svc.GetSummary(7, "", "")
	if err != nil { t.Fatalf("summary: %v", err) }
	if sum["net_balance"].(float64) != 50.0 { t.Fatalf("unexpected summary: %+v", sum) }
}
