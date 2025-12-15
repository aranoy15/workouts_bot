package bot

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"workouts_bot/pkg/logger"
	"workouts_bot/src/bot/handlers"
	"workouts_bot/src/bot/handlers/callbacks"
	"workouts_bot/src/bot/handlers/messages"
	"workouts_bot/src/bot/keyboards"
	"workouts_bot/src/config"
	"workouts_bot/src/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Bot struct {
	api              *tgbotapi.BotAPI
	messageHandlers  map[string]handlers.Handler
	callbackHandlers map[string]handlers.Handler
	webhookConfig    *config.WebhookConfig
}

func New(botToken string, database *gorm.DB, webhookCfg *config.WebhookConfig) (*Bot, error) {
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
		webhookConfig:    webhookCfg,
	}, nil
}

func (bot *Bot) Start(botContext context.Context, db *gorm.DB) error {
	if bot.webhookConfig != nil && bot.webhookConfig.Enabled {
		return bot.startWebhook(botContext, db)
	}
	return bot.startPolling(botContext, db)
}

func (bot *Bot) startPolling(botContext context.Context, db *gorm.DB) error {
	logger.Info("Starting bot in long polling mode...")

	if err := database.Migrate(db); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	_, _ = bot.api.Request(tgbotapi.DeleteWebhookConfig{})

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

func (bot *Bot) startWebhook(botContext context.Context, db *gorm.DB) error {
	logger.WithFields(logrus.Fields{
		"url":  bot.webhookConfig.URL + bot.webhookConfig.Path,
		"port": bot.webhookConfig.Port,
	}).Info("Starting bot in webhook mode...")

	webhookURL := bot.webhookConfig.URL + bot.webhookConfig.Path
	webhookConfig, err := tgbotapi.NewWebhook(webhookURL)
	if err != nil {
		logger.Error("Failed to create webhook config:", err)
		return fmt.Errorf("failed to create webhook config: %w", err)
	}

	_, err = bot.api.Request(webhookConfig)
	if err != nil {
		logger.Error("Failed to set webhook:", err)
		return fmt.Errorf("failed to set webhook: %w", err)
	}
	logger.Info("Webhook registered successfully")

	updates := bot.api.ListenForWebhook(bot.webhookConfig.Path)

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", bot.webhookConfig.Port),
	}

	http.HandleFunc("/migrate", func(w http.ResponseWriter, r *http.Request) {
		if err := database.Migrate(db); err != nil {
			logger.Error("Failed to migrate database:", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Failed to migrate database"))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("MIGRATING"))
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("PONG"))
	})

	go func() {
		logger.WithField("port", bot.webhookConfig.Port).Info("Starting HTTP server for webhook...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server error:", err)
		}
	}()

	for {
		select {
		case <-botContext.Done():
			logger.Info("Stopping webhook bot...")

			_, _ = bot.api.Request(tgbotapi.DeleteWebhookConfig{})

			if err := server.Shutdown(botContext); err != nil {
				logger.Warn("Error shutting down HTTP server:", err)
			}
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
