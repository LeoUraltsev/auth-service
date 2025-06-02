package app

import (
	"github.com/LeoUraltsev/auth-service/internal/config"
	"log/slog"
)

/*
1. Конфигурация
2. Логгирование
3. Модели
4. Storage
5. grpc
6. application service
*/

type App struct {
	log *slog.Logger
	cfg config.Config
}

func NewApp(log *slog.Logger, cfg config.Config) *App {
	return &App{
		log: log,
		cfg: cfg,
	}
}

func (a *App) Run() error {

	return nil
}
