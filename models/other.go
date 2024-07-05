package models

type Author struct {
	ID      int    `json:"id"`
	Slug    string `json:"slug"`
	Name    string `json:"name"`
	RusName string `json:"rus_name"`
}
type authorList []*Author

type mangaChapters struct {
	Uploaded int `json:"uploaded"`

	// Вообще не понятно что
	Total int `json:"total"`
}

type mangaType struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}
type mangaRating struct {
	Avg   string `json:"average"`
	Votes int    `json:"votes"`
}

type status struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

type mangaViews struct {
	Total    int    `json:"total"`
	Short    string `json:"short"`
	Fromated string `json:"formated"`
}

type mangaGenreTag struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Adult bool   `json:"adult"`
}
type genreTagList []*mangaGenreTag
