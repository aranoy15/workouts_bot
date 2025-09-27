package keyboards

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	CreateWorkoutMessage = "🏋️ Создать тренировку"
	MyWorkoutsMessage    = "📊 Мои тренировки"
	ExercisesMessage     = "💪 Упражнения"
	SettingsMessage      = "⚙️ Настройки"
)

func CreateMainMenu() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CreateWorkoutMessage),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(MyWorkoutsMessage),
			tgbotapi.NewKeyboardButton(ExercisesMessage),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(SettingsMessage),
		),
	)

	return keyboard
}
