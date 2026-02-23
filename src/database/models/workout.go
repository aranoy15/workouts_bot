package models

import (
	"time"

	"github.com/google/uuid"
)

func (WorkoutType) TableName() string {
	return "workouts.workout_types"
}

type WorkoutType struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

func (Workout) TableName() string {
	return "workouts.workouts"
}

type Workout struct {
	ID              uint              `gorm:"primaryKey" json:"id"`
	UserID          uuid.UUID         `gorm:"type:uuid;not null" json:"user_id"`
	Name            string            `json:"name"`
	WorkoutTypeID   uint              `gorm:"not null" json:"workout_type_id"`
	DurationMinutes int               `json:"duration_minutes"`
	Status          string            `gorm:"default:'planned'" json:"status"`
	WorkoutType     string            `json:"workout_type"`
	CreatedAt       time.Time         `json:"created_at"`
	CompletedAt     *time.Time        `json:"completed_at"`
	User            User              `gorm:"foreignKey:UserID" json:"user"`
	Exercises       []WorkoutExercise `gorm:"foreignKey:WorkoutID" json:"exercises"`
}

func (WorkoutExercise) TableName() string {
	return "workouts.workout_exercises"
}

type WorkoutExercise struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	WorkoutID   uint       `gorm:"not null" json:"workout_id"`
	ExerciseID  uint       `gorm:"not null" json:"exercise_id"`
	OrderIndex  int        `json:"order_index"`
	SetsCount   int        `json:"sets_count"`
	RepsCount   int        `json:"reps_count"`
	RestSeconds int        `json:"rest_seconds"`
	WeightKg    float64    `json:"weight_kg"`
	CompletedAt *time.Time `json:"completed_at"`
	Workout     Workout    `gorm:"foreignKey:WorkoutID" json:"workout"`
	Exercise    Exercise   `gorm:"foreignKey:ExerciseID" json:"exercise"`
	Sets        []Set      `gorm:"foreignKey:WorkoutExerciseID" json:"sets"`
}

func (Set) TableName() string {
	return "workouts.sets"
}

type Set struct {
	ID                uint            `gorm:"primaryKey" json:"id"`
	WorkoutExerciseID uint            `gorm:"not null" json:"workout_exercise_id"`
	SetNumber         int             `json:"set_number"`
	WeightKg          float64         `json:"weight_kg"`
	RepsDone          int             `json:"reps_done"`
	RestTakenSeconds  int             `json:"rest_taken_seconds"`
	CompletedAt       time.Time       `json:"completed_at"`
	WorkoutExercise   WorkoutExercise `gorm:"foreignKey:WorkoutExerciseID" json:"workout_exercise"`
}
