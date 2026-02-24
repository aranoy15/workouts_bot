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

	"workouts_bot/src/config"
	"workouts_bot/src/logger"
)

func loggerConfig(cfg *config.Config) logger.Config {
	return logger.Config{
		Level:      cfg.Logger.Level,
		Console:    cfg.Logger.Console,
		FilePath:   cfg.Logger.FilePath,
		MaxSize:    cfg.Logger.MaxSize,
		MaxBackups: cfg.Logger.MaxBackups,
		MaxAge:     cfg.Logger.MaxAge,
		Compress:   cfg.Logger.Compress,
		JSONFormat: cfg.Logger.JSONFormat,
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

	app, err := InitializeBot()
	if err != nil {
		log.Fatal("Failed to initialize bot:", err)
	}

	botContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Bot.Start(botContext, app.DB); err != nil {
			logger.Error("Bot error:", err)
			cancel()
		}
	}()

	sig := <-signalChan
	logger.Info("Received signal:", sig.String())
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
