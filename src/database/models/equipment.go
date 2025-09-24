package models

import (
	"time"

	"gorm.io/gorm"
)

type Equipment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Type      string    `gorm:"not null" json:"type"`
	IsHomeGym bool      `gorm:"default:false" json:"is_home_gym"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (e *Equipment) BeforeCreate(tx *gorm.DB) error {
	if e.Type == "" {
		e.Type = "general"
	}
	return nil
}
