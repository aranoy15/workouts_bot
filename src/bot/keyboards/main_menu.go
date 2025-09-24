package keyboards

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func CreateMainMenu() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
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
	)

	keyboard.ResizeKeyboard = true
	keyboard.OneTimeKeyboard = true

	return keyboard
}
