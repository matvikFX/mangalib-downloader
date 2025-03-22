package models

// Дополнительные поля в ссылке:
// fields[]=summary&fields[]=views&fields[]=genres&fields[]=tags&fields[]=releaseDate&fields[]=chap_count

type MangaInfoData struct {
	Data *MangaInfo `json:"data"`
}

type MangaInfo struct {
	Manga

	Translate status `json:"scanlateStatus"`
	Status    status `json:"status"`

	Description string `json:"summary"`
	ReleaseDate string `json:"releaseDate"`

	Genres genreTagList `json:"genres"`
	Tags   genreTagList `json:"tags"`

	Rating mangaRating `json:"rating"`
	Type   mangaType   `json:"type"`
	Views  mangaViews  `json:"views"`

	ChapterCount mangaChapters `json:"items_count"`
	Authors      authorList    `json:"authors"`
	Teams        TeamList      `json:"teams"`

	Branches BranchList
}
type MangaInfoList []*MangaInfo

func (m *MangaInfo) GetTags() []string {
	tags := []string{}
	for _, tag := range m.Tags {
		tags = append(tags, tag.Name)
	}

	return tags
}

func (m *MangaInfo) GetGenres() []string {
	genres := []string{}
	for _, genre := range m.Genres {
		genres = append(genres, genre.Name)
	}

	return genres
}
