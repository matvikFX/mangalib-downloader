package api

import (
	"bufio"
	"log/slog"
	"os"
)

var authToken = ""

func readAuthToken() string {
	if authToken != "" {
		return authToken
	}

	file, err := os.Open("auth_token")
	if err != nil {
		slog.Error("Error openning file", "Error", err)
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		authToken = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		slog.Error("Error reading string", "Error", err)
		return ""
	}

	return authToken
}
