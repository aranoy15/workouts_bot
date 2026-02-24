package callbacks

import (
	"fmt"
	"strconv"
	"strings"
	"workouts_bot/src/logger"
	"workouts_bot/src/bot/handlers"
	"workouts_bot/src/bot/keyboards"
	"workouts_bot/src/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const ExerciseCallbackType = "exercise"

type ExerciseHandler struct {
	bot             *tgbotapi.BotAPI
	exerciseService *services.ExerciseService
	workoutsService *services.WorkoutsService
	userService     *services.UserService
}

func NewExerciseHandler(bot *tgbotapi.BotAPI, database *gorm.DB) *ExerciseHandler {
	return &ExerciseHandler{
		bot:             bot,
		exerciseService: services.NewExerciseService(database),
		workoutsService: services.NewWorkoutsService(database),
		userService:     services.NewUserService(database),
	}
}

func (h *ExerciseHandler) Handle(update tgbotapi.Update) error {
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
	}).Info("Exercise callback received")

	if len(parts) < 3 {
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"data":    data,
		}).Error("Invalid exercise callback format")
		handlers.SendErrorMessage(h.bot, chatID, "Неверный формат команды")
		return nil
	}

	action := parts[1]
	exerciseID, err := strconv.ParseUint(parts[2], 10, 32)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id":     userID,
			"chat_id":     chatID,
			"exercise_id": parts[2],
			"error":       err,
		}).Error("Failed to parse exercise ID")
		handlers.SendErrorMessage(h.bot, chatID, "Неверный ID упражнения")
		return nil
	}

	switch action {
	case "details":
		return h.showDetails(userID, chatID, messageID, uint(exerciseID))
	case "video":
		return h.showVideo(userID, chatID, messageID, uint(exerciseID))
	case "add":
		return h.addToWorkout(userID, chatID, messageID, uint(exerciseID))
	default:
		logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
			"action":  action,
		}).Error("Unknown exercise action")
		handlers.SendErrorMessage(h.bot, chatID, "Неизвестное действие")
		return nil
	}
}

func (h *ExerciseHandler) showDetails(
	userID int64,
	chatID int64,
	messageID int,
	exerciseID uint,
) error {
	logger.WithFields(logrus.Fields{
		"user_id":     userID,
		"chat_id":     chatID,
		"exercise_id": exerciseID,
	}).Info("Showing exercise details")

	exercise, err := h.exerciseService.GetByID(exerciseID)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id":     userID,
			"chat_id":     chatID,
			"exercise_id": exerciseID,
			"error":       err,
		}).Error("Failed to get exercise details")
		handlers.SendErrorMessage(h.bot, chatID, "Упражнение не найдено")
		return nil
	}

	text := fmt.Sprintf(
		"📖 %s\n\n%s\n\n💪 Группы мышц: %s",
		exercise.Name,
		exercise.Description,
		strings.Join(exercise.MuscleGroups, ", "),
	)

	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	keyboard := keyboards.CreateExerciseKeyboard(exerciseID)
	editMsg.ReplyMarkup = &keyboard

	_, err = h.bot.Send(editMsg)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id":     userID,
			"chat_id":     chatID,
			"exercise_id": exerciseID,
			"error":       err,
		}).Error("Failed to send exercise details message")
	}
	return err
}

func (h *ExerciseHandler) showVideo(
	userID int64,
	chatID int64,
	_ int,
	exerciseID uint,
) error {
	logger.WithFields(logrus.Fields{
		"user_id":     userID,
		"chat_id":     chatID,
		"exercise_id": exerciseID,
	}).Info("Showing exercise video")

	exercise, err := h.exerciseService.GetByID(exerciseID)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id":     userID,
			"chat_id":     chatID,
			"exercise_id": exerciseID,
			"error":       err,
		}).Error("Failed to get exercise for video")
		handlers.SendErrorMessage(h.bot, chatID, "Упражнение не найдено")
		return nil
	}

	if exercise.VideoPath == "" {
		logger.WithFields(logrus.Fields{
			"user_id":     userID,
			"chat_id":     chatID,
			"exercise_id": exerciseID,
		}).Error("Exercise has no video")
		handlers.SendErrorMessage(h.bot, chatID, "Видео для этого упражнения не найдено")
		return nil
	}

	video := tgbotapi.NewVideo(chatID, tgbotapi.FilePath(exercise.VideoPath))
	video.Caption = fmt.Sprintf("🎥 %s", exercise.Name)

	_, err = h.bot.Send(video)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id":     userID,
			"chat_id":     chatID,
			"exercise_id": exerciseID,
			"error":       err,
		}).Error("Failed to send exercise video")
	}
	return err
}

func (h *ExerciseHandler) addToWorkout(
	userID int64,
	chatID int64,
	messageID int,
	exerciseID uint,
) error {
	logger.WithFields(logrus.Fields{
		"user_id":     userID,
		"chat_id":     chatID,
		"exercise_id": exerciseID,
	}).Info("Adding exercise to workout")

	user, err := h.userService.GetByTelegramID(userID)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id":     userID,
			"chat_id":     chatID,
			"exercise_id": exerciseID,
			"error":       err,
		}).Error("Failed to get user by telegram ID")
		handlers.SendErrorMessage(h.bot, chatID, "Ошибка при получении пользователя")
		return nil
	}

	workouts, err := h.workoutsService.GetUserWorkouts(user.ID)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id":     userID,
			"chat_id":     chatID,
			"exercise_id": exerciseID,
			"error":       err,
		}).Error("Failed to get user workouts")
		handlers.SendErrorMessage(h.bot, chatID, "Ошибка получения тренировок")
		return nil
	}

	if len(workouts) == 0 {
		logger.WithFields(logrus.Fields{
			"user_id":     userID,
			"chat_id":     chatID,
			"exercise_id": exerciseID,
		}).Error("User has no workouts")
		handlers.SendErrorMessage(h.bot, chatID, "У вас нет тренировок. Создайте тренировку сначала.")
		return nil
	}

	text := "📝 Выберите тренировку для добавления упражнения:"

	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)

	_, err = h.bot.Send(editMsg)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id":     userID,
			"chat_id":     chatID,
			"exercise_id": exerciseID,
			"error":       err,
		}).Error("Failed to send workout selection message")
	}
	return err
}
