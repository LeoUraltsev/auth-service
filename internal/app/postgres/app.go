package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"time"
)

type Pool interface {
	Close()
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	Ping(ctx context.Context) error
}
type Postgres struct {
	Pool Pool
	log  *slog.Logger
}

//todo: подумать над конфигурацией
//todo: подумать над тем чтобы не сразу возвращать ошибку, а дать несколько раз подключится

func NewPostgresPool(ctx context.Context, log *slog.Logger, connString string) (*Postgres, error) {
	log.Debug("connecting to postgres", slog.String("connection_string", connString))
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	config.MaxConns = 10
	config.MaxConnLifetime = 5 * time.Minute
	config.MaxConnIdleTime = 10 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Error("failed to connect to postgres", slog.String("error", err.Error()))
		return nil, err
	}
	log.Debug("connected to postgres", slog.String("connection_string", connString))
	err = pool.Ping(ctx)
	if err != nil {
		log.Error("failed to ping postgres", slog.String("error", err.Error()))
		return nil, err
	}
	log.Info("success ping to postgres")
	return &Postgres{
		Pool: pool,
		log:  log,
	}, nil
}

func (p *Postgres) Close() {
	p.log.Info("shutting down postgres")
	p.Pool.Close()
	p.log.Info("postgres pool closed")
}
