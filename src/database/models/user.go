package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Goal string

const (
	GoalMuscleGain Goal = "muscle_gain"
	GoalWeightLoss Goal = "weight_loss"
	GoalStrength   Goal = "strength"
)

type GoalsSlice []Goal

func (goalSlice GoalsSlice) Value() (driver.Value, error) {
	if len(goalSlice) == 0 {
		return "[]", nil
	}

	return json.Marshal(goalSlice)
}

func (goalSlice *GoalsSlice) Scan(value interface{}) error {
	if value == nil {
		*goalSlice = GoalsSlice{}
		return nil
	}

	var bytes []byte
	switch item := value.(type) {
	case string:
		bytes = []byte(item)
	case []byte:
		bytes = item
	default:
		return fmt.Errorf("cannot scan %T into GoalsSlice", value)
	}

	return json.Unmarshal(bytes, goalSlice)
}

type User struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	TelegramID   int64      `gorm:"uniqueIndex;not null" json:"telegram_id"`
	Username     string     `json:"username"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	EquipmentIDs []int      `gorm:"type:integer[]" json:"equipment_ids"`
	Goals        GoalsSlice `gorm:"type:jsonb" json:"goals"`
	Experience   int        `gorm:"default:1" json:"experience"`
	Limitations  []string   `gorm:"type:text[]" json:"limitations"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	if user.Experience == 0 {
		user.Experience = 1
	}
	return nil
}
