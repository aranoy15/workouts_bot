package keyboards

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func CreateMainMenu() tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏋️ Создать тренировку", "create_workout"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Мои тренировки", "my_workouts"),
			tgbotapi.NewInlineKeyboardButtonData("💪 Упражнения", "exercises"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⚙️ Настройки", "settings"),
		),
	/*
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🏋️ Создать тренировку"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📊 Мои тренировки"),
			tgbotapi.NewKeyboardButton("💪 Упражнения"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("⚙️ Настройки"),
		),
	*/
	)

	return keyboard
}
