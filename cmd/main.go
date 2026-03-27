package main

import (
	"log/slog"

	app "github.com/ermyar/pg-util/internal/app"
)

func main() {
	app, err := app.Parse()
	if err != nil {
		slog.Error("wrong input args", "err", err.Error())
		return
	}

	if err := app.Run(); err != nil {
		slog.Error("finished with error", "err", err.Error())
	}
}
