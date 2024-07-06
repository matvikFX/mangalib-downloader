package client

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// teams необязательно указывать
func (c *MangaLibClient) CreateChapterPath(teams, mangaName string, volume, number, chapName string) string {
	downloadsPath := filepath.Join(os.Getenv("USERPROFILE"), "Downloads", "MangaDownloader")

	teams = c.removeChars(teams)
	chapName = c.removeChars(chapName)

	var chapDir string
	if chapName == "" {
		chapDir = fmt.Sprintf("Том %s Глава %s", volume, number)
	} else {
		chapDir = fmt.Sprintf("Том %s Глава %s - %s", volume, number, chapName)
	}
	chapDir = strings.TrimSpace(chapDir)

	var chapterPath string
	if teams == "" {
		chapterPath = filepath.Join(downloadsPath, mangaName, chapDir)
	} else {
		chapterPath = filepath.Join(downloadsPath, mangaName, teams, chapDir)
	}

	return chapterPath
}

func (c *MangaLibClient) CheckExistence(filePath string) bool {
	var exists bool

	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		exists = true
	}

	return exists
}

func (c *MangaLibClient) createFolder(rusName, branchTeams, volume, number, name string) error {
	rusName = c.removeChars(rusName)

	chapPath := c.CreateChapterPath(branchTeams, rusName, volume, number, name)

	if err := os.MkdirAll(chapPath, os.ModeDir); err != nil {
		return err
	}

	return nil
}

func (c *MangaLibClient) createFile(data []byte, pagePath string) error {
	file, err := os.Create(pagePath)
	if err != nil {
		fmt.Println("Error creating file")
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		fmt.Println("Error writing to file")
		return err
	}

	return nil
}

func (c *MangaLibClient) createPageName(pageSlug int, pageImg string) string {
	return strconv.Itoa(pageSlug) + filepath.Ext(pageImg)
}

func (c *MangaLibClient) createPagePath(chapPath, pageName string) string {
	return filepath.Join(chapPath, pageName)
}

func (c *MangaLibClient) removeChars(text string) string {
	charsToReplace := []string{"<", ">", ":", "/", "|", "?", "*", "\"", "\\", "."}
	for _, char := range charsToReplace {
		text = strings.ReplaceAll(text, char, "")
	}

	return text
}
