package logger

import (
	"os"
	"path/filepath"
	"runtime"
)

func (l *Logger) SetDefaultLogsPath() {
	switch runtime.GOOS {
	case "windows":
		l.Path = filepath.Join(os.Getenv("USERPROFILE"), "Downloads", "MangaDownloader", "Logs")
	default:
		l.Path = filepath.Join(os.Getenv("HOME"), "MangaDownloader", "Logs")
	}
}

func (l *Logger) ChangePath(path string) {
	if isValidPath(path) {
		l.Path = path
	} else {
		l.SetDefaultLogsPath()
	}
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

func (l *Logger) checkPath() {
	switch l.Path {
	case "", ".", " ":
		l.SetDefaultLogsPath()
	case "/":
		// "Вы точно хотите сохранять логи в корневую папку?"
	}
}
