package callbacks

import (
	"fmt"
	"strconv"
	"strings"
	"workouts_bot/pkg/logger"
	"workouts_bot/src/bot/handlers"
	"workouts_bot/src/constants"
	"workouts_bot/src/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const DurationCallbackType = "duration"

type DurationHandler struct {
	bot             *tgbotapi.BotAPI
	workoutsService *services.WorkoutsService
	userService     *services.UserService
}

func NewDurationHandler(
	bot *tgbotapi.BotAPI,
	database *gorm.DB,
) *DurationHandler {
	return &DurationHandler{
		bot:             bot,
		workoutsService: services.NewWorkoutsService(database),
		userService:     services.NewUserService(database),
	}
}

func (h *DurationHandler) Handle(update tgbotapi.Update) error {
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
	}).Info("Duration callback received")

	if len(parts) < 3 {
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"data":    data,
		}).Error("Invalid duration callback format")
		handlers.SendErrorMessage(h.bot, chatID, "Неверный формат команды")
		return nil
	}

	workoutType := parts[1]
	durationStr := parts[2]

	duration, err := strconv.Atoi(durationStr)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id":      userID,
			"chat_id":      chatID,
			"duration_str": durationStr,
			"error":        err,
		}).Error("Failed to parse duration")
		handlers.SendErrorMessage(h.bot, chatID, "Неверная продолжительность")
		return nil
	}

	logger.WithFields(logrus.Fields{
		"user_id":      userID,
		"chat_id":      chatID,
		"workout_type": workoutType,
		"duration":     duration,
	}).Info("Processing duration selection")

	logger.WithFields(logrus.Fields{
		"user_id":      userID,
		"chat_id":      chatID,
		"workout_type": workoutType,
		"duration":     duration,
	}).Info("Creating workout with selected duration")

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

	var workoutName string
	switch workoutType {
	case "split":
		workoutName = "Классический сплит"
	case constants.WorkoutTypePushPull:
		workoutName = "Push/Pull/Legs"
	case "fullbody":
		workoutName = "Фулбади"
	case constants.WorkoutTypeCustom:
		workoutName = "Кастомная тренировка"
	default:
		workoutName = "Новая тренировка"
	}

	workout, err := h.workoutsService.CreateWorkout(user.ID, workoutName, workoutType)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id":      userID,
			"chat_id":      chatID,
			"duration":     duration,
			"workout_type": workoutType,
			"error":        err,
		}).Error("Failed to create workout")
		handlers.SendErrorMessage(h.bot, chatID, "Ошибка создания тренировки")
		return nil
	}

	err = h.workoutsService.Database.Model(workout).Update("duration_minutes", duration).Error
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id":    userID,
			"chat_id":    chatID,
			"workout_id": workout.ID,
			"duration":   duration,
			"error":      err,
		}).Error("Failed to update workout duration")
	}

	text := fmt.Sprintf("✅ Тренировка создана!\n\n"+
		"Тип: %s\n"+
		"Продолжительность: %s минут\n\n"+
		"Тренировка добавлена в ваш список.",
		workoutName, durationStr)

	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)

	_, err = h.bot.Send(editMsg)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id":  userID,
			"chat_id":  chatID,
			"duration": duration,
			"error":    err,
		}).Error("Failed to send workout creation confirmation")
	}
	return err
}
