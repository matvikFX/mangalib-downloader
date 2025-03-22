package services

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	cfgFile              = "config.yaml"
	DefaultLoggerPath    = "logs"
	DefaultBookmarksPath = "bookmarks.yaml"
)

type Config struct {
	DownloadPath  string `yaml:"download_path"`
	LoggerPath    string `yaml:"logger_path"`
	BookmarksPath string `yaml:"bookmarks_path"`
	CbzFormat     bool   `yaml:"cbz_format"`
}

func Load() (*Config, error) {
	config := &Config{}
	if err := config.Load(); err != nil {
		log.Println("Config file does not exists. Creating new file")
		config, err = DefaultConfig()
		if err != nil {
			log.Fatal("Can't create config file")
			return nil, err
		}
	}
	return config, nil
}

// Creates config with default paths:
// DownloadPath: $HOME/MangaDownloader,
// LoggerPath: $HOME/MangaDownloader/logs,
// BookmarksPath: $HOME/MangaDownloader/cbzFormatokmarks.yaml
func DefaultConfig() (*Config, error) {
	config := &Config{
		DownloadPath:  DefaultPath(""),
		LoggerPath:    DefaultPath(DefaultLoggerPath),
		BookmarksPath: DefaultPath(DefaultBookmarksPath),
		CbzFormat:     true,
	}

	file, err := os.Create(cfgFile)
	if err != nil {
		log.Println("Error creating file")
		return config, err
	}
	defer file.Close()

	jsonConf, err := yaml.Marshal(config)
	if err != nil {
		log.Println("Error marshaling config")
		return config, err
	}

	_, err = file.Write(jsonConf)
	if err != nil {
		log.Println("Error writing config to file")
		return config, err
	}

	return config, nil
}

func (c *Config) Save() error {
	jsonConf, err := yaml.Marshal(c)
	if err != nil {
		log.Println("Error marshaling config")
		return err
	}

	if err := os.WriteFile(cfgFile, jsonConf, 0o644); err != nil {
		log.Println("Error writing config to file")
		return err
	}

	return nil
}

func (c *Config) Load() error {
	content, err := os.ReadFile(cfgFile)
	if err != nil {
		log.Println("Config file does not exists")
		return err
	}

	if err := yaml.Unmarshal(content, c); err != nil {
		log.Println("Can't unmarshal config")
		return err
	}

	return nil
}

func (c *Config) ChangeDownloadPath(newDownloadPath string) string {
	if !IsPathValid(newDownloadPath) {
		c.DownloadPath = DefaultPath("")

		err := "Invalid path. Setting download path to default"
		log.Println(err)
		return err
	}

	c.DownloadPath = newDownloadPath
	return ""
}

func (c *Config) ChangeLogPath(newLogPath string) string {
	if !IsPathValid(newLogPath) {
		c.LoggerPath = DefaultPath(DefaultLoggerPath)

		err := "Invalid path. Setting logger path to default"
		log.Println(err)
		return err
	}

	c.LoggerPath = newLogPath
	return ""
}

func (c *Config) ChangeBookmarksPath(newBookmarksPath string) string {
	if !IsPathValid(newBookmarksPath) {
		c.BookmarksPath = DefaultPath(DefaultBookmarksPath)

		err := "Invalid path. Setting bookmarks path to default"
		log.Println(err)
		return err
	}

	c.BookmarksPath = newBookmarksPath
	return ""
}
