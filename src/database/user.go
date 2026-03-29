package database

import (
	"time"
	"workouts_bot/src/models"

	"workouts_bot/src/logger"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func GetUserByTelegramID(telegramID int64, db *gorm.DB) (*models.User, error) {
	var user models.User

	err := db.Where("telegram_id = ?", telegramID).First(&user).Error
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
	}).Info("User fetched successfully by telegram ID")

	return &user, nil
}

func UpsertUser(user *models.User, db *gorm.DB) error {
	existingUser, err := GetUserByTelegramID(user.TelegramID, db)

	if err == gorm.ErrRecordNotFound {
		logger.WithFields(logrus.Fields{
			"telegram_id": user.TelegramID,
		}).Info("User not found, creating new user")
		return db.Create(user).Error
	} else if err != nil {
		return err
	}

	existingUser.TelegramID = user.TelegramID
	existingUser.Experience = user.Experience
	existingUser.Username = user.Username
	existingUser.FirstName = user.FirstName
	existingUser.LastName = user.LastName
	existingUser.UpdatedAt = time.Now()
	err = db.Save(&existingUser).Error
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id":     existingUser.ID,
			"telegram_id": existingUser.TelegramID,
			"error":       err,
		}).Error("Failed to update user")
		return err
	}

	logger.WithFields(logrus.Fields{
		"user_id":     existingUser.ID,
		"telegram_id": existingUser.TelegramID,
	}).Info("User updated successfully")
	return nil
}
