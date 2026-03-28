package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
}

func (p *Config) connString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s", p.User, p.Password, p.Host, p.Port)
}

func connToDatabase(connStr, database string) string {
	return fmt.Sprintf("%s/%s", connStr, database)
}

func (p *Config) Connect(ctx context.Context) (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(ctx, p.connString())
	if err != nil {
		slog.Error("unable to connect to postgres", "error", err)
		return nil, err
	}
	slog.Info("connected to postgres")
	return pool, nil
}

func Drop(ctx context.Context, pool *pgxpool.Pool, database string) error {
	_, err := pool.Exec(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", database))
	return err
}

func backupFile(database string) string {
	return fmt.Sprintf("backup_%s.sql", database)
}

func Backup(ctx context.Context, pool *pgxpool.Pool, database string) error {
	outputFile := backupFile(database)
	cmd := exec.CommandContext(ctx, "pg_dump",
		"-d", connToDatabase(pool.Config().ConnString(), database),
		"-f", outputFile)
	err := cmd.Run()
	if err != nil {
		slog.Warn("unable to backup", "database", database, "err", err)
		return err
	}
	slog.Info("succesfully dumped", "database", database, "file", outputFile)
	return nil
}
