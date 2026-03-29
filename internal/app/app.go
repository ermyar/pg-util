package app

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/ermyar/pg-util/internal/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
)

type App struct {
	databases []string
	action    func(context.Context) error
	pool      *pgxpool.Pool
}

var (
	ErrUnknownOperation = errors.New("provided unknown operation")
	ErrUnknownLogLevel  = errors.New("provided unknown log level")
)

func (a *App) Run(ctx context.Context) error {
	defer a.pool.Close()
	return a.action(ctx)
}

func Parse() (*App, error) {
	var app App
	var cfg postgres.Config
	var operation string

	if err := godotenv.Load(); err != nil {
		slog.Debug("unable to load env", "err", err)
	}

	opts := &slog.HandlerOptions{}

	switch os.Getenv("LOGLEVEL") {
	case "DEBUG":
		opts.Level = slog.LevelDebug
	case "INFO":
		opts.Level = slog.LevelInfo
	case "WARN", "":
		opts.Level = slog.LevelWarn
	case "ERROR":
		opts.Level = slog.LevelError
	default:
		return nil, ErrUnknownLogLevel
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, opts)))

	pflag.StringSliceVarP(&app.databases, "databases", "d", nil, "list target databases names.")
	pflag.StringVarP(&operation, "operation", "o", "", "operation type, must be: backup, remove.")
	pflag.StringVarP(&cfg.Host, "host", "h", os.Getenv("PGHOST"), "Specifies the host name of the machine on which the server is running.")
	pflag.StringVarP(&cfg.Port, "port", "p", os.Getenv("PGPORT"), "Specifies the TCP port number on which the server is listening.")
	pflag.StringVarP(&cfg.User, "username", "U", os.Getenv("PGUSER"), "Connect to the PostgreSQL as the username.")
	pflag.StringVarP(&cfg.Password, "password", "W", os.Getenv("PGPASSWORD"), "Connect to the PostgreSQL with the password.")
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
	slog.Info("connected to postgres")

	app.pool = pool

	return &app, nil
}
