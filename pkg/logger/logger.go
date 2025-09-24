package logger

import (
	"io"
	"os"

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
}

func Init(config Config) {
	Log = logrus.New()
	Log.SetOutput(os.Stdout)

	logLevel, err := logrus.ParseLevel(config.Level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}

	Log.SetLevel(logLevel)
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

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
