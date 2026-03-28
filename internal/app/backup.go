package app

import (
	"context"
	"log/slog"

	"github.com/ermyar/pg-util/internal/postgres"
)

func (a *App) backup() error {
	slog.Info("backup started")
	for _, database := range a.databases {
		postgres.Backup(context.Background(), a.pool, database)
	}
	slog.Info("backup finished")
	return nil
}
