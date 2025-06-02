package main

import (
	"github.com/LeoUraltsev/auth-service/internal/app"
	"github.com/LeoUraltsev/auth-service/internal/app/logger"
	"github.com/LeoUraltsev/auth-service/internal/config"
	"log/slog"
	"os"
)

func main() {
	cfg, err := config.NewConfig("./config/config.local.yaml", "prod.env")
	if err != nil {
		slog.Warn(err.Error())
	}

	slog.Info(
		"loaded config",
		slog.String("env", cfg.App.Env),
		slog.String("grpc", cfg.GRPC.Address),
		slog.String("postgres", cfg.Postgres.DSN),
	)

	log, err := logger.NewLogger(cfg.App.Env)
	if err != nil {
		slog.Error("failed to initialize logger", slog.String("err", err.Error()))
		os.Exit(1)
	}

	log.Log.With(slog.String("env", cfg.App.Env))
	log.Log.Info("initialized logger")

	app.NewApp(
		log.Log,
		cfg,
	)
}
