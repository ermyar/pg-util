package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ermyar/pg-util/internal/app"
)

func main() {
	app, err := app.Parse()
	if err != nil {
		slog.Error("wrong input args", "err", err.Error())
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := app.Run(ctx); err != nil {
		slog.Error("finished with error", "err", err.Error())
		os.Exit(1)
	}
}
