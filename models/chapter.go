package models

type ChapterData struct {
	Data *Chapter `json:"data"`
}

type Chapter struct {
	ID    int `json:"id"`
	Index int `json:"index"`

	Volume string `json:"volume"`
	Number string `json:"number"`

	Name      string   `json:"name"`
	Pages     PageList `json:"pages"`
	CreatedAt string   `json:"created_at"`

	BranchID    int        `json:"branches_id"`
	BranchCount int        `json:"branches_count"`
	Branches    BranchList `json:"branches"`
}

type ChapterList []*Chapter

type ChaptersData struct {
	Data ChapterList `json:"data"`
}
