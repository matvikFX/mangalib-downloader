package core

import (
	"fmt"
	"log"
	"os"
	"strings"

	"mangalib-downloader/client"
	"mangalib-downloader/logger"
)

const cfgFile = "config"

var (
	dLen = len("DOWNLOAD_PATH=")
	lLen = len("LOGS_PATH=")
)

func (a *MangaApp) DefaultConfig() {
	a.Client.DownloadPath = client.DefaultDownloadPath()
	a.Client.Logger.Path = logger.DefaultLogsPath()

	a.SaveConfig()
}

func (a *MangaApp) LoadConfig() {
	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		a.DefaultConfig()
		return
	}

	file, err := os.ReadFile(cfgFile)
	if err != nil {
		log.Println("can not open config file")
	}

	paths := strings.Split(string(file), "\n")
	dPath := paths[0][dLen:]
	lPath := paths[1][lLen:]

	a.Client.DownloadPath = dPath
	a.Client.Logger.Path = lPath
}

func (a *MangaApp) SaveConfig() {
	text := fmt.Sprintf("DOWNLOAD_PATH=%s\nLOGS_PATH=%s",
		a.Client.DownloadPath, a.Client.Logger.Path)

	if err := os.WriteFile(cfgFile, []byte(text), os.ModePerm); err != nil {
		return
	}
}
