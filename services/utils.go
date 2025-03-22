package services

import (
	"os"
	"path/filepath"
	"runtime"
)

func DefaultPath(service string) string {
	var path string
	switch runtime.GOOS {
	case "windows":
		path = filepath.Join(os.Getenv("USERPROFILE"), "Downloads", "MangaDownloader", service)
	default:
		path = filepath.Join(os.Getenv("HOME"), "MangaDownloader", service)
	}

	return path
}

func IsPathValid(path string) bool {
	if filepath.IsAbs(path) {
		return true
	}

	if _, err := filepath.Abs(path); err != nil {
		return true
	}

	return false
}
