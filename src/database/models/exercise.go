package models

import (
	"time"

	"gorm.io/gorm"
)

type Exercise struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Name            string    `gorm:"not null" json:"name"`
	Description     string    `json:"description"`
	Category        string    `json:"category"`
	MuscleGroups    []string  `gorm:"type:text[]" json:"muscle_groups"`
	Difficulty      int       `gorm:"default:1" json:"difficulty"`
	DurationMinutes int       `json:"duration_minutes"`
	ImagePath       string    `json:"image_path"`
	VideoPath       string    `json:"video_path"`
	EquipmentIDs    []int     `gorm:"type:integer[]" json:"equipment_ids"`
	Instructions    []string  `gorm:"type:text[]" json:"instructions"`
	CommonMistakes  []string  `gorm:"type:text[]" json:"common_mistakes"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (exercise *Exercise) BeforeCreate(tx *gorm.DB) error {
	if exercise.Difficulty == 0 {
		exercise.Difficulty = 1
	}
	return nil
}
