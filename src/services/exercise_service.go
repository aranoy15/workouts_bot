package services

import (
	"workouts_bot/src/logger"
	"workouts_bot/src/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ExerciseService struct {
	database *gorm.DB
}

func NewExerciseService(database *gorm.DB) *ExerciseService {
	return &ExerciseService{database: database}
}

func (exerciseService *ExerciseService) GetAll() ([]models.Exercise, error) {
	var exercises []models.Exercise

	err := exerciseService.database.Find(&exercises).Error
	if err != nil {
		logger.Error("Failed to get all exercises", err)
		return nil, err
	}

	return exercises, nil
}

func (exerciseService *ExerciseService) GetByMuscleGroups(
	muscleGroups []string,
) ([]models.Exercise, error) {
	var exercises []models.Exercise

	err := exerciseService.database.
		Where("? = ANY(muscle_groups)", muscleGroups).
		Find(&exercises).
		Error

	if err != nil {
		logger.WithFields(logrus.Fields{
			"muscle_groups": muscleGroups,
			"error":         err,
		}).Error("Failed to get exercises by muscle groups")
		return nil, err
	}

	return exercises, nil
}

func (exerciseService *ExerciseService) GetByEquipment(
	equipmentIDs []int,
) ([]models.Exercise, error) {
	var exercises []models.Exercise

	if len(equipmentIDs) == 0 {
		err := exerciseService.database.
			Where("equipment_ids = ?", equipmentIDs).
			Find(&exercises).
			Error

		logger.WithFields(logrus.Fields{
			"equipment_ids": equipmentIDs,
			"error":         err,
		}).Error("Failed to get exercises by equipment")

		return exercises, err
	}

	err := exerciseService.database.
		Where("equipment_ids <@ ?", equipmentIDs).
		Find(&exercises).
		Error

	if err != nil {
		logger.WithFields(logrus.Fields{
			"equipment_ids": equipmentIDs,
			"error":         err,
		}).Error("Failed to get exercises by equipment")
		return nil, err
	}

	return exercises, nil
}

func (exerciseService *ExerciseService) GetByDifficulty(
	maxDifficulty int,
) ([]models.Exercise, error) {
	var exercises []models.Exercise

	err := exerciseService.database.
		Where("difficulty <= ?", maxDifficulty).
		Find(&exercises).
		Error

	if err != nil {
		logger.WithFields(logrus.Fields{
			"max_difficulty": maxDifficulty,
			"error":          err,
		}).Error("Failed to get exercises by difficulty")
		return nil, err
	}

	return exercises, nil
}

func (exerciseService *ExerciseService) GetByCategory(
	category string,
) ([]models.Exercise, error) {
	var exercises []models.Exercise

	err := exerciseService.database.
		Where("category = ?", category).
		Find(&exercises).
		Error

	if err != nil {
		logger.WithFields(logrus.Fields{
			"category": category,
			"error":    err,
		}).Error("Failed to get exercises by category")
		return nil, err
	}

	return exercises, nil
}

func (exerciseService *ExerciseService) GetByID(
	exerciseID uint,
) (*models.Exercise, error) {
	var exercise models.Exercise

	err := exerciseService.database.
		First(&exercise, exerciseID).
		Error

	if err != nil {
		logger.WithFields(logrus.Fields{
			"exercise_id": exerciseID,
			"error":       err,
		}).Error("Failed to get exercise by ID")
		return nil, err
	}

	return &exercise, nil
}
