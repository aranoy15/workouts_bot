package bot

import (
	"context"
	"workouts_bot/pkg/logger"
	"workouts_bot/src/bot/handlers"
	"workouts_bot/src/bot/keyboards"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Bot struct {
	api      *tgbotapi.BotAPI
	handlers map[string]handlers.Handler
}

func New(botToken string, database *gorm.DB) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		logger.Error("Failed to create bot API:", err)
		return nil, err
	}

	logger.Info("Bot API created successfully")
	return &Bot{
		api: bot,
		handlers: map[string]handlers.Handler{
			handlers.StartCommand:          handlers.NewStartHandler(bot, database),
			keyboards.CreateWorkoutMessage: handlers.NewCreateWorkoutHandler(bot),
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
	} else if update.CallbackQuery != nil {
		bot.handleCallbackQuery(update)
	}
}

func (bot *Bot) handleMessage(update tgbotapi.Update) {
	message := update.Message

	logger.WithFields(logrus.Fields{
		"user_id": message.From.ID,
		"chat_id": message.Chat.ID,
		"message": message.Text,
	}).Info("Message:")

	var err error
	switch message.Text {
	case handlers.StartCommand:
		err = bot.handlers[handlers.StartCommand].Handle(update)
	case keyboards.CreateWorkoutMessage:
		err = bot.handlers[keyboards.CreateWorkoutMessage].Handle(update)
	default:
	}

	if err != nil {
		logger.WithFields(logrus.Fields{
			"user_id": message.From.ID,
			"chat_id": message.Chat.ID,
			"message": message.Text,
			"error":   err,
		}).Error("Failed to handle message")
	}
}

func (bot *Bot) handleCallbackQuery(update tgbotapi.Update) {
	callbackQuery := update.CallbackQuery

	logger.WithFields(logrus.Fields{
		"user_id": callbackQuery.From.ID,
		"chat_id": callbackQuery.Message.Chat.ID,
		"data":    callbackQuery.Data,
	}).Info("Callback query:")

}
