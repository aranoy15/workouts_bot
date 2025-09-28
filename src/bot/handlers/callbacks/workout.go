package callbacks

import (
	"fmt"
	"strconv"
	"strings"
	"workouts_bot/pkg/logger"
	"workouts_bot/src/bot/handlers"
	"workouts_bot/src/bot/keyboards"
	"workouts_bot/src/database/models"
	"workouts_bot/src/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const WorkoutCallbackType = "workout"

type WorkoutHandler struct {
	bot             *tgbotapi.BotAPI
	workoutsService *services.WorkoutsService
	userService     *services.UserService
}

func NewWorkoutHandler(bot *tgbotapi.BotAPI, database *gorm.DB) *WorkoutHandler {
	return &WorkoutHandler{
		bot:             bot,
		workoutsService: services.NewWorkoutsService(database),
		userService:     services.NewUserService(database),
	}
}

func (h *WorkoutHandler) Handle(update tgbotapi.Update) error {
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
	}).Info("Workout callback received")

	if len(parts) < 3 {
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"data":    data,
		}).Error("Invalid workout callback format")
		handlers.SendErrorMessage(h.bot, chatID, "Неверный формат команды")
		return nil
	}

	action := parts[1]
	workoutID, err := strconv.ParseUint(parts[2], 10, 32)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id":    userID,
			"chat_id":    chatID,
			"workout_id": parts[2],
			"error":      err,
		}).Error("Failed to parse workout ID")
		handlers.SendErrorMessage(h.bot, chatID, "Неверный ID тренировки")
		return nil
	}

	switch action {
	case "start":
		return h.startWorkout(userID, chatID, messageID, workoutID)
	case "edit":
		return h.editWorkout(userID, chatID, messageID, workoutID)
	case "delete":
		return h.deleteWorkout(userID, chatID, messageID, workoutID)
	default:
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"action":  action,
		}).Error("Unknown workout action")
		handlers.SendErrorMessage(h.bot, chatID, "Неизвестное действие")
		return nil
	}
}

func (h *WorkoutHandler) startWorkout(
	userID int64,
	chatID int64,
	messageID int,
	workoutID uint64,
) error {
	logger.WithFields(logrus.Fields{
		"user_id":    userID,
		"chat_id":    chatID,
		"workout_id": workoutID,
	}).Info("Starting workout")

	workout, err := h.workoutsService.GetWorkoutByID(workoutID)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id":    userID,
			"chat_id":    chatID,
			"workout_id": workoutID,
			"error":      err,
		}).Error("Failed to get workout")
		handlers.SendErrorMessage(h.bot, chatID, "Тренировка не найдена")
		return nil
	}

	user, err := h.userService.GetByTelegramID(userID)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"error":   err,
		}).Error("Failed to get user by telegram ID")
		handlers.SendErrorMessage(h.bot, chatID, "Ошибка при получении пользователя")
		return nil
	}

	if workout.UserID != user.ID {
		logger.WithFields(logrus.Fields{
			"user_id":    userID,
			"chat_id":    chatID,
			"workout_id": workoutID,
			"owner_id":   workout.UserID,
		}).Error("User trying to access workout they don't own")
		handlers.SendErrorMessage(h.bot, chatID, "Нет доступа к этой тренировке")
		return nil
	}

	if err := h.workoutsService.StartWorkout(workoutID); err != nil {
		logger.WithFields(logrus.Fields{
			"user_id":    userID,
			"chat_id":    chatID,
			"workout_id": workoutID,
			"error":      err,
		}).Error("Failed to start workout")
		handlers.SendErrorMessage(h.bot, chatID, "Ошибка начала тренировки")
		return nil
	}

	return h.showFirstExercise(userID, chatID, messageID, workout)
}

func (h *WorkoutHandler) editWorkout(
	userID int64,
	chatID int64,
	messageID int,
	workoutID uint64,
) error {
	logger.WithFields(logrus.Fields{
		"user_id":    userID,
		"chat_id":    chatID,
		"workout_id": workoutID,
	}).Info("Editing workout")

	workout, err := h.workoutsService.GetWorkoutByID(workoutID)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id":    userID,
			"chat_id":    chatID,
			"workout_id": workoutID,
			"error":      err,
		}).Error("Failed to get workout for editing")
		handlers.SendErrorMessage(h.bot, chatID, "Тренировка не найдена")
		return nil
	}

	user, err := h.userService.GetByTelegramID(userID)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"error":   err,
		}).Error("Failed to get user by telegram ID")
		handlers.SendErrorMessage(h.bot, chatID, "Ошибка при получении пользователя")
		return nil
	}

	if workout.UserID != user.ID {
		logger.WithFields(logrus.Fields{
			"user_id":    userID,
			"chat_id":    chatID,
			"workout_id": workoutID,
			"owner_id":   workout.UserID,
		}).Error("User trying to edit workout they don't own")
		handlers.SendErrorMessage(h.bot, chatID, "Нет доступа к этой тренировке")
		return nil
	}

	text := fmt.Sprintf(
		"✏️ Редактирование тренировки: %s\n\nВыберите что изменить:",
		workout.Name,
	)

	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)

	_, err = h.bot.Send(editMsg)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id":    userID,
			"chat_id":    chatID,
			"workout_id": workoutID,
			"error":      err,
		}).Error("Failed to send edit workout message")
	}
	return err
}

func (h *WorkoutHandler) deleteWorkout(
	userID int64,
	chatID int64,
	messageID int,
	workoutID uint64,
) error {
	logger.WithFields(logrus.Fields{
		"user_id":    userID,
		"chat_id":    chatID,
		"workout_id": workoutID,
	}).Info("Deleting workout")

	workout, err := h.workoutsService.GetWorkoutByID(workoutID)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id":    userID,
			"chat_id":    chatID,
			"workout_id": workoutID,
			"error":      err,
		}).Error("Failed to get workout for deletion")
		handlers.SendErrorMessage(h.bot, chatID, "Тренировка не найдена")
		return nil
	}

	user, err := h.userService.GetByTelegramID(userID)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"error":   err,
		}).Error("Failed to get user by telegram ID")
		handlers.SendErrorMessage(h.bot, chatID, "Ошибка при получении пользователя")
		return nil
	}

	if workout.UserID != user.ID {
		logger.WithFields(logrus.Fields{
			"user_id":    userID,
			"chat_id":    chatID,
			"workout_id": workoutID,
			"owner_id":   workout.UserID,
		}).Error("User trying to delete workout they don't own")
		handlers.SendErrorMessage(h.bot, chatID, "Нет доступа к этой тренировке")
		return nil
	}

	text := fmt.Sprintf(
		"🗑️ Удалить тренировку %q?\n\nЭто действие нельзя отменить!",
		workout.Name,
	)

	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	keyboard := keyboards.CreateConfirmationKeyboard(
		fmt.Sprintf("delete_workout:%d", workoutID),
	)
	editMsg.ReplyMarkup = &keyboard

	_, err = h.bot.Send(editMsg)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id":    userID,
			"chat_id":    chatID,
			"workout_id": workoutID,
			"error":      err,
		}).Error("Failed to send delete confirmation message")
	}
	return err
}

func (h *WorkoutHandler) showFirstExercise(
	userID int64,
	chatID int64,
	messageID int,
	workout *models.Workout,
) error {
	logger.WithFields(logrus.Fields{
		"user_id":    userID,
		"chat_id":    chatID,
		"workout_id": workout.ID,
	}).Info("Showing first exercise")

	if len(workout.Exercises) == 0 {
		logger.WithFields(logrus.Fields{
			"user_id":    userID,
			"chat_id":    chatID,
			"workout_id": workout.ID,
		}).Error("Workout has no exercises")
		handlers.SendErrorMessage(h.bot, chatID, "В тренировке нет упражнений")
		return nil
	}

	firstExercise := workout.Exercises[0]

	text := fmt.Sprintf("🏋️ Тренировка: %s\n\n"+
		"📖 Упражнение: %s\n"+
		"📊 Подходы: %d\n"+
		"🔄 Повторения: %d\n"+
		"⏱️ Отдых: %d сек\n"+
		"🏋️ Вес: %.1f кг\n\n"+
		"Готовы начать?",
		workout.Name,
		firstExercise.Exercise.Name,
		firstExercise.SetsCount,
		firstExercise.RepsCount,
		firstExercise.RestSeconds,
		firstExercise.WeightKg,
	)

	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	keyboard := keyboards.CreateSetKeyboard(firstExercise.ID, 1)
	editMsg.ReplyMarkup = &keyboard

	_, err := h.bot.Send(editMsg)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id":     userID,
			"chat_id":     chatID,
			"workout_id":  workout.ID,
			"exercise_id": firstExercise.ID,
			"error":       err,
		}).Error("Failed to send first exercise message")
	}
	return err
}
