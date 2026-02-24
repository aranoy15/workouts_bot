package messages

import (
	"fmt"
	"workouts_bot/src/logger"
	"workouts_bot/src/bot/handlers"
	"workouts_bot/src/bot/keyboards"
	"workouts_bot/src/constants"
	"workouts_bot/src/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MyWorkoutsHandler struct {
	bot             *tgbotapi.BotAPI
	workoutsService *services.WorkoutsService
	userService     *services.UserService
}

func NewMyWorkoutsHandler(
	bot *tgbotapi.BotAPI,
	database *gorm.DB,
) *MyWorkoutsHandler {
	return &MyWorkoutsHandler{
		bot:             bot,
		workoutsService: services.NewWorkoutsService(database),
		userService:     services.NewUserService(database),
	}
}

func (handler *MyWorkoutsHandler) Handle(update tgbotapi.Update) error {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	logger.WithFields(logrus.Fields{
		"chat_id": chatID,
		"user_id": userID,
	}).Info("My workouts handler")

	user, err := handler.userService.GetByTelegramID(userID)
	if err != nil {
		logger.Error("Failed to get user by telegram ID: ", err)
		handlers.SendErrorMessage(
			handler.bot, chatID,
			"Ошибка при получении пользователя",
		)
		return err
	}

	workouts, err := handler.workoutsService.GetUserWorkouts(user.ID)
	if err != nil {
		logger.Error("Failed to get user workouts: ", err)
		handlers.SendErrorMessage(
			handler.bot, chatID,
			"Ошибка при получении тренировок",
		)
		return err
	}

	if len(workouts) == 0 {
		message := "📊 У вас пока нет тренировок.\n\n" +
			"Создайте первую тренировку, нажав кнопку \"🏋️ Создать тренировку\""

		msg := tgbotapi.NewMessage(chatID, message)
		_, err = handler.bot.Send(msg)
		return err
	}

	message := "📊 Ваши тренировки:\n\n"

	for index := range workouts {
		workout := &workouts[index]
		status := getWorkoutStatusEmoji(workout.Status)
		message += fmt.Sprintf("%d. %s %s\n", index+1, status, workout.Name)

		if workout.Status == constants.WorkoutStatusCompleted && workout.CompletedAt != nil {
			message += fmt.Sprintf("   ✅ Завершена: %s\n",
				workout.CompletedAt.Format("02.01.2006 15:04"),
			)
		} else if workout.Status == constants.WorkoutStatusInProgress {
			message += "   🔄 В процессе\n"
		} else {
			message += "   📋 Запланирована\n"
		}

		message += fmt.Sprintf("   ⏱️ Длительность: %d мин\n", workout.DurationMinutes)
		message += fmt.Sprintf("   💪 Упражнений: %d\n\n", len(workout.Exercises))
	}

	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyMarkup = keyboards.CreateMyWorkoutsKeyboard()

	_, err = handler.bot.Send(msg)
	return err
}

func getWorkoutStatusEmoji(status string) string {
	switch status {
	case "completed":
		return "✅"
	case "in_progress":
		return "🔄"
	case "planned":
		return "📋"
	default:
		return "❓"
	}
}
