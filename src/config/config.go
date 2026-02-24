package config

import (
	"net/url"
	"os"
	"strconv"
	"strings"
)

type LoggerConfig struct {
	Level      string
	FilePath   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	Console    bool
	JSONFormat bool
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type WebhookConfig struct {
	Enabled     bool
	URL         string
	Path        string
	Port        int
	SecretToken string
}

type S3Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	Region          string
}

type Config struct {
	BotToken string
	Logger   LoggerConfig
	Database DatabaseConfig
	Webhook  WebhookConfig
	S3       S3Config
}

func Load() (*Config, error) {
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
			JSONFormat: getEnvBool("LOG_JSON_FORMAT", false),
		},
		Database: parseDatabaseConfig(),
		Webhook: WebhookConfig{
			Enabled:     getEnvBool("WEBHOOK_ENABLED", true),
			URL:         getEnv("WEBHOOK_URL", ""),
			Path:        getEnv("WEBHOOK_PATH", "/webhook"),
			Port:        getEnvInt("PORT", 8080),
			SecretToken: getEnv("WEBHOOK_SECRET_TOKEN", ""),
		},
		S3: S3Config{
			Endpoint:        getEnv("S3_ENDPOINT", "https://storage.yandexcloud.net"),
			AccessKeyID:     getEnv("S3_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("S3_SECRET_ACCESS_KEY", ""),
			BucketName:      getEnv("S3_BUCKET_NAME", ""),
			Region:          getEnv("S3_REGION", "ru-central1"),
		},
	}

	return config, nil
}

func parseDatabaseConfig() DatabaseConfig {
	if databaseURL := getEnv("DATABASE_URL", ""); databaseURL != "" {
		if config, err := parseDatabaseURL(databaseURL); err == nil {
			return config
		}
	}

	return DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnvInt("DB_PORT", 5432),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "postgres"),
		SSLMode:  getEnv("DB_SSL_MODE", "disable"),
	}
}

func parseDatabaseURL(databaseURL string) (DatabaseConfig, error) {
	u, err := url.Parse(databaseURL)
	if err != nil {
		return DatabaseConfig{}, err
	}

	config := DatabaseConfig{
		Host:    u.Hostname(),
		Port:    5432,
		DBName:  strings.TrimPrefix(u.Path, "/"),
		SSLMode: "disable",
	}

	if u.Port() != "" {
		if port, err := strconv.Atoi(u.Port()); err == nil {
			config.Port = port
		}
	}

	if u.User != nil {
		config.User = u.User.Username()
		if password, ok := u.User.Password(); ok {
			config.Password = password
		}
	}

	if sslMode := u.Query().Get("sslmode"); sslMode != "" {
		config.SSLMode = sslMode
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
