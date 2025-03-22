package downloader

import (
	"fmt"
	"log/slog"
	"manga-downloader/models"
	"manga-downloader/services"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (d *Downloader) ChangeConfig(path string, cbzFormat bool) error {
	log := slog.With("Downloader", "ChangePath")

	if !isValidPath(path) {
		log.Error("Setting path to default")
		d.downloadPath = services.DefaultPath("")

		err := "invalid path"
		return fmt.Errorf(err)
	}

	d.downloadPath = path
	d.cbzFormat = cbzFormat
	return nil
}

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

	if err := os.MkdirAll(chapPath, 0o644); err != nil {
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

func createFile(data []byte, pagePath string) error {
	file, err := os.Create(pagePath)
	if err != nil {
		slog.Error("Error creating file", "Error", err)
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		slog.Error("Error writing to file", "Error", err)
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

func downloadedChapters(chapters models.ChapterList) []string {
	chapList := make([]string, len(chapters))
	for idx, chapter := range chapters {
		var fullName string
		if chapter.Name == "" {
			fullName = fmt.Sprintf(
				"Том %s Глава %s",
				chapter.Volume, chapter.Number,
			)
		} else {
			fullName = fmt.Sprintf(
				"Том %s Глава %s - %s",
				chapter.Volume, chapter.Number, chapter.Name,
			)
		}
		fullName = strings.TrimSpace(fullName)
		chapList[idx] = fullName
	}

	return chapList
}
