package handlers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type CreateWorkoutHandler struct {
	bot *tgbotapi.BotAPI
}

func NewCreateWorkoutHandler(bot *tgbotapi.BotAPI) *CreateWorkoutHandler {
	return &CreateWorkoutHandler{
		bot: bot,
	}
}

func (h *CreateWorkoutHandler) Handle(update tgbotapi.Update) {
	// TODO: implement create workout handler
	chatId := GetChatId(update)

	message := tgbotapi.NewMessage(chatId, "Создание тренировки")
	h.bot.Send(message)
}
