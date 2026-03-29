package messages

import (
	"fmt"
	"workouts_bot/src/bot/handlers"
	"workouts_bot/src/bot/keyboards"
	"workouts_bot/src/database"
	"workouts_bot/src/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SettingsHandler struct {
	bot      *tgbotapi.BotAPI
	database *gorm.DB
}

func NewSettingsHandler(bot *tgbotapi.BotAPI, database *gorm.DB) *SettingsHandler {
	return &SettingsHandler{
		bot:      bot,
		database: database,
	}
}

func (handler *SettingsHandler) Handle(update tgbotapi.Update) error {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	logger.WithFields(logrus.Fields{
		"chat_id": chatID,
		"user_id": userID,
	}).Info("Settings handler")

	user, err := database.GetUserByTelegramID(userID, handler.database)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"chat_id": chatID,
			"user_id": userID,
			"error":   err,
		}).Error("Failed to get user by telegram ID")
		handlers.SendErrorMessage(handler.bot, chatID, "Пользователь не найден")
		return nil
	}

	message := "⚙️ Ваши настройки:\n\n"
	message += fmt.Sprintf("📈 Уровень опыта: %s\n\n", getExperienceLevel(user.Experience))
	message += "Используйте кнопки ниже для изменения настроек:"

	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyMarkup = keyboards.CreateSettingsKeyboard()

	_, err = handler.bot.Send(msg)
	return err
}

func getExperienceLevel(experience int) string {
	switch {
	case experience < 1:
		return "🟢 Начинающий"
	case experience < 3:
		return "🟡 Средний"
	case experience < 5:
		return "🟠 Продвинутый"
	default:
		return "🔴 Эксперт"
	}
}
