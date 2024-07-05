package models

type Page struct {
	// Название картинки
	Image string `json:"image"`

	// Номер страницы
	Slug int `json:"slug"`

	Height int `json:"height"`
	Width  int `json:"width"`

	// Ссылка на картинку
	//
	// chapterID + image
	URL string `json:"url"`
}
type PageList []*Page
