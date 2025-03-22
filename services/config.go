package services

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

const cfgFile = "config.yaml"

type Config struct {
	DownloadPath  string `yaml:"download_path"`
	LoggerPath    string `yaml:"logger_path"`
	BookmarksPath string `yaml:"bookmarks_path"`
	CbzFormat     bool   `yaml:"cbz_format"`
}

func NewConfig() *Config {
	config := &Config{}
	if err := config.Load(); err != nil {
		config = &Config{
			DownloadPath:  DefaultPath(""),
			LoggerPath:    DefaultPath(DefaultLoggerPath),
			BookmarksPath: DefaultPath(DefaultBookmarksPath),
			CbzFormat:     true,
		}
		config.Save()
	}
	return config
}

func (c *Config) Default() {
	c = NewConfig()
}

func (c *Config) Save() error {
	// log := c.logger.With("Config", "Save")

	jsonConf, err := yaml.Marshal(c)
	if err != nil {
		// log.Fatal(err)
		// log.Error("Can't marshal config", "Error", err)
		return err
	}

	if err := os.WriteFile(cfgFile, jsonConf, 0o644); err != nil {
		// log.Fatal(err)
		// log.Error("Can't save config file", "Error", err)
		return err
	}

	return nil
}

func (c *Config) Load() error {
	// log := c.logger.With("Config", "Load")

	content, err := os.ReadFile(cfgFile)
	if err != nil {
		// log.Error("Config file does not exists", "Error", err)
		// return err
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(content, c); err != nil {
		// log.Error("Can't unmarshal config", "Error", err)
		// return err
		log.Fatal(err)
	}

	return nil
}

func (c *Config) ChangeDownloadPath(newPath string) string {
	// log := c.logger.With("Config", "ChangeDownloadPath")
	//
	if !IsPathValid(newPath) {
		c.DownloadPath = DefaultPath("")
		err := "Invalid path. Setting download path to default"
		// log.Error(err)
		return err
	}

	c.DownloadPath = newPath
	// log.Warn("Download path changed", "DownloadPath", c.DownloadPath)

	return ""
}

func (c *Config) ChangeLogPath(newPath string) string {
	// log := c.logger.With("Config", "ChangeLogPath")

	if !IsPathValid(newPath) {
		c.LoggerPath = DefaultPath(DefaultLoggerPath)
		err := "Invalid path. Setting logger path to default"
		// log.Error(err)
		return err
	}

	c.LoggerPath = newPath
	// log.Warn("Logger path changed", "LoggerPath", c.LoggerPath)

	return ""
}

func (c *Config) ChangeBookmarkPath(newPath string) string {
	// log := c.logger.With("Config", "ChangeBookmarkPath")

	if !IsPathValid(newPath) {
		c.BookmarksPath = DefaultPath(DefaultBookmarksPath)
		err := "Invalid path. Setting bookmarks path to default"
		// log.Error(err)
		return err
	}

	c.BookmarksPath = newPath
	// log.Warn("Bookmarks path changed", "BookmarksPath", c.BookmarksPath)

	return ""
}
