package callbacks

import (
	"strconv"
	"strings"
	"workouts_bot/src/logger"
	"workouts_bot/src/bot/handlers"
	"workouts_bot/src/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const ExperienceCallbackType = "experience"

type ExperienceHandler struct {
	bot         *tgbotapi.BotAPI
	userService *services.UserService
}

func NewExperienceHandler(bot *tgbotapi.BotAPI, database *gorm.DB) *ExperienceHandler {
	return &ExperienceHandler{
		bot:         bot,
		userService: services.NewUserService(database),
	}
}

func (h *ExperienceHandler) Handle(update tgbotapi.Update) error {
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
	}).Info("Experience callback received")

	if len(parts) < 2 {
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"data":    data,
		}).Error("Invalid experience callback format")
		handlers.SendErrorMessage(h.bot, chatID, "Неверный формат команды")
		return nil
	}

	experienceStr := parts[1]

	experience, err := strconv.Atoi(experienceStr)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id":        userID,
			"chat_id":        chatID,
			"experience_str": experienceStr,
			"error":          err,
		}).Error("Failed to parse experience level")
		handlers.SendErrorMessage(h.bot, chatID, "Неверный уровень опыта")
		return nil
	}

	logger.WithFields(logrus.Fields{
		"user_id":    userID,
		"chat_id":    chatID,
		"experience": experience,
	}).Info("Updating user experience level")

	user, err := h.userService.GetByTelegramID(userID)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"error":   err,
		}).Error("Failed to get user for experience update")
		handlers.SendErrorMessage(h.bot, chatID, "Пользователь не найден")
		return nil
	}

	user.Experience = experience

	if err := h.userService.CreateOrUpdate(user); err != nil {
		logger.WithFields(logrus.Fields{
			"user_id":    userID,
			"chat_id":    chatID,
			"experience": experience,
			"error":      err,
		}).Error("Failed to update user experience")
		handlers.SendErrorMessage(h.bot, chatID, "Ошибка обновления уровня опыта")
		return nil
	}

	logger.WithFields(logrus.Fields{
		"user_id":    userID,
		"chat_id":    chatID,
		"experience": experience,
	}).Info("User experience updated successfully")

	msg := tgbotapi.NewMessage(chatID, "✅ Уровень опыта обновлен!")
	_, err = h.bot.Send(msg)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id":    userID,
			"chat_id":    chatID,
			"experience": experience,
			"error":      err,
		}).Error("Failed to send experience update confirmation")
	}
	return err
}
