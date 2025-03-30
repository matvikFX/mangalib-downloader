package api

import (
	"bufio"
	"os"
)

var authToken = ""

func readAuthToken() string {
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
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		authToken = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return ""
	}

	return authToken
}
