package callbacks

import (
	"strconv"
	"strings"
	"workouts_bot/pkg/logger"
	"workouts_bot/src/bot/handlers"
	"workouts_bot/src/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const DurationCallbackType = "duration"

type DurationHandler struct {
	bot             *tgbotapi.BotAPI
	workoutsService *services.WorkoutsService
	userService     *services.UserService
}

func NewDurationHandler(
	bot *tgbotapi.BotAPI,
	database *gorm.DB,
) *DurationHandler {
	return &DurationHandler{
		bot:             bot,
		workoutsService: services.NewWorkoutsService(database),
		userService:     services.NewUserService(database),
	}
}

func (h *DurationHandler) Handle(update tgbotapi.Update) error {
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
	}).Info("Duration callback received")

	if len(parts) < 2 {
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"data":    data,
		}).Error("Invalid duration callback format")
		handlers.SendErrorMessage(h.bot, chatID, "Неверный формат команды")
		return nil
	}

	durationStr := parts[1]

	duration, err := strconv.Atoi(durationStr)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id":      userID,
			"chat_id":      chatID,
			"duration_str": durationStr,
			"error":        err,
		}).Error("Failed to parse duration")
		handlers.SendErrorMessage(h.bot, chatID, "Неверная продолжительность")
		return nil
	}

	logger.WithFields(logrus.Fields{
		"user_id":  userID,
		"chat_id":  chatID,
		"duration": duration,
	}).Info("Processing duration selection")

	logger.WithFields(logrus.Fields{
		"user_id":  userID,
		"chat_id":  chatID,
		"duration": duration,
	}).Info("Creating workout with selected duration")

	// TODO: Implement actual workout creation logic
	// For now, just send a confirmation message
	text := "✅ Тренировка создана!\n\n" +
		"Продолжительность: " + durationStr + " минут\n\n" +
		"Тренировка добавлена в ваш список."

	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)

	_, err = h.bot.Send(editMsg)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id":  userID,
			"chat_id":  chatID,
			"duration": duration,
			"error":    err,
		}).Error("Failed to send workout creation confirmation")
	}
	return err
}
