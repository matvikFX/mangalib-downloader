package main

import (
	"log"
	"log/slog"
	"manga-downloader/components"
	"manga-downloader/downloader"
	"manga-downloader/services"
	"os"
	"path/filepath"
	"time"
)

func main() {
	cfg, err := services.Load()
	if err != nil {
		log.Fatal(err)
	}

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
	slog.SetDefault(logger)

	if err := Init(logger, cfg); err != nil {
		log.Fatal(err)
	}
}

func Init(logger *slog.Logger, cfg *services.Config) error {
	log := slog.With("App", "Init")

	log.Info("Loading downloader")
	downloader := downloader.New(cfg.DownloadPath, cfg.CbzFormat)

	log.Info("Loading bookmarks")
	bookmarks := services.NewBookmarks(logger)
	if err := bookmarks.Load(cfg.BookmarksPath); err != nil {
		log.Error("Error loading bookmarks", "Error", err)
		return err
	}

	log.Info("Starting application")
	app := components.NewTViewApp(cfg, bookmarks, downloader)
	if err := app.Start(); err != nil {
		return err
	}
	defer app.Stop()
	log.Info("Application stopped")

	return nil
}
