package callbacks

import (
	"strings"
	"workouts_bot/src/bot/handlers"
	"workouts_bot/src/bot/keyboards"
	"workouts_bot/src/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const SettingsCallbackType = "settings"

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

func (h *SettingsHandler) Handle(update tgbotapi.Update) error {
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
	}).Info("Settings callback received")

	if len(parts) < 2 {
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"data":    data,
		}).Error("Invalid settings callback format")
		handlers.SendErrorMessage(h.bot, chatID, "Неверный формат команды")
		return nil
	}

	setting := parts[1]

	switch setting {
	case "experience":
		return h.showExperienceMenu(userID, chatID, messageID)
	case "experience_back":
		return h.showMainSettingsMenu(userID, chatID, messageID)
	default:
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"setting": setting,
		}).Error("Unknown settings option")
		handlers.SendErrorMessage(h.bot, chatID, "Неизвестная настройка")
		return nil
	}
}

func (h *SettingsHandler) showMainSettingsMenu(
	userID int64,
	chatID int64,
	messageID int,
) error {
	logger.WithFields(logrus.Fields{
		"user_id": userID,
		"chat_id": chatID,
	}).Info("Showing main settings menu")

	text := "⚙️ Настройки"

	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	keyboard := keyboards.CreateSettingsKeyboard()
	editMsg.ReplyMarkup = &keyboard

	_, err := h.bot.Send(editMsg)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"error":   err,
		}).Error("Failed to send main settings menu")
	}
	return err
}

func (h *SettingsHandler) showExperienceMenu(
	userID int64,
	chatID int64,
	messageID int,
) error {
	logger.WithFields(logrus.Fields{
		"user_id": userID,
		"chat_id": chatID,
	}).Info("Showing experience menu")

	text := "📈 Какой у вас уровень опыта в тренировках?"

	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	keyboard := keyboards.CreateExperienceLevelKeyboard()
	editMsg.ReplyMarkup = &keyboard

	_, err := h.bot.Send(editMsg)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"error":   err,
		}).Error("Failed to send experience menu")
	}
	return err
}
