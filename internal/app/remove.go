package app

import (
	"context"
	"log/slog"

	"github.com/ermyar/pg-util/internal/postgres"
)

func (a *App) remove() error {
	slog.Info("remove started")
	defer slog.Info("remove finished")

	for _, database := range a.databases {
		if err := postgres.Drop(context.Background(), a.pool, database); err != nil {
			slog.Warn("unable to drop", "database", database, "err", err)
			continue
		}
		slog.Info("database removed", "database", database)
	}

	return nil
}
