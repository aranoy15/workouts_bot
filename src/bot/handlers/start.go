package handlers

import (
	"workouts_bot/src/bot/keyboards"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	helloMessage = "Привет! Я бот для тренировок 🏋️\n\n" +
		"Я помогу тебе:\n" +
		"• Подобрать упражнения\n" +
		"• Составить программу тренировок\n" +
		"• Отслеживать прогресс\n" +
		"• Записывать подходы и веса\n\n" +
		"Выбери действие:"
)

type StartHandler struct {
	bot *tgbotapi.BotAPI
}

func NewStartHandler(bot *tgbotapi.BotAPI) *StartHandler {
	return &StartHandler{
		bot: bot,
	}
}

func (startHandler *StartHandler) Handle(update tgbotapi.Update) {
	keyboard := keyboards.CreateMainMenu()

	chatId := GetChatId(update)

	message := tgbotapi.NewMessage(chatId, helloMessage)
	message.ReplyMarkup = &keyboard
	startHandler.bot.Send(message)
}
