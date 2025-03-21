package utils

import (
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
