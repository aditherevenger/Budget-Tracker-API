package repository

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"github.com/aditherevenger/Budget-Tracker-API/database"
	"github.com/aditherevenger/Budget-Tracker-API/models"
)

func setupTestDBTransaction(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite in-memory: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Category{}, &models.Transaction{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	database.DB = db
	return db
}

func TestTransactionRepository_CRUD_List_Filters_Summary(t *testing.T) {
	setupTestDBTransaction(t)
	trepo := NewTransactionRepository()
	crepo := NewCategoryRepository()
	urepo := NewUserRepository()

	// create a user and categories
	u := &models.User{Email: "owner@example.com", Password: "hash", FirstName: "Own", LastName: "Er"}
	if err := urepo.Create(u); err != nil { t.Fatalf("create user: %v", err) }
	catFood := &models.Category{UserID: u.ID, Name: "Food"}
	if err := crepo.Create(catFood); err != nil { t.Fatalf("create category: %v", err) }
	catSalary := &models.Category{UserID: u.ID, Name: "Salary"}
	if err := crepo.Create(catSalary); err != nil { t.Fatalf("create category: %v", err) }

	d1 := time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC)
	d2 := time.Date(2025, 9, 10, 0, 0, 0, 0, time.UTC)

	// create transactions
	tx1 := &models.Transaction{UserID: u.ID, CategoryID: catFood.ID, Amount: 50, Type: models.Expense, Description: "Lunch", Date: d1}
	if err := trepo.Create(tx1); err != nil { t.Fatalf("create tx1: %v", err) }
	tx2 := &models.Transaction{UserID: u.ID, CategoryID: catSalary.ID, Amount: 1000, Type: models.Income, Description: "Pay", Date: d2}
	if err := trepo.Create(tx2); err != nil { t.Fatalf("create tx2: %v", err) }

	// get by id (with preload)
	got, err := trepo.GetByID(tx1.ID, u.ID)
	if err != nil { t.Fatalf("get by id: %v", err) }
	if got.CategoryID != catFood.ID || got.Category.Name != "Food" { t.Fatalf("expected preload category Food, got %+v", got.Category) }

	// list by user no filters
	items, err := trepo.GetByUserID(u.ID, &models.TransactionFilter{})
	if err != nil { t.Fatalf("list: %v", err) }
	if len(items) != 2 { t.Fatalf("expected 2 transactions, got %d", len(items)) }

	// filter by type
	items, err = trepo.GetByUserID(u.ID, &models.TransactionFilter{Type: models.Expense})
	if err != nil { t.Fatalf("list by type: %v", err) }
	if len(items) != 1 || items[0].Type != models.Expense { t.Fatalf("expected 1 expense, got %+v", items) }

	// filter by category
	items, err = trepo.GetByUserID(u.ID, &models.TransactionFilter{CategoryID: catSalary.ID})
	if err != nil { t.Fatalf("list by category: %v", err) }
	if len(items) != 1 || items[0].CategoryID != catSalary.ID { t.Fatalf("expected 1 salary tx, got %+v", items) }

	// filter by date range
	items, err = trepo.GetByUserID(u.ID, &models.TransactionFilter{StartDate: d2.Add(-time.Hour), EndDate: d2.Add(time.Hour)})
	if err != nil { t.Fatalf("list by date: %v", err) }
	if len(items) != 1 || !items[0].Date.Equal(d2) { t.Fatalf("expected 1 tx at d2, got %+v", items) }

	// update
	got.Amount = 60
	if err := trepo.Update(got); err != nil { t.Fatalf("update: %v", err) }
	reloaded, err := trepo.GetByID(tx1.ID, u.ID)
	if err != nil { t.Fatalf("reload: %v", err) }
	if reloaded.Amount != 60 { t.Fatalf("expected amount 60, got %v", reloaded.Amount) }

	// delete
	if err := trepo.Delete(tx2.ID, u.ID); err != nil { t.Fatalf("delete: %v", err) }
	items, err = trepo.GetByUserID(u.ID, &models.TransactionFilter{})
	if err != nil { t.Fatalf("re-list: %v", err) }
	if len(items) != 1 { t.Fatalf("expected 1 tx after delete, got %d", len(items)) }

	// summary (no date filter; function takes strings)
	sum, err := trepo.GetSummary(u.ID, "", "")
	if err != nil { t.Fatalf("summary: %v", err) }
	income := sum["total_income"].(float64)
	expense := sum["total_expense"].(float64)
	net := sum["net_balance"].(float64)
	if income < 0 || expense < 0 { t.Fatalf("expected non-negative totals, got income=%v expense=%v", income, expense) }
	if (income - expense) != net { t.Fatalf("expected net = income - expense, got income=%v expense=%v net=%v", income, expense, net) }
}
