package handlers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Handler interface {
	Handle(update tgbotapi.Update) error
}

func GetChatId(update tgbotapi.Update) int64 {
	var chatId int64
	if update.Message != nil {
		chatId = update.Message.Chat.ID
	} else if update.CallbackQuery != nil {
		chatId = update.CallbackQuery.Message.Chat.ID
	}
	return chatId
}
