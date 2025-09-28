package bot

import (
	"context"
	"strings"
	"workouts_bot/pkg/logger"
	"workouts_bot/src/bot/handlers"
	"workouts_bot/src/bot/handlers/callbacks"
	"workouts_bot/src/bot/handlers/messages"
	"workouts_bot/src/bot/keyboards"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Bot struct {
	api              *tgbotapi.BotAPI
	messageHandlers  map[string]handlers.Handler
	callbackHandlers map[string]handlers.Handler
}

func New(botToken string, database *gorm.DB) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		logger.Error("Failed to create bot API:", err)
		return nil, err
	}

	logger.Info("Bot API created successfully")
	messageHandlers := map[string]handlers.Handler{
		keyboards.StartMessage: messages.NewStartHandler(
			bot, database,
		),
		keyboards.CreateWorkoutMessage: messages.NewCreateWorkoutHandler(
			bot, database,
		),
		keyboards.MyWorkoutsMessage: messages.NewMyWorkoutsHandler(
			bot, database,
		),
		keyboards.ExercisesMessage: messages.NewExercisesHandler(
			bot, database,
		),
		keyboards.SettingsMessage: messages.NewSettingsHandler(
			bot, database,
		),
	}

	callbackHandlers := map[string]handlers.Handler{
		callbacks.WorkoutCallbackType: callbacks.NewWorkoutHandler(
			bot, database,
		),
		callbacks.ExerciseCallbackType: callbacks.NewExerciseHandler(
			bot, database,
		),
		callbacks.SetCallbackType: callbacks.NewSetHandler(
			bot, database,
		),
		callbacks.GoalCallbackType: callbacks.NewGoalHandler(
			bot, database,
		),
		callbacks.EquipmentCallbackType: callbacks.NewEquipmentHandler(
			bot, database,
		),
		callbacks.ExperienceCallbackType: callbacks.NewExperienceHandler(
			bot, database,
		),
		callbacks.SettingsCallbackType: callbacks.NewSettingsHandler(
			bot, database,
		),
		callbacks.WorkoutTypeCallbackType: callbacks.NewWorkoutTypeHandler(
			bot, database,
		),
		callbacks.DurationCallbackType: callbacks.NewDurationHandler(
			bot, database,
		),
	}

	return &Bot{
		api:              bot,
		messageHandlers:  messageHandlers,
		callbackHandlers: callbackHandlers,
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

	handler, ok := bot.messageHandlers[message.Text]
	if !ok {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Invalid command")
		_, _ = bot.api.Send(msg)
		return
	}

	if err := handler.Handle(update); err != nil {
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

	callback := tgbotapi.NewCallback(callbackQuery.ID, "")
	_, _ = bot.api.Request(callback)

	logger.WithFields(logrus.Fields{
		"user_id": callbackQuery.From.ID,
		"chat_id": callbackQuery.Message.Chat.ID,
		"data":    callbackQuery.Data,
	}).Info("Callback query:")

	data := callbackQuery.Data
	parts := strings.Split(data, ":")

	if len(parts) < 1 {
		msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "❌ Неверный формат команды")
		_, _ = bot.api.Send(msg)
		return
	}

	handlerType := parts[0]
	handler, ok := bot.callbackHandlers[handlerType]
	if !ok {
		msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "❌ Неизвестная команда")
		_, _ = bot.api.Send(msg)
		return
	}

	if err := handler.Handle(update); err != nil {
		logger.WithFields(logrus.Fields{
			"user_id": callbackQuery.From.ID,
			"chat_id": callbackQuery.Message.Chat.ID,
			"data":    callbackQuery.Data,
			"error":   err,
		}).Error("Failed to handle callback query")
	}
}
