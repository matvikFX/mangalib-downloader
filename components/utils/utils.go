package utils

import (
	"path/filepath"
	"strconv"
	"strings"
)

func ParseURL(url string) (int, string) {
	s := strings.Split(url, "/")
	s = strings.Split(s[len(s)-1], "--")

	id, _ := strconv.Atoi(s[0])
	slug := strings.Split(s[1], "?")[0]

	return id, slug
}

func GetMatches(currentText string) (entries []string) {
	const hintsNum = 10

	if len(currentText) == 0 {
		return nil
	}

	matchesWithPrefix, err := filepath.Glob(currentText + "*")
	if err != nil {
		return nil
	}

	if len(matchesWithPrefix) == 1 {
		if currentText == matchesWithPrefix[0] {
			return nil
		}
	}

	var matches []string
	for _, match := range matchesWithPrefix {
		dirs := strings.Split(match, "/")
		if strings.HasPrefix(dirs[len(dirs)-1], ".") {
			continue
		}
		matches = append(matches, match)
	}

	if len(matches) > hintsNum {
		matches = matches[:hintsNum]
	}

	return matches
}
