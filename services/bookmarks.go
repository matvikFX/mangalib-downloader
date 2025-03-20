package services

import (
	"log"
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
	// if err := os.MkdirAll(b.Path, 0o640); err != nil {
	// 	return err
	// }

	content, err := yaml.Marshal(b.List)
	if err != nil {
		return err
	}

	if err := os.WriteFile(b.Path, content, 0o640); err != nil {
		return err
	}

	return nil
}

func (b *Bookmarks) Load(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Println("Bookmarks file does not exists. Creating...")
		if _, err := os.Create(path); err != nil {
			log.Println("Error creating bookmarks file: ", err)
			return err
		}

		return nil
	}

	if err = yaml.Unmarshal(content, &b.List); err != nil {
		return err
	}

	return nil
}
