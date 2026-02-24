package services

import (
	"time"
	"workouts_bot/src/logger"
	"workouts_bot/src/models"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserService struct {
	database *gorm.DB
}

func NewUserService(database *gorm.DB) *UserService {
	return &UserService{database: database}
}

func (userService *UserService) GetByTelegramID(
	telegramID int64,
) (*models.User, error) {
	var user models.User

	err := userService.database.
		Where("telegram_id = ?", telegramID).
		First(&user).Error
	if err != nil {
		logger.WithFields(logrus.Fields{
			"telegram_id": telegramID,
			"error":       err,
		}).Error("Failed to get user by telegram ID")
		return nil, err
	}

	logger.WithFields(logrus.Fields{
		"user_id":     user.ID,
		"telegram_id": user.TelegramID,
	}).Info("User fetched successfully")

	return &user, nil
}

func (userService *UserService) CreateOrUpdate(user *models.User) error {
	var existingUser models.User

	result := userService.database.
		Where("telegram_id = ?", user.TelegramID).
		First(&existingUser)
	if result.Error == gorm.ErrRecordNotFound {
		err := userService.database.Create(user).Error
		if err != nil {
			logger.WithFields(logrus.Fields{
				"telegram_id": user.TelegramID,
				"error":       err,
			}).Error("Failed to create user")
			return err
		}

		logger.WithFields(logrus.Fields{
			"user_id":     user.ID,
			"telegram_id": user.TelegramID,
		}).Info("User created successfully")

		return nil
	} else if result.Error != nil {
		return result.Error
	}

	user.ID = existingUser.ID
	user.UpdatedAt = time.Now()
	err := userService.database.Save(user).Error
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id":     user.ID,
			"telegram_id": user.TelegramID,
		}).Error("Failed to update user")
		return err
	}

	logger.WithFields(logrus.Fields{
		"user_id": user.ID,
	}).Info("User updated successfully")

	return nil
}

func (userService *UserService) UpdateEquipment(
	userID uuid.UUID,
	equipmentIDs []int,
) error {
	err := userService.database.Model(&models.User{}).
		Where("id = ?", userID).
		Update("equipment_ids", equipmentIDs).
		Error

	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id":       userID,
			"equipment_ids": equipmentIDs,
			"error":         err,
		}).Error("Failed to update equipment")
		return err
	}

	logger.WithFields(logrus.Fields{
		"user_id":       userID,
		"equipment_ids": equipmentIDs,
	}).Info("Equipment updated successfully")

	return nil
}

func (userService *UserService) UpdateGoals(
	userID uuid.UUID,
	goals []string,
) error {
	err := userService.database.Model(&models.User{}).
		Where("id = ?", userID).Update("goals", goals).
		Error

	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"goals":   goals,
			"error":   err,
		}).Error("Failed to update goals")
		return err
	}

	logger.WithFields(logrus.Fields{
		"user_id": userID,
		"goals":   goals,
	}).Info("Goals updated successfully")

	return nil
}
