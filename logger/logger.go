package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Logger struct {
	Path string
}

func NewLogger(path string) *Logger {
	logger := &Logger{
		Path: path,
	}

	if path != "" {
		if !isValidPath(path) {
			log.Println("Введенный текст не является абсолютным путем")
			logger.SetDefaultLogsPath()
		}
	}

	if path == "" {
		logger.SetDefaultLogsPath()
	}

	return logger
}

func (l *Logger) WriteLog(text string) {
	localTime := time.Now().Local()

	err := os.MkdirAll(l.Path, 0o755)
	if err != nil {
		log.Println("Error creating log folder: ", err)
		return
	}

	fileName := localTime.Format(time.DateOnly) + ".log"
	filePath := filepath.Join(l.Path, fileName)

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Println("Error opening file: ", err)
		return
	}
	defer file.Close()

	formated := fmt.Sprintf("%s %s\n", localTime.Format(time.TimeOnly), text)
	_, err = file.WriteString(formated)
	if err != nil {
		log.Println("Error writing into file: ", err)
		return
	}
}

func (l *Logger) WriteJSON(mangaStruct any) {
	localTime := time.Now().Local()

	err := os.MkdirAll(l.Path, 0o755)
	if err != nil {
		log.Println("Error creating log folder: ", err)
		return
	}

	fileName := localTime.Format(time.DateOnly) + ".json"
	filePath := filepath.Join(l.Path, fileName)

	jsonData, err := json.Marshal(mangaStruct)
	if err != nil {
		log.Println("Error to marshal struct: ", err)
		return
	}

	if err := os.WriteFile(filePath, jsonData, os.ModePerm); err != nil {
		log.Println("Error writing into file: ", err)
		return
	}
}
