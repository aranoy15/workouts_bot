package services

import (
	"workouts_bot/pkg/logger"
	"workouts_bot/src/database/models"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type WorkoutsService struct {
	Database *gorm.DB
}

func NewWorkoutsService(database *gorm.DB) *WorkoutsService {
	return &WorkoutsService{Database: database}
}

func (workoutsService *WorkoutsService) CreateWorkout(
	userID uuid.UUID,
	name string,
	workoutType string,
) (*models.Workout, error) {
	workout := models.Workout{
		UserID:      userID,
		Name:        name,
		Status:      "planned",
		WorkoutType: workoutType,
	}

	err := workoutsService.Database.Create(&workout).Error
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id":      userID,
			"workout_type": workoutType,
			"error":        err,
		}).Error("Failed to create workout")
		return nil, err
	}

	logger.WithFields(logrus.Fields{
		"workout_id": workout.ID,
		"user_id":    userID,
	}).Info("Workout created successfully")

	return &workout, nil
}

func (workoutsService *WorkoutsService) GetUserWorkouts(
	userID uuid.UUID,
) ([]models.Workout, error) {
	var workouts []models.Workout

	err := workoutsService.Database.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&workouts).
		Error

	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   err,
		}).Error("Failed to get user workouts")
		return nil, err
	}

	return workouts, nil
}

func (workoutsService *WorkoutsService) GetWorkoutByID(
	workoutID uint64,
) (*models.Workout, error) {
	var workout models.Workout

	err := workoutsService.Database.
		Preload("Exercises.Exercise").
		Preload("Exercises.Sets").
		First(&workout, workoutID).
		Error

	if err != nil {
		logger.WithFields(logrus.Fields{
			"workout_id": workoutID,
			"error":      err,
		}).Error("Failed to get workout by ID")
		return nil, err
	}

	return &workout, nil
}

func (workoutsService *WorkoutsService) StartWorkout(workoutID uint64) error {
	err := workoutsService.Database.
		Model(&models.Workout{}).
		Where("id = ?", workoutID).
		Update("status", "in_progress").
		Error

	if err != nil {
		logger.WithFields(logrus.Fields{
			"workout_id": workoutID,
			"error":      err,
		}).Error("Failed to start workout")
		return err
	}

	logger.WithFields(logrus.Fields{
		"workout_id": workoutID,
	}).Info("Workout started")

	return nil
}

func (workoutsService *WorkoutsService) CompleteWorkout(workoutID uint64) error {
	err := workoutsService.Database.
		Model(&models.Workout{}).
		Where("id = ?", workoutID).
		Updates(map[string]interface{}{
			"status":       "completed",
			"completed_at": "NOW()",
		}).
		Error

	if err != nil {
		logger.WithFields(logrus.Fields{
			"workout_id": workoutID,
			"error":      err,
		}).Error("Failed to complete workout")
		return err
	}

	logger.WithFields(logrus.Fields{
		"workout_id": workoutID,
	}).Info("Workout completed")

	return nil
}

func (workoutsService *WorkoutsService) AddExerciseToWorkout(
	workoutID uint,
	exerciseID uint,
	orderIndex int,
	setsCount int,
	repsCount int,
	restSeconds int,
	weightKg float64,
) error {
	workoutExercise := models.WorkoutExercise{
		WorkoutID:   workoutID,
		ExerciseID:  exerciseID,
		OrderIndex:  orderIndex,
		SetsCount:   setsCount,
		RepsCount:   repsCount,
		RestSeconds: restSeconds,
		WeightKg:    weightKg,
	}

	err := workoutsService.Database.Create(&workoutExercise).Error
	if err != nil {
		logger.WithFields(logrus.Fields{
			"workout_id":  workoutID,
			"exercise_id": exerciseID,
			"order_index": orderIndex,
			"error":       err,
		}).Error("Failed to add exercise to workout")
		return err
	}

	logger.WithFields(logrus.Fields{
		"workout_exercise_id": workoutExercise.ID,
		"workout_id":          workoutID,
		"exercise_id":         exerciseID,
	}).Info("Exercise added to workout")

	return nil
}

func (workoutsService *WorkoutsService) DeleteWorkout(workoutID uint) error {
	err := workoutsService.Database.
		Where("id = ?", workoutID).
		Delete(&models.Workout{}).
		Error

	if err != nil {
		logger.WithFields(logrus.Fields{
			"workout_id": workoutID,
			"error":      err,
		}).Error("Failed to delete workout")
		return err
	}

	logger.WithFields(logrus.Fields{
		"workout_id": workoutID,
	}).Info("Workout deleted")

	return nil
}
