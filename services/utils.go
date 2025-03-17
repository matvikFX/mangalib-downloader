package services

import (
	"errors"
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

func ChangePath(path string) (string, error) {
	var newPath string
	if isPathValid(path) {
		newPath = path
	} else {
		return "", errors.New("invalid path")
	}
	return newPath, nil
}

func isPathValid(path string) bool {
	if filepath.IsAbs(path) {
		return true
	}

	if _, err := filepath.Abs(path); err != nil {
		return true
	}

	return false
}
