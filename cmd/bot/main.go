package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"workouts_bot/pkg/logger"
	"workouts_bot/src/bot"
	"workouts_bot/src/config"
	"workouts_bot/src/database"
)

func loggerConfig(config *config.Config) logger.Config {
	return logger.Config{
		Level:      config.Logger.Level,
		Console:    config.Logger.Console,
		FilePath:   config.Logger.FilePath,
		MaxSize:    config.Logger.MaxSize,
		MaxBackups: config.Logger.MaxBackups,
		MaxAge:     config.Logger.MaxAge,
		Compress:   config.Logger.Compress,
	}
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	logger.Init(loggerConfig(cfg))
	logger.Info("Starting workouts bot...")

	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	bot, err := bot.New(cfg.BotToken, db)
	if err != nil {
		log.Fatal("Failed to create bot:", err)
	}

	botContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := bot.Start(botContext); err != nil {
			logger.Error("Bot error:", err)
			cancel()
		}
	}()

	signal := <-signalChan
	logger.Info("Received signal:", signal.String())
	logger.Info("Shutting down gracefully...")

	cancel()

	shutdownTimeout := 10 * time.Second
	shutdownContext, shutdownCancel := context.WithTimeout(
		context.Background(),
		shutdownTimeout,
	)
	defer shutdownCancel()

	select {
	case <-shutdownContext.Done():
		logger.Warn("Shutdown timeout exceeded, forcing exit")
	case <-time.After(5 * time.Second):
		logger.Info("Shutdown completed")
	}
}
