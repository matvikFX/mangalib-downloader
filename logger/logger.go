package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Logger struct{}

var logPath = filepath.Join(os.Getenv("HOME"), "MangaDownloader", "Logs")

func (l *Logger) WriteLog(text string) {
	localTime := time.Now().Local()

	err := os.MkdirAll(logPath, 0o755)
	if err != nil {
		fmt.Println("Error creating log folder: ", err)
	}

	fileName := localTime.Format(time.DateOnly) + ".log"
	filePath := filepath.Join(logPath, fileName)

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Println("Error opening file: ", err)
	}
	defer file.Close()

	formated := fmt.Sprintf("%s %s\n", localTime.Format(time.TimeOnly), text)
	_, err = file.WriteString(formated)
	if err != nil {
		fmt.Println("Error writing into file: ", err)
	}
}
