package app

import (
	"context"
	"log/slog"
	"sync"

	"github.com/ermyar/pg-util/internal/postgres"
)

func (a *App) remove(ctx context.Context) error {
	slog.Info("remove started")
	wg := sync.WaitGroup{}

	for _, database := range a.findDatabases(ctx) {
		wg.Add(1)
		go func(database string) {
			defer wg.Done()
			postgres.Drop(ctx, a.pool, database)
		}(database)
	}

	wg.Wait()
	slog.Info("remove finished")
	return nil
}
