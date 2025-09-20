package repository

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"github.com/aditherevenger/Budget-Tracker-API/database"
	"github.com/aditherevenger/Budget-Tracker-API/models"
)

func setupTestDBCategory(t *testing.T) *gorm.DB {
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

func TestCategoryRepository_CRUD_And_ListByUser(t *testing.T) {
	setupTestDBCategory(t)
	crepo := NewCategoryRepository()
	urepo := NewUserRepository()

	// create a user
	u := &models.User{Email: "owner@example.com", Password: "hash", FirstName: "Own", LastName: "Er"}
	if err := urepo.Create(u); err != nil { t.Fatalf("create user: %v", err) }

	// create categories
	c1 := &models.Category{UserID: u.ID, Name: "Food", Description: "Meals", Color: "#F00"}
	if err := crepo.Create(c1); err != nil { t.Fatalf("create cat1: %v", err) }
	c2 := &models.Category{UserID: u.ID, Name: "Rent", Description: "House", Color: "#0F0"}
	if err := crepo.Create(c2); err != nil { t.Fatalf("create cat2: %v", err) }

	// list by user
	cats, err := crepo.GetByUserID(u.ID, nil)
	if err != nil { t.Fatalf("list by user: %v", err) }
	if len(cats) != 2 { t.Fatalf("expected 2 categories, got %d", len(cats)) }

	// get by id scoped to user
	got, err := crepo.GetByID(c1.ID, u.ID)
	if err != nil { t.Fatalf("get by id: %v", err) }
	if got.Name != "Food" { t.Fatalf("unexpected name: %s", got.Name) }

	// update
	got.Color = "#00AAFF"
	if err := crepo.Update(got); err != nil { t.Fatalf("update: %v", err) }
	reloaded, err := crepo.GetByID(c1.ID, u.ID)
	if err != nil { t.Fatalf("reload: %v", err) }
	if reloaded.Color != "#00AAFF" { t.Fatalf("expected updated color, got %s", reloaded.Color) }

	// delete
	if err := crepo.Delete(c2.ID, u.ID); err != nil { t.Fatalf("delete: %v", err) }
	cats, err = crepo.GetByUserID(u.ID, nil)
	if err != nil { t.Fatalf("re-list: %v", err) }
	if len(cats) != 1 { t.Fatalf("expected 1 category after delete, got %d", len(cats)) }
}
