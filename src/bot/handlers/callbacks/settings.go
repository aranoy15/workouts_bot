package callbacks

import (
	"strings"
	"workouts_bot/pkg/logger"
	"workouts_bot/src/bot/handlers"
	"workouts_bot/src/bot/keyboards"
	"workouts_bot/src/constants"
	"workouts_bot/src/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const SettingsCallbackType = "settings"

type SettingsHandler struct {
	bot         *tgbotapi.BotAPI
	userService *services.UserService
}

func NewSettingsHandler(bot *tgbotapi.BotAPI, database *gorm.DB) *SettingsHandler {
	return &SettingsHandler{
		bot:         bot,
		userService: services.NewUserService(database),
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
	case "goals":
		return h.showGoalsMenu(userID, chatID, messageID)
	case "equipment":
		return h.showEquipmentMenu(userID, chatID, messageID)
	case "experience":
		return h.showExperienceMenu(userID, chatID, messageID)
	case "limitations":
		return h.showLimitationsMenu(userID, chatID, messageID)
	case constants.NavigationMain:
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

func (h *SettingsHandler) showGoalsMenu(
	userID int64,
	chatID int64,
	messageID int,
) error {
	logger.WithFields(logrus.Fields{
		"user_id": userID,
		"chat_id": chatID,
	}).Info("Showing goals menu")

	text := "🎯 Выберите ваши цели тренировок:"

	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	keyboard := keyboards.CreateGoalSelectionKeyboard()
	editMsg.ReplyMarkup = &keyboard

	_, err := h.bot.Send(editMsg)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"error":   err,
		}).Error("Failed to send goals menu")
	}
	return err
}

func (h *SettingsHandler) showEquipmentMenu(
	userID int64,
	chatID int64,
	messageID int,
) error {
	logger.WithFields(logrus.Fields{
		"user_id": userID,
		"chat_id": chatID,
	}).Info("Showing equipment menu")

	text := "🏋️ Какое оборудование у вас есть?"

	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	keyboard := keyboards.CreateEquipmentKeyboard()
	editMsg.ReplyMarkup = &keyboard

	_, err := h.bot.Send(editMsg)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"error":   err,
		}).Error("Failed to send equipment menu")
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

func (h *SettingsHandler) showLimitationsMenu(
	userID int64,
	chatID int64,
	messageID int,
) error {
	logger.WithFields(logrus.Fields{
		"user_id": userID,
		"chat_id": chatID,
	}).Info("Showing limitations menu")

	text := "⚠️ Есть ли у вас ограничения по здоровью?"

	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	keyboard := keyboards.CreateBackKeyboard("settings:main")
	editMsg.ReplyMarkup = &keyboard

	_, err := h.bot.Send(editMsg)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"error":   err,
		}).Error("Failed to send limitations menu")
	}
	return err
}
