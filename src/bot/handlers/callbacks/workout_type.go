package callbacks

import (
	"fmt"
	"strings"
	"workouts_bot/pkg/logger"
	"workouts_bot/src/bot/handlers"
	"workouts_bot/src/bot/keyboards"
	"workouts_bot/src/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const WorkoutTypeCallbackType = "workout_type"

type WorkoutTypeHandler struct {
	bot             *tgbotapi.BotAPI
	workoutsService *services.WorkoutsService
	userService     *services.UserService
}

func NewWorkoutTypeHandler(
	bot *tgbotapi.BotAPI,
	database *gorm.DB,
) *WorkoutTypeHandler {
	return &WorkoutTypeHandler{
		bot:             bot,
		workoutsService: services.NewWorkoutsService(database),
		userService:     services.NewUserService(database),
	}
}

func (h *WorkoutTypeHandler) Handle(update tgbotapi.Update) error {
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
	}).Info("Workout type callback received")

	if len(parts) < 2 {
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"data":    data,
		}).Error("Invalid workout type callback format")
		handlers.SendErrorMessage(h.bot, chatID, "Неверный формат команды")
		return nil
	}

	workoutType := parts[1]

	logger.WithFields(logrus.Fields{
		"user_id":      userID,
		"chat_id":      chatID,
		"workout_type": workoutType,
	}).Info("Processing workout type selection")

	switch workoutType {
	case "main":
		return h.showMainWorkoutTypeMenu(userID, chatID, messageID)
	case "split":
		return h.createSplitWorkout(userID, chatID, messageID)
	case "push_pull":
		return h.createPushPullWorkout(userID, chatID, messageID)
	case "fullbody":
		return h.createFullBodyWorkout(userID, chatID, messageID)
	case "custom":
		return h.createCustomWorkout(userID, chatID, messageID)
	default:
		logger.WithFields(logrus.Fields{
			"user_id":      userID,
			"chat_id":      chatID,
			"workout_type": workoutType,
		}).Error("Unknown workout type")
		handlers.SendErrorMessage(h.bot, chatID, "Неизвестный тип тренировки")
		return nil
	}
}

func (h *WorkoutTypeHandler) createWorkoutWithDuration(
	userID int64,
	chatID int64,
	messageID int,
	workoutType string,
	title string,
	emoji string,
) error {
	logger.WithFields(logrus.Fields{
		"user_id":      userID,
		"chat_id":      chatID,
		"workout_type": workoutType,
	}).Info("Creating workout")

	text := fmt.Sprintf("%s %s\n\nВыберите продолжительность тренировки:", emoji, title)

	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	keyboard := keyboards.CreateWorkoutDurationKeyboard(workoutType)
	editMsg.ReplyMarkup = &keyboard

	_, err := h.bot.Send(editMsg)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id":      userID,
			"chat_id":      chatID,
			"workout_type": workoutType,
			"error":        err,
		}).Error("Failed to send workout duration selection")
	}
	return err
}

func (h *WorkoutTypeHandler) createSplitWorkout(
	userID int64,
	chatID int64,
	messageID int,
) error {
	return h.createWorkoutWithDuration(userID, chatID, messageID, "split", "Классический сплит", "🏋️")
}

func (h *WorkoutTypeHandler) createPushPullWorkout(
	userID int64,
	chatID int64,
	messageID int,
) error {
	return h.createWorkoutWithDuration(userID, chatID, messageID, "push_pull", "Push/Pull/Legs", "🔄")
}

func (h *WorkoutTypeHandler) createFullBodyWorkout(
	userID int64,
	chatID int64,
	messageID int,
) error {
	logger.WithFields(logrus.Fields{
		"user_id": userID,
		"chat_id": chatID,
	}).Info("Creating full body workout")

	text := "💪 Фулбади\n\n" +
		"Выберите продолжительность тренировки:"

	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	keyboard := keyboards.CreateWorkoutDurationKeyboard("fullbody")
	editMsg.ReplyMarkup = &keyboard

	_, err := h.bot.Send(editMsg)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"error":   err,
		}).Error("Failed to send full body workout duration selection")
	}
	return err
}

func (h *WorkoutTypeHandler) createCustomWorkout(
	userID int64,
	chatID int64,
	messageID int,
) error {
	logger.WithFields(logrus.Fields{
		"user_id": userID,
		"chat_id": chatID,
	}).Info("Creating custom workout")

	text := "🎯 Кастомная тренировка\n\n" +
		"Выберите продолжительность тренировки:"

	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	keyboard := keyboards.CreateWorkoutDurationKeyboard("custom")
	editMsg.ReplyMarkup = &keyboard

	_, err := h.bot.Send(editMsg)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"error":   err,
		}).Error("Failed to send custom workout duration selection")
	}
	return err
}

func (h *WorkoutTypeHandler) showMainWorkoutTypeMenu(
	userID int64,
	chatID int64,
	messageID int,
) error {
	logger.WithFields(logrus.Fields{
		"user_id": userID,
		"chat_id": chatID,
	}).Info("Showing main workout type menu")

	text := "🏋️ Создание тренировки\n\n" +
		"Выберите тип тренировки:"

	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	keyboard := keyboards.CreateWorkoutTypeKeyboard()
	editMsg.ReplyMarkup = &keyboard

	_, err := h.bot.Send(editMsg)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"error":   err,
		}).Error("Failed to send main workout type menu")
	}
	return err
}
