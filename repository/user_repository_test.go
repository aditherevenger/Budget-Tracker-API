package repository

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"github.com/aditherevenger/Budget-Tracker-API/database"
	"github.com/aditherevenger/Budget-Tracker-API/models"
)

func setupTestDB(t *testing.T) *gorm.DB {
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

func TestUserRepository_CRUD(t *testing.T) {
	setupTestDB(t)
	repo := NewUserRepository()

	u := &models.User{Email: "jane@example.com", Password: "hash", FirstName: "Jane", LastName: "Doe"}
	if err := repo.Create(u); err != nil {
		t.Fatalf("create user: %v", err)
	}
	if u.ID == 0 { t.Fatalf("expected user ID to be set") }

	gotByEmail, err := repo.GetByEmail("jane@example.com")
	if err != nil { t.Fatalf("get by email: %v", err) }
	if gotByEmail.ID != u.ID { t.Fatalf("expected same ID, got %d vs %d", gotByEmail.ID, u.ID) }

	gotByID, err := repo.GetByID(u.ID)
	if err != nil { t.Fatalf("get by id: %v", err) }
	if gotByID.Email != "jane@example.com" { t.Fatalf("unexpected email: %s", gotByID.Email) }

	gotByID.FirstName = "Janet"
	if err := repo.Update(gotByID); err != nil { t.Fatalf("update: %v", err) }

	reloaded, err := repo.GetByID(u.ID)
	if err != nil { t.Fatalf("reload: %v", err) }
	if reloaded.FirstName != "Janet" { t.Fatalf("expected updated first name, got %s", reloaded.FirstName) }

	if err := repo.Delete(reloaded); err != nil { t.Fatalf("delete: %v", err) }
	if _, err := repo.GetByID(u.ID); err == nil {
		t.Fatalf("expected error when fetching deleted user")
	}
}
