//go:build wireinject
// +build wireinject

package main

import (
	"workouts_bot/src/bot"
	"workouts_bot/src/config"
	"workouts_bot/src/database"

	"github.com/google/wire"
	"gorm.io/gorm"
)

func provideDatabaseConfig(cfg *config.Config) *config.DatabaseConfig {
	return &cfg.Database
}

func provideBotToken(cfg *config.Config) string {
	return cfg.BotToken
}

func provideWebhookConfig(cfg *config.Config) *config.WebhookConfig {
	return &cfg.Webhook
}

type BotApp struct {
	Bot *bot.Bot
	DB  *gorm.DB
}

func InitializeBot() (*BotApp, error) {
	wire.Build(
		config.Load,
		provideDatabaseConfig,
		database.Connect,
		provideBotToken,
		provideWebhookConfig,
		bot.New,
		wire.Struct(new(BotApp), "Bot", "DB"),
	)
	return nil, nil
}
