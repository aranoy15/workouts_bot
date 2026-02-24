package callbacks

import (
	"strconv"
	"strings"
	"time"
	"workouts_bot/src/logger"
	"workouts_bot/src/bot/handlers"
	"workouts_bot/src/models"
	"workouts_bot/src/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const SetCallbackType = "set"

type SetHandler struct {
	bot             *tgbotapi.BotAPI
	workoutsService *services.WorkoutsService
	userService     *services.UserService
	database        *gorm.DB
}

func NewSetHandler(bot *tgbotapi.BotAPI, database *gorm.DB) *SetHandler {
	return &SetHandler{
		bot:             bot,
		workoutsService: services.NewWorkoutsService(database),
		userService:     services.NewUserService(database),
		database:        database,
	}
}

func (h *SetHandler) Handle(update tgbotapi.Update) error {
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
	}).Info("Set callback received")

	if len(parts) < 4 {
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"data":    data,
		}).Error("Invalid set callback format")
		handlers.SendErrorMessage(h.bot, chatID, "Неверный формат команды")
		return nil
	}

	action := parts[1]
	workoutExerciseID, err := strconv.ParseUint(parts[2], 10, 32)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id":             userID,
			"chat_id":             chatID,
			"workout_exercise_id": parts[2],
			"error":               err,
		}).Error("Failed to parse workout exercise ID")
		handlers.SendErrorMessage(h.bot, chatID, "Неверный ID упражнения в тренировке")
		return nil
	}

	setNumber, err := strconv.Atoi(parts[3])
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id":    userID,
			"chat_id":    chatID,
			"set_number": parts[3],
			"error":      err,
		}).Error("Failed to parse set number")
		handlers.SendErrorMessage(h.bot, chatID, "Неверный номер подхода")
		return nil
	}

	switch action {
	case "complete":
		return h.completeSet(update, uint(workoutExerciseID), setNumber)
	case "skip":
		return h.skipSet(update, uint(workoutExerciseID), setNumber)
	case "pause":
		return h.pauseWorkout(update, uint(workoutExerciseID))
	default:
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"action":  action,
		}).Error("Unknown set action")
		handlers.SendErrorMessage(h.bot, chatID, "Неизвестное действие")
		return nil
	}
}

func (h *SetHandler) completeSet(
	update tgbotapi.Update,
	workoutExerciseID uint,
	setNumber int,
) error {
	set := &models.Set{
		WorkoutExerciseID: workoutExerciseID,
		SetNumber:         setNumber,
		WeightKg:          0,
		RepsDone:          0,
		RestTakenSeconds:  0,
		CompletedAt:       time.Now(),
	}

	err := h.database.Create(set).Error
	if err != nil {
		handlers.SendErrorMessage(
			h.bot,
			update.CallbackQuery.Message.Chat.ID,
			"Ошибка записи подхода",
		)
		return nil
	}

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "✅ Подход завершен!")
	_, err = h.bot.Send(msg)
	return err
}

func (h *SetHandler) skipSet(update tgbotapi.Update, workoutExerciseID uint, setNumber int) error {
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "⏭️ Подход пропущен")
	_, err := h.bot.Send(msg)
	return err
}

func (h *SetHandler) pauseWorkout(update tgbotapi.Update, workoutExerciseID uint) error {
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "⏸️ Тренировка приостановлена")
	_, err := h.bot.Send(msg)
	return err
}
