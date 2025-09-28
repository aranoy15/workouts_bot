package handlers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Handler interface {
	Handle(update tgbotapi.Update) error
}

func SendErrorMessage(bot *tgbotapi.BotAPI, chatID int64, errorText string) {
	msg := tgbotapi.NewMessage(chatID, "❌ "+errorText)
	bot.Send(msg)
}
