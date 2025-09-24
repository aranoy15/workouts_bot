package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	TelegramID   int64     `gorm:"uniqueIndex;not null" json:"telegram_id"`
	Username     string    `json:"username"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	EquipmentIDs []int     `gorm:"type:integer[]" json:"equipment_ids"`
	Goals        []string  `gorm:"type:text[]" json:"goals"`
	Experience   int       `gorm:"default:1" json:"experience"`
	Limitations  []string  `gorm:"type:text[]" json:"limitations"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	if user.Experience == 0 {
		user.Experience = 1
	}
	return nil
}
