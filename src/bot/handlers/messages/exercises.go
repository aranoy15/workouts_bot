package messages

import (
	"fmt"
	"workouts_bot/pkg/logger"
	"workouts_bot/src/bot/handlers"
	"workouts_bot/src/bot/keyboards"
	"workouts_bot/src/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const maxExercisesPerCategory = 5

type ExercisesHandler struct {
	bot             *tgbotapi.BotAPI
	exerciseService *services.ExerciseService
	userService     *services.UserService
}

func NewExercisesHandler(
	bot *tgbotapi.BotAPI,
	database *gorm.DB,
) *ExercisesHandler {
	return &ExercisesHandler{
		bot:             bot,
		exerciseService: services.NewExerciseService(database),
		userService:     services.NewUserService(database),
	}
}

func (handler *ExercisesHandler) Handle(update tgbotapi.Update) error {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	logger.WithFields(logrus.Fields{
		"chat_id": chatID,
		"user_id": userID,
	}).Info("Exercises handler")

	user, err := handler.userService.GetByTelegramID(userID)
	if err != nil {
		logger.Error("Failed to get user by telegram ID: ", err)
		handlers.SendErrorMessage(
			handler.bot, chatID,
			"Ошибка при получении пользователя",
		)
		return err
	}

	exercises, err := handler.exerciseService.GetByEquipment(user.EquipmentIDs)
	if err != nil {
		logger.Error("Failed to get exercises by equipment: ", err)
		handlers.SendErrorMessage(
			handler.bot, chatID,
			"Ошибка при получении упражнений",
		)
		return err
	}

	if len(exercises) == 0 {
		message := "💪 Упражнения не найдены.\n\n" +
			"Возможно, нужно настроить оборудование в настройках."

		msg := tgbotapi.NewMessage(chatID, message)
		_, err = handler.bot.Send(msg)
		return err
	}

	categories := make(map[string][]string)
	for index := range exercises {
		exercise := &exercises[index]
		categories[exercise.Category] = append(
			categories[exercise.Category],
			exercise.Name,
		)
	}

	message := "💪 Доступные упражнения:\n\n"
	for category, exerciseNames := range categories {
		message += fmt.Sprintf("📂 %s:\n", getCategoryEmoji(category))
		for index, name := range exerciseNames {
			if index < maxExercisesPerCategory {
				message += fmt.Sprintf("   • %s\n", name)
			}
		}

		if len(exerciseNames) > maxExercisesPerCategory {
			message += fmt.Sprintf(
				"   ... и еще %d упражнений\n",
				len(exerciseNames)-maxExercisesPerCategory,
			)
		}
		message += "\n"
	}

	message += fmt.Sprintf("Всего упражнений: %d\n\n", len(exercises))
	message += "Используйте кнопки ниже для навигации:"

	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyMarkup = keyboards.CreateExercisesKeyboard()

	_, err = handler.bot.Send(msg)
	return err
}

func getCategoryEmoji(category string) string {
	switch category {
	case "compound":
		return "🏋️ Базовые"
	case "isolation":
		return "🎯 Изолированные"
	case "strength":
		return "💪 Силовые"
	case "cardio":
		return "🏃 Кардио"
	case "bodyweight":
		return "🤸 С собственным весом"
	case "hiit":
		return "⚡ HIIT"
	case "endurance":
		return "🔄 Выносливость"
	default:
		return "📋 " + category
	}
}
