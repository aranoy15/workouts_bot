package callbacks

import (
	"strings"
	"workouts_bot/pkg/logger"
	"workouts_bot/src/bot/handlers"
	"workouts_bot/src/database/models"
	"workouts_bot/src/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const GoalCallbackType = "goal"

type GoalHandler struct {
	bot         *tgbotapi.BotAPI
	userService *services.UserService
}

func NewGoalHandler(bot *tgbotapi.BotAPI, database *gorm.DB) *GoalHandler {
	return &GoalHandler{
		bot:         bot,
		userService: services.NewUserService(database),
	}
}

func (h *GoalHandler) Handle(update tgbotapi.Update) error {
	callbackQuery := update.CallbackQuery
	userID := callbackQuery.From.ID
	chatID := callbackQuery.Message.Chat.ID
	messageID := callbackQuery.Message.MessageID
	data := callbackQuery.Data
	parts := strings.Split(data, ":")

	logger.WithFields(logrus.Fields{
		"user_id":    userID,
		"chat_id":    chatID,
		"message_id": messageID,
		"data":       data,
	}).Info("Goal callback received")

	if len(parts) < 2 {
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"data":    data,
		}).Error("Invalid goal callback format")
		handlers.SendErrorMessage(h.bot, chatID, "Неверный формат команды")
		return nil
	}

	goal := parts[1]

	logger.WithFields(logrus.Fields{
		"user_id": userID,
		"chat_id": chatID,
		"goal":    goal,
	}).Info("Updating user goal")

	user, err := h.userService.GetByTelegramID(userID)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"error":   err,
		}).Error("Failed to get user for goal update")
		handlers.SendErrorMessage(h.bot, chatID, "Пользователь не найден")
		return nil
	}

	switch goal {
	case "muscle_gain":
		user.Goals = models.GoalsSlice{models.GoalMuscleGain}
	case "strength":
		user.Goals = models.GoalsSlice{models.GoalStrength}
	case "endurance":
		user.Goals = models.GoalsSlice{models.GoalStrength}
	case "weight_loss":
		user.Goals = models.GoalsSlice{models.GoalWeightLoss}
	default:
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"goal":    goal,
		}).Error("Unknown goal type")
		handlers.SendErrorMessage(h.bot, chatID, "Неизвестная цель")
		return nil
	}

	if err := h.userService.CreateOrUpdate(user); err != nil {
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"goal":    goal,
			"error":   err,
		}).Error("Failed to update user goal")
		handlers.SendErrorMessage(h.bot, chatID, "Ошибка обновления целей")
		return nil
	}

	logger.WithFields(logrus.Fields{
		"user_id": userID,
		"chat_id": chatID,
		"goal":    goal,
	}).Info("User goal updated successfully")

	msg := tgbotapi.NewMessage(chatID, "✅ Цель обновлена!")
	_, err = h.bot.Send(msg)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"goal":    goal,
			"error":   err,
		}).Error("Failed to send goal update confirmation")
	}
	return err
}
