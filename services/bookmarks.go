package services

import (
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
	Path string
	List map[string][]title
}

func NewBookmarks() *Bookmarks {
	return &Bookmarks{
		Path: DefaultPath(DefaultBookmarksPath),
		List: make(map[string][]title),
	}
}

func (b *Bookmarks) Save() error {
	// if err := os.MkdirAll(b.Path, 0o644); err != nil {
	// 	return err
	// }

	content, err := yaml.Marshal(b.List)
	if err != nil {
		return err
	}

	if err := os.WriteFile(b.Path, content, 0o644); err != nil {
		return err
	}

	return nil
}

func (b *Bookmarks) Load() error {
	content, err := os.ReadFile(b.Path)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(content, &b.List); err != nil {
		return err
	}

	return nil
}
