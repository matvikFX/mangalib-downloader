package services

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"
)

const DefaultLoggerPath = "logs"

type Logger struct {
	Path string
}

func NewLogger() *Logger {
	return &Logger{
		Path: DefaultPath(DefaultLoggerPath),
	}
}

func (l *Logger) Write(logStruct any) {
	localtime := time.Now().Local()

	if err := os.MkdirAll(l.Path, 0o750); err != nil {
		log.Println("Error creating log folder: ", err)
		return
	}

	fileName := localtime.Format(time.DateOnly) + ".json"
	filePath := filepath.Join(l.Path, fileName)

	jsonLog, err := json.Marshal(logStruct)
	if err != nil {
		log.Println("Error to marshal struct: ", err)
		return
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o640)
	if err != nil {
		log.Println("Error to open log file: ", err)
		return
	}
	defer file.Close()

	_, err = file.Write(jsonLog)
	if err != nil {
		log.Println("Error writing into file: ", err)
		return
	}
}

func (l *Logger) ChangePath(path string) string {
	if !IsPathValid(path) {
		l.Path = DefaultPath(DefaultLoggerPath)

		log.Println("Error: invalid path. Setting logger path to default")
		return "Invalid path. Setting logger path to default"
	} else {
		l.Path = path
	}
	return ""
}
