package keyboards

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CreateSettingsKeyboard() tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				SettingsExperience,
				"settings:experience",
			),
		),
	)

	return keyboard
}

func CreateExperienceLevelKeyboard() tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				ExpBeginner,
				"experience:0",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				ExpIntermediate,
				"experience:1",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				ExpAdvanced,
				"experience:3",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				ExpExpert,
				"experience:5",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				NavBack,
				"settings:experience_back",
			),
		),
	)

	return keyboard
}

func CreateConfirmationKeyboard(action string) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				NavYes,
				fmt.Sprintf("confirm:%s:yes", action),
			),
			tgbotapi.NewInlineKeyboardButtonData(
				NavNo,
				fmt.Sprintf("confirm:%s:no", action),
			),
		),
	)

	return keyboard
}

func CreateBackKeyboard(callbackData string) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				NavBack,
				callbackData,
			),
		),
	)

	return keyboard
}
