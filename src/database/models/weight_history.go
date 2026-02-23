package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WeightHistory struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	ExerciseID uint      `gorm:"not null" json:"exercise_id"`
	WeightKg   float64   `gorm:"not null" json:"weight_kg"`
	RepsCount  int       `gorm:"not null" json:"reps_count"`
	RecordedAt time.Time `json:"recorded_at"`
	User       User      `gorm:"foreignKey:UserID" json:"user"`
	Exercise   Exercise  `gorm:"foreignKey:ExerciseID" json:"exercise"`
}

func (WeightHistory) TableName() string {
	return "workouts.weight_histories"
}

func (wh *WeightHistory) BeforeCreate(tx *gorm.DB) error {
	if wh.RecordedAt.IsZero() {
		wh.RecordedAt = time.Now()
	}
	return nil
}
