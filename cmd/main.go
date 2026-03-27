package main

import (
	"log"
	"log/slog"

	app "github.com/ermyar/pg-util/internal/app"
)

func main() {
	app, err := app.Parse()
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := app.Run(); err != nil {
		slog.Error(err.Error())
	}
}
