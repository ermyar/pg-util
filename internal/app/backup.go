package app

import (
	"context"
	"log/slog"
	"sync"

	"github.com/ermyar/pg-util/internal/postgres"
)

func (a *App) backup(ctx context.Context) error {
	slog.Info("backup started")
	wg := sync.WaitGroup{}

	for _, database := range a.findDatabases(ctx) {
		wg.Add(1)
		go func(database string) {
			defer wg.Done()
			postgres.Backup(ctx, a.pool, database)
		}(database)
	}
	wg.Wait()
	slog.Info("backup finished")
	return nil
}
