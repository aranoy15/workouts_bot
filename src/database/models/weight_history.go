package models

import (
	"time"

	"gorm.io/gorm"
)

type WeightHistory struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null" json:"user_id"`
	ExerciseID uint      `gorm:"not null" json:"exercise_id"`
	WeightKg   float64   `gorm:"not null" json:"weight_kg"`
	RepsCount  int       `gorm:"not null" json:"reps_count"`
	RecordedAt time.Time `json:"recorded_at"`
	User       User      `gorm:"foreignKey:UserID" json:"user"`
	Exercise   Exercise  `gorm:"foreignKey:ExerciseID" json:"exercise"`
}

func (wh *WeightHistory) BeforeCreate(tx *gorm.DB) error {
	if wh.RecordedAt.IsZero() {
		wh.RecordedAt = time.Now()
	}
	return nil
}
