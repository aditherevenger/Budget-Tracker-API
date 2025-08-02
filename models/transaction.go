package models

import (
	"gorm.io/gorm"
	"time"
)

type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

type Transaction struct {
	ID          uint            `json:"id" gorm:"primaryKey"`
	UserID      uint            `json:"user_id" gorm:"not null"`
	CategoryID  uint            `json:"category_id" gorm:"not null"`
	Amount      float64         `json:"amount" gorm:"not null"`
	Type        TransactionType `json:"type" gorm:"not null"`
	Description string          `json:"description"`
	Date        time.Time       `json:"date" gorm:"not null"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `json:"-" gorm:"index"`

	// Relationships
	User     User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Category Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}

type CreateTransactionRequest struct {
	CategoryID  uint            `json:"category_id" binding:"required"`
	Amount      float64         `json:"amount" binding:"required,gt=0"`
	Type        TransactionType `json:"type" binding:"required,oneof=income expense"`
	Description string          `json:"description"`
	Date        time.Time       `json:"date" binding:"required"`
}

type UpdateTransactionRequest struct {
	CategoryID  *uint            `json:"category_id,omitempty"`
	Amount      *float64         `json:"amount,omitempty" binding:"omitempty,gt=0"`
	Type        *TransactionType `json:"type,omitempty" binding:"omitempty,oneof=income expense"`
	Description *string          `json:"description,omitempty"`
	Date        *time.Time       `json:"date,omitempty"`
}

type TransactionFilter struct {
	Type       TransactionType `form:"type"`
	CategoryID uint            `form:"category_id"`
	StartDate  time.Time       `form:"start_date"`
	EndDate    time.Time       `form:"end_date"`
	Limit      int             `form:"limit"`
	Offset     int             `form:"offset"`
}
