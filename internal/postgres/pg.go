package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PgConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

func (p *PgConfig) connString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s", p.User, p.Password, p.Host, p.Port)
}

func (p *PgConfig) Connect(ctx context.Context) (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(ctx, p.connString())
	if err != nil {
		slog.Error("unable to connect to postgres", "error", err)
		return nil, err
	}
	slog.Info("connected to postgres")
	return pool, nil
}
