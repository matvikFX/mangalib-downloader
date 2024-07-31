package logger

import (
	"os"
	"path/filepath"
	"runtime"
)

func DefaultLogsPath() string {
	var path string
	switch runtime.GOOS {
	case "windows":
		path = filepath.Join(os.Getenv("USERPROFILE"), "Downloads", "MangaDownloader", "Logs")
	default:
		path = filepath.Join(os.Getenv("HOME"), "MangaDownloader", "Logs")
	}
	return path
}

func isValidPath(path string) bool {
	if filepath.IsAbs(path) {
		return true
	}

	if _, err := filepath.Abs(path); err == nil {
		return true
	}

	return false
}
