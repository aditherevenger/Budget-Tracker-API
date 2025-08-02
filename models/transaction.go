package models

import (
	"gorm.io/gorm"
	"time"
)

type Transaction struct {
	gorm.Model
	Amount      float64   `gorm:"not null"`
	Date        time.Time `gorm:"not null"`
	Description string
	Type        string    `gorm:"not null"` // e.g., "income" or "expense"
	CategoryID  uint      `gorm:"not null"` // Reference to Category
	UserID      uint      `gorm:"not null"` // Reference to User
	// Add other transaction fields as needed
}
