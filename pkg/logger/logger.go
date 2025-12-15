package logger

import (
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *logrus.Logger

type Config struct {
	Level      string
	FilePath   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	Console    bool
	JSONFormat bool
}

type YandexCloudFormatter struct {
	logrus.JSONFormatter
}

func (f *YandexCloudFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	entry.Data["level"] = strings.ToUpper(entry.Level.String())
	entry.Data["message"] = entry.Message
	entry.Data["timestamp"] = entry.Time.Format("2006-01-02T15:04:05.000Z07:00")

	// Delete duplicate fields
	delete(entry.Data, "msg")
	delete(entry.Data, "time")

	return f.JSONFormatter.Format(entry)
}

func Init(config Config) {
	Log = logrus.New()
	Log.SetOutput(os.Stdout)

	logLevel, err := logrus.ParseLevel(config.Level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}

	Log.SetLevel(logLevel)

	if config.JSONFormat {
		Log.SetFormatter(&YandexCloudFormatter{
			JSONFormatter: logrus.JSONFormatter{
				TimestampFormat:   "2006-01-02T15:04:05.000Z07:00",
				DisableTimestamp:  true,
				DisableHTMLEscape: true,
				FieldMap: logrus.FieldMap{
					logrus.FieldKeyLevel: "severity",
					logrus.FieldKeyMsg:   "message",
					logrus.FieldKeyTime:  "timestamp",
				},
			},
		})
	} else {
		Log.SetFormatter(&logrus.TextFormatter{
			DisableTimestamp: true,
		})
	}

	var writers []io.Writer

	if config.Console {
		writers = append(writers, os.Stdout)
	}

	if config.FilePath != "" {
		fileWriter := &lumberjack.Logger{
			Filename:   config.FilePath,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		}
		writers = append(writers, fileWriter)
	}
	if len(writers) > 1 {
		Log.SetOutput(io.MultiWriter(writers...))
	} else if len(writers) == 1 {
		Log.SetOutput(writers[0])
	} else {
		Log.SetOutput(os.Stdout)
	}
}

func InitSimple(level string) {
	config := Config{
		Level:      level,
		Console:    true,
		FilePath:   "",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}
	Init(config)
}

func Info(args ...any) {
	Log.Info(args...)
}

func Error(args ...any) {
	Log.Error(args...)
}

func Debug(args ...any) {
	Log.Debug(args...)
}

func Warn(args ...any) {
	Log.Warn(args...)
}

func Fatal(args ...any) {
	Log.Fatal(args...)
}

func WithField(key string, value any) *logrus.Entry {
	return Log.WithField(key, value)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return Log.WithFields(fields)
}
