package main

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"manga-downloader/components"
	"manga-downloader/services"
)

func main() {
	cfg := services.NewConfig()

	localtime := time.Now().Local()
	fileName := localtime.Format(time.DateOnly) + ".json"
	filePath := filepath.Join(cfg.LoggerPath, fileName)

	if err := os.MkdirAll(cfg.LoggerPath, 0o750); err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o640)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileHandler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})
	logger := slog.New(fileHandler)

	if err := Init(logger, cfg); err != nil {
		panic(err)
	}
}

func Init(logger *slog.Logger, cfg *services.Config) error {
	log := logger.With("App", "Init")

	// log.Info("Loading config")
	// cfg := services.NewConfig()
	// if err := cfg.Load(); err != nil {
	// 	log.Error("Error loading config", "Error", err)
	// 	return err
	// }
	// log.Debug("Config object", "config", cfg)

	log.Info("Loading bookmarks")
	bookmarks := services.NewBookmarks(logger)
	if err := bookmarks.Load(cfg.BookmarksPath); err != nil {
		log.Error("Error loading bookmarks", "Error", err)
		return err
	}

	log.Info("Starting application")
	app := components.NewTViewApp(logger, cfg, bookmarks)
	if err := app.Start(); err != nil {
		return err
	}
	defer app.Stop()
	log.Info("Application stopped")

	return nil
}
