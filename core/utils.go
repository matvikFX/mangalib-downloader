package core

import (
	"os"
	"path/filepath"
	"runtime"
)

func SetDefaultDownloadPath() string {
	var defaultDownloadPath string
	switch runtime.GOOS {
	case "windows":
		defaultDownloadPath = filepath.Join(os.Getenv("USERPROFILE"), "Downloads", "MangaDownloader")
	default:
		defaultDownloadPath = filepath.Join(os.Getenv("HOME"), "MangaDownloader")
	}
	return defaultDownloadPath
}

func SetDefaultLogsPath(downloadPath string) string {
	if downloadPath != "" {
		defaultPath := SetDefaultDownloadPath()
		return filepath.Join(defaultPath, "Logs")
	}
	return filepath.Join(downloadPath, "Logs")
}
