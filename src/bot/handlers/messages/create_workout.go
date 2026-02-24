package messages

import (
	"workouts_bot/src/logger"
	"workouts_bot/src/bot/handlers"
	"workouts_bot/src/bot/keyboards"
	"workouts_bot/src/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CreateWorkoutHandler struct {
	bot             *tgbotapi.BotAPI
	workoutsService *services.WorkoutsService
	userService     *services.UserService
}

func NewCreateWorkoutHandler(bot *tgbotapi.BotAPI, database *gorm.DB) *CreateWorkoutHandler {
	return &CreateWorkoutHandler{
		bot:             bot,
		workoutsService: services.NewWorkoutsService(database),
		userService:     services.NewUserService(database),
	}
}

func (handler *CreateWorkoutHandler) Handle(update tgbotapi.Update) error {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	logger.WithFields(logrus.Fields{
		"chat_id": chatID,
		"user_id": userID,
	}).Info("Create workout handler called")

	user, err := handler.userService.GetByTelegramID(userID)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"error":   err,
		}).Error("Failed to get user by telegram ID")
		handlers.SendErrorMessage(
			handler.bot, chatID,
			"Ошибка при получении пользователя",
		)
		return err
	}

	logger.WithFields(logrus.Fields{
		"user_id":    userID,
		"chat_id":    chatID,
		"user_db_id": user.ID,
	}).Info("User found, showing workout type selection")

	message := "🏋️ Создание тренировки\n\n" +
		"Выберите тип тренировки:"

	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyMarkup = keyboards.CreateWorkoutTypeKeyboard()

	_, err = handler.bot.Send(msg)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"error":   err,
		}).Error("Failed to send workout type selection message")
	}
	return err
}
