package main

import (
	"log/slog"
	"os"

	"github.com/samuel-poirier/toggl-to-dynamics-timesheet-converter/cmd/internal/app"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	app := app.New(logger)

	allArgs := os.Args

	config, err := app.ParseArgs(allArgs)

	if err != nil {
		logger.Error("failed to parse arguments", slog.Any("error", err))
		return
	}

	togglExportLines, err := app.LoadTogglCsvExportLines(config.InputFile)

	if err != nil {
		logger.Error("failed to parse arguments", slog.Any("error", err))
		return
	}

	logger.Info("loaded entries", slog.Any("count", len(togglExportLines)))

	err = app.Export(config, togglExportLines)

	if err != nil {
		logger.Error("failed to export", slog.Any("error", err))
		return
	}

	logger.Info("export completed successfully", slog.Any("output", config.OutputFile))

}
