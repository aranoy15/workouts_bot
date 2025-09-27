package handlers

import (
	"workouts_bot/pkg/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type CreateWorkoutHandler struct {
	bot *tgbotapi.BotAPI
}

func NewCreateWorkoutHandler(bot *tgbotapi.BotAPI) *CreateWorkoutHandler {
	return &CreateWorkoutHandler{
		bot: bot,
	}
}

func (h *CreateWorkoutHandler) Handle(update tgbotapi.Update) error {
	// TODO: implement create workout handler
	chatId := GetChatId(update)

	logger.WithFields(logrus.Fields{
		"chat_id": chatId,
	}).Info("Create workout handler")

	message := tgbotapi.NewMessage(chatId, "Создание тренировки")

	_, err := h.bot.Send(message)
	return err
}
