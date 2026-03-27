package app

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/ermyar/pg-util/internal/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/pflag"
)

type App struct {
	database []string
	action   func() error
	pool     *pgxpool.Pool
}

var (
	ErrUnknownOperation = errors.New("provided unknown operation")
)

func (a *App) Run() error {
	defer a.pool.Close()
	return a.action()
}

func Parse() (*App, error) {
	var app App
	var cfg postgres.PgConfig
	var operation string

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr,
		&slog.HandlerOptions{
			Level: slog.LevelWarn,
		},
	)))

	pflag.StringSliceVarP(&app.database, "databases", "d", nil, "list target databases names")
	pflag.StringVarP(&operation, "operation", "o", "", "operation type, must be: backup, remove")
	pflag.StringVarP(&cfg.Host, "host", "h", os.Getenv("PG_HOST"), "Specifies the host name of the machine on which the server is running.")
	pflag.StringVarP(&cfg.Port, "port", "p", os.Getenv("PG_PORT"), "Specifies the TCP port number on which the server is listening.")
	pflag.StringVarP(&cfg.User, "username", "U", os.Getenv("PG_USER"), "Connect to the PostgreSQL as the username.")
	pflag.StringVarP(&cfg.Password, "password", "W", os.Getenv("PG_PASS"), "Connect to the PostgreSQL with the password.")
	pflag.Parse()

	switch operation {
	case "backup":
		app.action = app.backup
	case "remove":
		app.action = app.remove
	default:
		return nil, ErrUnknownOperation
	}

	pool, err := cfg.Connect(context.Background())
	if err != nil {
		return nil, err
	}

	app.pool = pool

	return &app, nil
}
