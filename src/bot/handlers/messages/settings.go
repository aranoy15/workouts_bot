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

type SettingsHandler struct {
	bot         *tgbotapi.BotAPI
	userService *services.UserService
}

func NewSettingsHandler(bot *tgbotapi.BotAPI, database *gorm.DB) *SettingsHandler {
	return &SettingsHandler{
		bot:         bot,
		userService: services.NewUserService(database),
	}
}

func (handler *SettingsHandler) Handle(update tgbotapi.Update) error {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	logger.WithFields(logrus.Fields{
		"chat_id": chatID,
		"user_id": userID,
	}).Info("Settings handler")

	user, err := handler.userService.GetByTelegramID(userID)
	if err != nil {
		logger.Error("Failed to get user by telegram ID: ", err)
		handlers.SendErrorMessage(
			handler.bot, chatID,
			"Ошибка при получении пользователя",
		)
		return err
	}

	message := "⚙️ Ваши настройки:\n\n"

	message += "🎯 Цели:\n"
	if len(user.Goals) == 0 {
		message += "   Не установлены\n"
	} else {
		for _, goal := range user.Goals {
			message += fmt.Sprintf("   • %s\n", getGoalEmoji(string(goal)))
		}
	}
	message += "\n"

	message += fmt.Sprintf("📈 Уровень опыта: %s\n\n", getExperienceLevel(user.Experience))
	message += "🏋️ Оборудование:\n"
	if len(user.EquipmentIDs) == 0 {
		message += "   Не настроено\n"
	} else {
		message += fmt.Sprintf("   Настроено %d типов оборудования\n", len(user.EquipmentIDs))
	}
	message += "\n"

	message += "⚠️ Ограничения:\n"
	if len(user.Limitations) == 0 {
		message += "   Не указаны\n"
	} else {
		for _, limitation := range user.Limitations {
			message += fmt.Sprintf("   • %s\n", limitation)
		}
	}
	message += "\n"

	message += "Используйте кнопки ниже для изменения настроек:"

	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyMarkup = keyboards.CreateSettingsKeyboard()

	_, err = handler.bot.Send(msg)
	return err
}

func getGoalEmoji(goal string) string {
	switch goal {
	case "muscle_gain":
		return "💪 Набор массы"
	case "strength":
		return "🏋️ Сила"
	case "endurance":
		return "🏃 Выносливость"
	case "weight_loss":
		return "🔥 Похудение"
	default:
		return "📋 " + goal
	}
}

func getExperienceLevel(experience int) string {
	switch {
	case experience <= 1:
		return "🟢 Начинающий"
	case experience <= 3:
		return "🟡 Средний"
	case experience <= 5:
		return "🟠 Продвинутый"
	default:
		return "🔴 Эксперт"
	}
}
