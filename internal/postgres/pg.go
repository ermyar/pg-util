package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
}

func (p *Config) ConnString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s", p.User, p.Password, p.Host, p.Port)
}

func (p *Config) Connect(ctx context.Context) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, p.ConnString())
	if err != nil {
		slog.Error("unable to create connection to postgres", "error", err)
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	return pool, nil
}

func Drop(ctx context.Context, pool *pgxpool.Pool, database string) error {
	_, err := pool.Exec(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", database))
	if err != nil {
		slog.Warn("unable to drop", "database", database, "err", err)
	} else {
		slog.Info("database removed", "database", database)
	}
	return err
}

func backupFile(database string) string {
	return fmt.Sprintf("backup_%s.sql", database)
}

func Backup(ctx context.Context, pool *pgxpool.Pool, database string) error {
	outputFile := backupFile(database)
	connStr := fmt.Sprintf("%s/%s", pool.Config().ConnString(), database)

	cmd := exec.CommandContext(ctx, "pg_dump",
		"-d", connStr,
		"-f", outputFile)

	err := cmd.Run()
	if err != nil {
		slog.Warn("unable to backup", "database", database, "err", err)
		return err
	}
	slog.Info("succesfully dumped", "database", database, "file", outputFile)
	return nil
}

func ListDatabases(ctx context.Context, conn *pgxpool.Pool) ([]string, error) {
	var databases []string
	rows, err := conn.Query(ctx, "SELECT datname FROM pg_database")
	if err != nil {
		slog.Error("unable to list databases", "error", err)
		return nil, err
	}

	for rows.Next() {
		var database string
		if err := rows.Scan(&database); err != nil {
			slog.Error("unable to scan database", "error", err)
			return nil, err
		}
		databases = append(databases, database)
	}
	return databases, nil
}
