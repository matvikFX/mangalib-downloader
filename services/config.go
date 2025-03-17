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
}

func NewConfig() *Config {
	config := &Config{
		DownloadPath:  DefaultPath(""),
		LoggerPath:    DefaultPath(DefaultLoggerPath),
		BookmarksPath: DefaultPath(DefaultBookmarksPath),
	}

	config.Save()
	return config
}

func (c *Config) Default() {
	c = NewConfig()
}

func (c *Config) Save() {
	jsonConf, err := yaml.Marshal(c)
	if err != nil {
		log.Println("can't marshal config")
		return
	}

	if err := os.WriteFile(cfgFile, jsonConf, 0o644); err != nil {
		log.Println("can't save config file")
		return
	}
}

func (c *Config) Load() error {
	content, err := os.ReadFile(cfgFile)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(content, c); err != nil {
		return err
	}

	return nil
}
