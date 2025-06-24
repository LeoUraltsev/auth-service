package app

import (
	"context"
	"github.com/LeoUraltsev/auth-service/internal/app/grpc"
	"github.com/LeoUraltsev/auth-service/internal/app/postgres"
	"github.com/LeoUraltsev/auth-service/internal/application"
	"github.com/LeoUraltsev/auth-service/internal/config"
	"github.com/LeoUraltsev/auth-service/internal/infrastructure/hasher"
	pgUserStorage "github.com/LeoUraltsev/auth-service/internal/infrastructure/storage/postgres"
	"log/slog"
	"os/signal"
	"sync"
	"syscall"
)

type App struct {
	log *slog.Logger
	cfg *config.Config
}

func NewApp(log *slog.Logger, cfg *config.Config) *App {
	return &App{
		log: log,
		cfg: cfg,
	}
}

func (a *App) Run() error {
	log := a.log
	log.Info("starting app")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pg, err := postgres.NewPostgresPool(ctx, log, a.cfg.Postgres.DSN)
	if err != nil {
		log.Info("failed to connect database")
		return err
	}

	userStorage := pgUserStorage.NewUsersStorage(pg, log)
	hash := hasher.NewHasher()

	userService := application.NewUserService(userStorage, hash, log)

	rpc := grpc.NewApp(userService, log, a.cfg.GRPC.Address)

	chErrRpc := make(chan error)
	go func() {
		defer close(chErrRpc)
		err = rpc.Start()
		if err != nil {
			chErrRpc <- err
		}
	}()

	wg := &sync.WaitGroup{}

	select {
	case <-ctx.Done():
		log.Info("shutting down app")
		wg.Add(2)
		go func() {
			defer wg.Done()
			rpc.Stop()
		}()
		go func() {
			defer wg.Done()
			pg.Close()
		}()
		wg.Wait()
		log.Info("app stopped")
	case err := <-chErrRpc:
		if err != nil {
			pg.Close()
			return err
		}
	}

	return nil
}
