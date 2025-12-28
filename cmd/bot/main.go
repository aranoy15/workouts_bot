package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
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
		JSONFormat: config.Logger.JSONFormat,
	}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--health-check" {
		healthCheck()
		return
	}

	cfg, err := config.Load()
	if err != nil {
		log.Println("Failed to load config:", err)
	}

	logger.Init(loggerConfig(cfg))
	logger.Info("Starting workouts bot...")

	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	bot, err := bot.New(cfg.BotToken, db, &cfg.Webhook)
	if err != nil {
		log.Fatal("Failed to create bot:", err)
	}

	botContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := bot.Start(botContext, db); err != nil {
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
	case <-time.After(0 * time.Second):
		logger.Info("Shutdown completed")
	}
}

func healthCheck() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	url := fmt.Sprintf("http://localhost:%s/health", port)

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Health check failed: %v", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Health check failed: status %d", resp.StatusCode)
		os.Exit(1)
	}

	log.Println("Health check passed")
	os.Exit(0)
}
