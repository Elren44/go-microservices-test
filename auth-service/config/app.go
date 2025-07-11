package config

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

type App struct {
	Config *AuthConfig
	Logger *slog.Logger
}

func NewApp(config *AuthConfig) *App {
	logger := slog.New()
	h1 := handler.NewConsoleHandler(slog.AllLevels)
	logger.AddHandlers(h1)
	return &App{
		Config: config,
		Logger: logger,
	}
}
