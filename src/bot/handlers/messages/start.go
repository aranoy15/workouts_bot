package messages

import (
	"workouts_bot/src/bot/handlers"
	"workouts_bot/src/bot/keyboards"
	"workouts_bot/src/database"
	"workouts_bot/src/logger"
	"workouts_bot/src/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	helloMessage = "Привет! Я бот для тренировок 🏋️\n\n" +
		"Я помогу тебе:\n" +
		"• Подобрать упражнения\n" +
		"• Составить программу тренировок\n" +
		"• Отслеживать прогресс\n" +
		"• Записывать подходы и веса\n\n" +
		"Выбери действие:"
)

const (
	StartCommand = "/start"
)

type StartHandler struct {
	bot      *tgbotapi.BotAPI
	database *gorm.DB
}

func NewStartHandler(
	bot *tgbotapi.BotAPI,
	database *gorm.DB,
) *StartHandler {
	return &StartHandler{
		bot:      bot,
		database: database,
	}
}

func (startHandler *StartHandler) Handle(
	update tgbotapi.Update,
) error {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID
	userName := update.Message.From.UserName
	firstName := update.Message.From.FirstName

	logger.WithFields(logrus.Fields{
		"user_id":    userID,
		"user_name":  userName,
		"first_name": firstName,
	}).Info("New user started bot")

	user := &models.User{
		TelegramID: userID,
		Username:   userName,
		FirstName:  firstName,
	}

	if err := database.UpsertUser(user, startHandler.database); err != nil {
		logger.Error("Failed to create or update user:", err)
		handlers.SendErrorMessage(
			startHandler.bot, chatID,
			"Ошибка при получении пользователя",
		)
		return err
	}

	msg := tgbotapi.NewMessage(chatID, helloMessage)
	msg.ReplyMarkup = keyboards.CreateMainMenu()

	_, err := startHandler.bot.Send(msg)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"chat_id": chatID,
			"user_id": userID,
			"error":   err,
		}).Error("Failed to send start message")
		return err
	}
	return nil
}

func (startHandler *StartHandler) MainMenu(
	userID int64,
	chatID int64,
	messageID int,
) error {
	msg := tgbotapi.NewMessage(chatID, "Главное меню")
	msg.ReplyMarkup = keyboards.CreateMainMenu()
	_, err := startHandler.bot.Send(msg)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"chat_id": chatID,
			"user_id": userID,
			"error":   err,
		}).Error("Failed to send main menu")
		return err
	}
	return nil
}
