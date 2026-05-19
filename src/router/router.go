package router

import (
	"log"
	"workouts_bot/src/config"

	"github.com/gin-gonic/gin"
)

const (
	defaultPath = "/api/v1"
)

func NewRouter(cfg *config.Config) *gin.Engine {
	engine := gin.New()

	log.Printf("Service started port %d", cfg.Port)

	engine.Use(gin.Logger())

	return engine
}
