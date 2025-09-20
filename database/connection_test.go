package database

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"github.com/aditherevenger/Budget-Tracker-API/models"
)

func TestMigrate_WithInMemorySQLite(t *testing.T) {
	// override global DB with in-memory sqlite
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil { t.Fatalf("open sqlite: %v", err) }
	DB = db

	// call migrate
	Migrate()

	// check tables exist by trying basic operations
	if err := DB.Create(&models.User{Email: "a@b.com", Password: "x", FirstName: "A", LastName: "B"}).Error; err != nil {
		t.Fatalf("user create failed after migrate: %v", err)
	}
}
