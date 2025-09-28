package messages

import (
	"time"
	"workouts_bot/pkg/logger"
	"workouts_bot/src/bot/handlers"
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

func NewStartHandler(
	bot *tgbotapi.BotAPI,
	database *gorm.DB,
) *StartHandler {
	return &StartHandler{
		bot:         bot,
		userService: services.NewUserService(database),
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

	user, err := startHandler.userService.GetByTelegramID(userID)
	if err != nil {
		logger.Info("User not found, creating new user:", chatID)
		user = &models.User{
			TelegramID: userID,
			Username:   userName,
			FirstName:  firstName,
			Goals:      models.GoalsSlice{models.GoalMuscleGain},
			Experience: 1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
	}

	if err := startHandler.userService.CreateOrUpdate(user); err != nil {
		logger.Error("Failed to create or update user:", err)
		handlers.SendErrorMessage(
			startHandler.bot, chatID,
			"Ошибка при получении пользователя",
		)
		return err
	}

	msg := tgbotapi.NewMessage(chatID, helloMessage)
	msg.ReplyMarkup = keyboards.CreateMainMenu()

	_, err = startHandler.bot.Send(msg)
	return err
}
