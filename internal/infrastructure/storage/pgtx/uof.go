package pgtx

import (
	"context"
	pg "github.com/LeoUraltsev/auth-service/internal/app/postgres"
	"github.com/LeoUraltsev/auth-service/internal/domain/users"
	"github.com/LeoUraltsev/auth-service/internal/helper/logger"
	"github.com/LeoUraltsev/auth-service/internal/infrastructure/storage/postgres"
	"log/slog"
)

type StorageUnitOfWork struct {
	pg  *pg.Postgres
	log *slog.Logger
}

func NewStorageUnitOfWork(pg *pg.Postgres, log *slog.Logger) *StorageUnitOfWork {
	return &StorageUnitOfWork{
		pg:  pg,
		log: log,
	}
}

func (s *StorageUnitOfWork) Execute(ctx context.Context, fn func(repository users.UserRepository) error) error {
	log := logger.LogWithContext(ctx, s.log)
	log.Info("starting transaction")
	tx, err := s.pg.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	
	repo := postgres.NewUsersStorage(s.pg, log)

	if err = fn(repo); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
