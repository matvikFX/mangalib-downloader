package downloader

import (
	"archive/zip"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

func CreateCBZArchive(sourceDir string) error {
	log := slog.With("Downloader", "CreateCBZArchive")

	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		log.Error(fmt.Sprintf("Dir %s doesn't exists", sourceDir), "Error", err)
		return err
	}

	cbzPath := sourceDir + ".cbz"
	archiveFile, err := os.Create(cbzPath)
	if err != nil {
		log.Error("Error creating archive", "Error", err)
		return err
	}
	defer archiveFile.Close()

	zipWriter := zip.NewWriter(archiveFile)
	defer zipWriter.Close()

	err = filepath.Walk(sourceDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			log.Error("Error occured", "Error", err)
			return err
		}

		if filePath == sourceDir {
			log.Error("Trying to archive an archive", "Error", err)
			return nil
		}

		relPath, err := filepath.Rel(sourceDir, filePath)
		if err != nil {
			log.Error("Error getting relative path", "Error", err)
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			log.Error("Error creating header", "Error", err)
			return err
		}

		header.Name = filepath.ToSlash(relPath) // Для кросс-платформенности
		header.Method = zip.Deflate

		if info.IsDir() {
			header.Name += "/"
			header.Method = zip.Store
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			log.Error("Error creating writer with header",
				"header", header, "Error", err)
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(filePath)
			if err != nil {
				log.Error("Error opening archive", "Error", err)
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			if err != nil {
				log.Error("Error writing data into archive", "Error", err)
				return err
			}
		}

		return nil
	})
	if err != nil {
		log.Error("Error creating archive", "Error", err)
		return err
	}

	return nil
}
