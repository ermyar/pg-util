package app

import (
	"errors"
	"log/slog"
	"os"

	"github.com/spf13/pflag"
)

type App struct {
	database []string
	action   func() error
}

var (
	ErrUnknownOperation = errors.New("provided unknown operation")
)

func (a *App) Run() error {
	return a.action()
}

func Parse() (*App, error) {
	var app App
	var operation string

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	})))

	pflag.StringSliceVar(&app.database, "databases", nil, "list target databases names")
	pflag.StringVar(&operation, "operation", "", "operation type, must be: backup, remove")
	pflag.Parse()

	switch operation {
	case "backup":
		app.action = app.backup
	case "remove":
		app.action = app.remove
	default:
		return nil, ErrUnknownOperation
	}

	return &app, nil
}
