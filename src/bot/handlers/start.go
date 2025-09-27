package handlers

import (
	"workouts_bot/pkg/logger"
	"workouts_bot/src/bot/keyboards"
	"workouts_bot/src/database/models"
	"workouts_bot/src/services"

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
	bot         *tgbotapi.BotAPI
	userService *services.UserService
}

func NewStartHandler(bot *tgbotapi.BotAPI, database *gorm.DB) *StartHandler {
	return &StartHandler{
		bot:         bot,
		userService: services.NewUserService(database),
	}
}

func (startHandler *StartHandler) Handle(update tgbotapi.Update) error {
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
		Goals:      models.GoalsSlice{models.GoalMuscleGain},
		Experience: 1,
	}

	if err := startHandler.userService.CreateOrUpdate(user); err != nil {
		logger.Error("Failed to create or update user:", err)
		return err
	}

	keyboard := keyboards.CreateMainMenu()

	chatId := GetChatId(update)

	message := tgbotapi.NewMessage(chatId, helloMessage)
	message.ReplyMarkup = &keyboard

	_, err := startHandler.bot.Send(message)
	return err
}
