package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	Type        string `gorm:"not null"` // e.g., "income" or "expense"
	UserID      uint   `gorm:"not null"` // To associate category with a user
	// Add other category fields as needed
}
