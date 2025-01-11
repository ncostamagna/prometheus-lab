package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User model
type Product struct {
	ID          string         `json:"id" gorm:"type:char(36);not null;primary_key;unique_index"`
	Name        string         `json:"name" gorm:"type:char(50);not null"`
	Description string         `json:"description" gorm:"type:Blob"`
	Price       float64        `json:"price"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-"`
}

func (c *Product) BeforeCreate(tx *gorm.DB) error {

	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}
