package downloader

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func CreateCBZArchive(sourceDir string) error {
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		return err
	}

	cbzPath := sourceDir + ".cbz"
	archiveFile, err := os.Create(cbzPath)
	if err != nil {
		return err
	}
	defer archiveFile.Close()

	zipWriter := zip.NewWriter(archiveFile)
	defer zipWriter.Close()

	err = filepath.Walk(sourceDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filePath == sourceDir {
			return nil
		}

		relPath, err := filepath.Rel(sourceDir, filePath)
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
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
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
