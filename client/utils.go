package client

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func (c *MangaLibClient) GetBranchTeams(ctx context.Context, branchID int) string {
	branchTeams := make(map[int]string)
	if c.Branch != 0 {
		branches, err := c.GetMangaBranches(ctx, branchID)
		if err != nil {
			c.Logger.WriteLog(err.Error())
		}

		branchTeams = branches.BranchTeams()
	}

	return branchTeams[c.Branch]
}

// teams необязательно указывать
func (c *MangaLibClient) CreateChapterPath(teams, mangaName string, volume, number, chapName string) string {
	mangaName = removeChars(mangaName)
	teams = removeChars(teams)
	chapName = removeChars(chapName)

	var chapDir string
	if chapName == "" {
		chapDir = fmt.Sprintf("Том %s Глава %s", volume, number)
	} else {
		chapDir = fmt.Sprintf("Том %s Глава %s - %s", volume, number, chapName)
	}
	chapDir = strings.TrimSpace(chapDir)

	var chapterPath string
	if teams == "" {
		chapterPath = filepath.Join(c.DownloadPath, mangaName, chapDir)
	} else {
		chapterPath = filepath.Join(c.DownloadPath, mangaName, teams, chapDir)
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
	rusName = removeChars(rusName)

	chapPath := c.CreateChapterPath(branchTeams, rusName, volume, number, name)

	if err := os.MkdirAll(chapPath, 0o755); err != nil {
		return err
	}

	return nil
}

func (c *MangaLibClient) ChangePath(path string) {
	if isValidPath(path) {
		c.DownloadPath = path
	} else {
		DefaultDownloadPath()
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

func DefaultDownloadPath() string {
	var path string
	switch runtime.GOOS {
	case "windows":
		path = filepath.Join(os.Getenv("USERPROFILE"), "Downloads", "MangaDownloader")
	default:
		path = filepath.Join(os.Getenv("HOME"), "MangaDownloader")
	}
	return path
}

func createFile(data []byte, pagePath string) error {
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

func createPageName(pageSlug int, pageImg string) string {
	return strconv.Itoa(pageSlug) + filepath.Ext(pageImg)
}

func createPagePath(chapPath, pageName string) string {
	return filepath.Join(chapPath, pageName)
}

func removeChars(text string) string {
	charsToReplace := []string{"<", ">", ":", "/", "|", "?", "*", "\"", "\\", "."}
	for _, char := range charsToReplace {
		text = strings.ReplaceAll(text, char, "")
	}

	return text
}
