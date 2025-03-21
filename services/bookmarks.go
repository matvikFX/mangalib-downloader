package services

import (
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

const DefaultBookmarksPath = "bookmarks"

type title struct {
	Name    string `yaml:"name"`
	URL     string `yaml:"url"`
	Chapter int    `yaml:"chapter"`
}

type Bookmarks struct {
	logger *slog.Logger

	Path string
	List map[string][]title
}

func NewBookmarks(logger *slog.Logger) *Bookmarks {
	return &Bookmarks{
		logger: logger,

		Path: DefaultPath(DefaultBookmarksPath),
		List: make(map[string][]title),
	}
}

func (b *Bookmarks) Save() error {
	log := b.logger.With("Bookmarks", "Save")

	content, err := yaml.Marshal(b.List)
	if err != nil {
		log.Error("Error marshaling bookmarks", "Error", err)
		return err
	}

	if err := os.WriteFile(b.Path, content, 0o640); err != nil {
		log.Error("Error creating bookmarks file", "Error", err)
		return err
	}

	return nil
}

func (b *Bookmarks) Load(path string) error {
	log := b.logger.With("Bookmarks", "Save")

	content, err := os.ReadFile(path)
	if err != nil {
		log.Warn("Bookmarks file does not exists. Creating...")
		if _, err := os.Create(path); err != nil {
			log.Error("Error creating bookmarks file", "Error", err)
			return err
		}

		return nil
	}

	if err = yaml.Unmarshal(content, &b.List); err != nil {
		log.Error("Error unmarshaling bookmarks", "Error", err)
		return err
	}

	return nil
}
