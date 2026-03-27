package app

import "log/slog"

func (a *App) backup() error {
	slog.Info("backup started")
	return nil
}
