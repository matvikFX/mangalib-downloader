package downloader

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func (c *MangaLibDownloader) ChangePath(path string) string {
	if !isValidPath(path) {
		c.DownloadPath = DefaultDownloadPath()

		c.Logger.Write("Error: invalid path. Setting path to default")
		return "Invalid path. Setting path to default"
	}

	c.DownloadPath = path
	return ""
}

// Находится в API
// func (c *MangaLibDownloader) GetBranchTeams(ctx context.Context, branchID int) string {
// 	branchTeams := make(map[int]string)
// 	if branchID != 0 {
// 		branches, err := api.GetMangaBranches(ctx, branchID)
// 		if err != nil {
// 			c.Logger.Write(err.Error())
// 		}
//
// 		branchTeams = branches.BranchTeams()
// 	}
//
// 	return branchTeams[branchID]
// }

// teams необязательно указывать
func CreateChapterPath(
	downloadPath string,
	teams, mangaName string, volume, number, chapName string,
) string {
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
		chapterPath = filepath.Join(downloadPath, mangaName, chapDir)
	} else {
		chapterPath = filepath.Join(downloadPath, mangaName, teams, chapDir)
	}

	return chapterPath
}

func CheckExistence(filePath string) bool {
	var exists bool

	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		exists = true
	}

	return exists
}

func createFolder(downloadPath, rusName, branchTeams, volume, number, name string) error {
	rusName = removeChars(rusName)

	chapPath := CreateChapterPath(downloadPath, branchTeams, rusName, volume, number, name)

	if err := os.MkdirAll(chapPath, 0o755); err != nil {
		return err
	}

	return nil
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

func createPageURL(image string) string {
	// Download URLs
	const (
		FirstURL      = "https://img2.mixlib.me"
		SecondURL     = "https://img4.imgslib.link"   // Работает
		CompressedURL = "https://img33.imgslib.link/" // Работает
		DownloadURL   = "https://img4.imgslib.org"
	)

	return CompressedURL + image
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
