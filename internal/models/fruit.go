package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Fruit struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Price     float64   `gorm:"not null" json:"price"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (f *Fruit) BeforeCreate(tx *gorm.DB) (err error) { f.ID = uuid.New(); return nil }
