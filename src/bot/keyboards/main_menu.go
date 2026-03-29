package keyboards

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	StartMessage    = "/start"
	SettingsMessage = "⚙️ Настройки"
)

func CreateMainMenu() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(SettingsMessage),
		),
	)

	return keyboard
}
