package services

import (
	"time"
	"workouts_bot/pkg/logger"
	"workouts_bot/src/database/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProgressService struct {
	database *gorm.DB
}

type NextSetInfo struct {
	SetNumber    int
	TotalSets    int
	ExerciseName string
	WeightKg     float64
	RepsCount    int
	RestSeconds  int
}

func NewProgressService(database *gorm.DB) *ProgressService {
	return &ProgressService{database: database}
}

func (progressService *ProgressService) RecordSet(set *models.Set) error {
	set.CompletedAt = time.Now()

	err := progressService.database.Create(set).Error
	if err != nil {
		logger.WithFields(logrus.Fields{
			"workout_exercise_id": set.WorkoutExerciseID,
			"set_number":          set.SetNumber,
			"weight":              set.WeightKg,
			"reps":                set.RepsDone,
			"error":               err,
		}).Error("Failed to record set")
		return err
	}

	if err := progressService.updateWeightHistory(set); err != nil {
		logger.Error("Failed to update weight history: ", err)
	}

	logger.WithFields(logrus.Fields{
		"set_id": set.ID,
		"weight": set.WeightKg,
		"reps":   set.RepsDone,
	}).Info("Set recorded successfully")

	return nil
}

func (progressService *ProgressService) updateWeightHistory(
	set *models.Set,
) error {
	var workoutExercise models.WorkoutExercise
	err := progressService.database.
		Preload("Workout").
		First(&workoutExercise, set.WorkoutExerciseID).
		Error

	if err != nil {
		return err
	}

	weightHistory := models.WeightHistory{
		UserID:     workoutExercise.Workout.UserID,
		ExerciseID: workoutExercise.ExerciseID,
		WeightKg:   set.WeightKg,
		RepsCount:  set.RepsDone,
		RecordedAt: set.CompletedAt,
	}

	return progressService.database.Create(&weightHistory).Error
}

func (progressService *ProgressService) GetNextSet(
	workoutExerciseID uint,
) (*NextSetInfo, error) {
	var workoutExercise models.WorkoutExercise
	err := progressService.database.
		Preload("Exercise").
		Preload("Sets").
		First(&workoutExercise, workoutExerciseID).
		Error

	if err != nil {
		return nil, err
	}

	completedSets := len(workoutExercise.Sets)
	nextSetNumber := completedSets + 1

	if nextSetNumber > workoutExercise.SetsCount {
		return nil, nil
	}

	nextSet := &NextSetInfo{
		SetNumber:    nextSetNumber,
		TotalSets:    workoutExercise.SetsCount,
		ExerciseName: workoutExercise.Exercise.Name,
		WeightKg:     workoutExercise.WeightKg,
		RepsCount:    workoutExercise.RepsCount,
		RestSeconds:  workoutExercise.RestSeconds,
	}

	return nextSet, nil
}
