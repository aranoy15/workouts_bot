package bot

import (
	"context"
	"workouts_bot/pkg/logger"
	"workouts_bot/src/bot/handlers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	StartCommand = "/start"
)

type Bot struct {
	api      *tgbotapi.BotAPI
	handlers map[string]handlers.Handler

	//TODO: implement db integration

	//TODO: implement all services

	// handlers
}

func New(botToken string) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		logger.Error("Failed to create bot API:", err)
		return nil, err
	}

	logger.Info("Bot API created successfully")
	return &Bot{
		api: bot,
		handlers: map[string]handlers.Handler{
			StartCommand: handlers.NewStartHandler(bot),
		},
	}, nil
}

func (bot *Bot) Start(botContext context.Context) error {
	logger.Info("Starting bot...")

	botUpdate := tgbotapi.NewUpdate(0)
	botUpdate.Timeout = 60

	updates := bot.api.GetUpdatesChan(botUpdate)
	for {
		select {
		case <-botContext.Done():
			logger.Info("Stopping bot...")
			return nil
		case update := <-updates:
			go bot.handleUpdate(update)
		}
	}
}

func (bot *Bot) handleUpdate(update tgbotapi.Update) {
	if update.Message != nil {
		bot.handleMessage(update)
	}
}

func (bot *Bot) handleMessage(update tgbotapi.Update) {
	message := update.Message

	switch message.Text {
	case StartCommand:
		bot.handlers[StartCommand].Handle(update)
	default:
	}

	/*
		logger.WithFields(logrus.Fields{
			"user_id":  message.From.ID,
			"username": message.From.UserName,
			"text":     message.Text,
		}).Info("Received message")

		answer := tgbotapi.NewMessage(message.Chat.ID, message.Text)
		bot.api.Send(answer)
	*/
}
