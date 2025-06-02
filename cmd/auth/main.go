package main

import (
	"github.com/LeoUraltsev/auth-service/internal/app"
	"github.com/LeoUraltsev/auth-service/internal/config"
	"log/slog"
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

	app.NewApp(nil, cfg)
}
