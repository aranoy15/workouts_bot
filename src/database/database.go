package database

import (
	"fmt"
	"time"
	"workouts_bot/src/config"
	"workouts_bot/src/database/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	if cfg.Host == "sqlite" {
		db, err = gorm.Open(sqlite.Open(cfg.DBName), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	} else {
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	}

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.WorkoutType{},
		&models.Exercise{},
		&models.Equipment{},
		&models.Workout{},
		&models.WorkoutExercise{},
		&models.Set{},
		&models.WeightHistory{},
	)
}
