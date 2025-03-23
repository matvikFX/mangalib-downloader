package api

import (
	"bufio"
	"log/slog"
	"os"
)

var authToken = ""

func readAuthToken() string {
	log := slog.With("API", "readAuthToken")

	var exists bool
	authToken, exists = os.LookupEnv("AUTH_TOKEN")
	if exists {
		return authToken
	}

	if authToken != "" {
		return authToken
	}

	file, err := os.Open("auth_token")
	if err != nil {
		log.Error("Error openning file", "Error", err)
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		authToken = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Error("Error reading string", "Error", err)
		return ""
	}

	return authToken
}
