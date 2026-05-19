//go:build wireinject
// +build wireinject

package main

import (
	"workouts_bot/src/config"
	"workouts_bot/src/router"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type ServiceApp struct {
	Engine *gin.Engine
	Cfg    *config.Config
}

func NewServiceApp(cfg *config.Config) *ServiceApp {
	return &ServiceApp{
		Cfg:    cfg,
		Engine: router.NewRouter(cfg),
	}
}

func InitializeService() (*ServiceApp, error) {
	wire.Build(
		config.Load,
		NewServiceApp,
	)
	return nil, nil
}
