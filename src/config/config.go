package config

import (
	"os"
	"strconv"
	"workouts_bot/pkg/logger"

	"github.com/joho/godotenv"
)

type LoggerConfig struct {
	Level      string
	FilePath   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	Console    bool
}

type Config struct {
	BotToken string
	Logger   LoggerConfig
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		logger.Warn("No .env file found")
	}

	config := &Config{
		BotToken: getEnv("BOT_TOKEN", ""),
		Logger: LoggerConfig{
			Level:      getEnv("LOG_LEVEL", "info"),
			FilePath:   getEnv("LOG_FILE_PATH", ""),
			MaxSize:    getEnvInt("LOG_MAX_SIZE", 50),
			MaxBackups: getEnvInt("LOG_MAX_BACKUPS", 5),
			MaxAge:     getEnvInt("LOG_MAX_AGE", 30),
			Compress:   getEnvBool("LOG_COMPRESS", true),
			Console:    getEnvBool("LOG_CONSOLE", true),
		},
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
