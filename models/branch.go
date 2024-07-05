package models

type Team struct {
	ID      int         `json:"id"`
	Slug    string      `json:"slug"`
	Name    string      `json:"name"`
	Details teamDetails `json:"details"`
}
type TeamList []*Team

type teamDetails struct {
	BranchID int  `json:"branch_id"`
	IsActive bool `json:"is_active"`
}

type Branch struct {
	ID        int      `json:"id"`
	BranchID  int      `json:"branch_id"`
	CreatedAt string   `json:"created_at"`
	Teams     TeamList `json:"teams"`
}
type BranchList []*Branch

type BranchesData struct {
	Data BranchList `json:"data"`
}
