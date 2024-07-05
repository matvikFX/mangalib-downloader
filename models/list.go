package models

type MangaListData struct {
	Meta  *Meta     `json:"meta"`
	Manga MangaList `json:"data"`
}

type Manga struct {
	RusName string `json:"rus_name"`
	EngName string `json:"eng_name"`
	Name    string `json:"name"`
	Slug    string `json:"slug"`
}
type MangaList []*Manga

type Meta struct {
	Page int `json:"page"`
	From int `json:"from"`
	To   int `json:"to"`
}

func (m *Manga) RusNameChange() {
	if m.RusName == "" {
		m.RusName = m.EngName
	}

	if m.EngName == "" {
		m.RusName = m.Name
	}
}
