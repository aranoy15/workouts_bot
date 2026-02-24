package callbacks

import (
	"strings"
	"workouts_bot/src/logger"
	"workouts_bot/src/bot/handlers"
	"workouts_bot/src/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const EquipmentCallbackType = "equipment"

type EquipmentHandler struct {
	bot         *tgbotapi.BotAPI
	userService *services.UserService
}

func NewEquipmentHandler(bot *tgbotapi.BotAPI, database *gorm.DB) *EquipmentHandler {
	return &EquipmentHandler{
		bot:         bot,
		userService: services.NewUserService(database),
	}
}

func (h *EquipmentHandler) Handle(update tgbotapi.Update) error {
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
	}).Info("Equipment callback received")

	if len(parts) < 2 {
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"data":    data,
		}).Error("Invalid equipment callback format")
		handlers.SendErrorMessage(h.bot, chatID, "Неверный формат команды")
		return nil
	}

	equipment := parts[1]

	logger.WithFields(logrus.Fields{
		"user_id":   userID,
		"chat_id":   chatID,
		"equipment": equipment,
	}).Info("Updating user equipment")

	user, err := h.userService.GetByTelegramID(userID)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"error":   err,
		}).Error("Failed to get user for equipment update")
		handlers.SendErrorMessage(h.bot, chatID, "Пользователь не найден")
		return nil
	}

	switch equipment {
	case "home":
		err := h.userService.UpdateEquipment(user.ID, []int{1})
		if err != nil {
			logger.WithFields(logrus.Fields{
				"user_id":   userID,
				"chat_id":   chatID,
				"equipment": equipment,
				"error":     err,
			}).Error("Failed to update equipment to home")
			handlers.SendErrorMessage(h.bot, chatID, "Ошибка обновления оборудования")
			return nil
		}
	case "gym":
		err := h.userService.UpdateEquipment(user.ID, []int{2})
		if err != nil {
			logger.WithFields(logrus.Fields{
				"user_id":   userID,
				"chat_id":   chatID,
				"equipment": equipment,
				"error":     err,
			}).Error("Failed to update equipment to gym")
			handlers.SendErrorMessage(h.bot, chatID, "Ошибка обновления оборудования")
			return nil
		}
	case "none":
		err := h.userService.UpdateEquipment(user.ID, []int{0})
		if err != nil {
			logger.WithFields(logrus.Fields{
				"user_id":   userID,
				"chat_id":   chatID,
				"equipment": equipment,
				"error":     err,
			}).Error("Failed to update equipment to none")
			handlers.SendErrorMessage(h.bot, chatID, "Ошибка обновления оборудования")
			return nil
		}
	case "custom":
		return h.showCustomEquipmentMenu(userID, chatID, messageID)
	default:
		logger.WithFields(logrus.Fields{
			"user_id":   userID,
			"chat_id":   chatID,
			"equipment": equipment,
		}).Error("Unknown equipment type")
		handlers.SendErrorMessage(h.bot, chatID, "Неизвестное оборудование")
		return nil
	}

	logger.WithFields(logrus.Fields{
		"user_id":   userID,
		"chat_id":   chatID,
		"equipment": equipment,
	}).Info("User equipment updated successfully")

	msg := tgbotapi.NewMessage(chatID, "✅ Оборудование обновлено!")
	_, err = h.bot.Send(msg)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id":   userID,
			"chat_id":   chatID,
			"equipment": equipment,
			"error":     err,
		}).Error("Failed to send equipment update confirmation")
	}
	return err
}

func (h *EquipmentHandler) showCustomEquipmentMenu(
	userID int64,
	chatID int64,
	messageID int,
) error {
	logger.WithFields(logrus.Fields{
		"user_id": userID,
		"chat_id": chatID,
	}).Info("Showing custom equipment menu")

	text := "⚙️ Настройка оборудования"

	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)

	_, err := h.bot.Send(editMsg)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"error":   err,
		}).Error("Failed to send custom equipment menu")
	}
	return err
}
