package models

import (
	"strings"
)

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
	ID        int         `json:"id"`
	CreatedAt string      `json:"created_at"`
	Teams     TeamList    `json:"teams"`
	Details   teamDetails `json:"details"`

	// Используется при получении глав
	BranchID int `json:"branch_id"`
}
type BranchList []*Branch

type BranchesData struct {
	Data BranchList `json:"data"`
}

func (b BranchList) GetTeams() []string {
	teams := make([]string, len(b))
	for i, branch := range b {
		var builder strings.Builder
		for j, team := range branch.Teams {
			if j > 0 {
				builder.WriteString(", ")
			}
			builder.WriteString(team.Name)
		}
		teams[i] = builder.String()
	}
	return teams
}

func (b BranchList) BranchTeams() map[int]string {
	teams := make(map[int]string)
	for _, branch := range b {
		var builder strings.Builder
		for idx, team := range branch.Teams {
			if idx > 0 {
				builder.WriteString(", ")
			}
			builder.WriteString(team.Name)
		}
		teams[branch.ID] = builder.String()
	}
	return teams
}

func (b BranchList) BranchTeamList() map[int][]string {
	teams := make(map[int][]string)
	for _, branch := range b {
		teams[branch.ID] = make([]string, len(branch.Teams))
		for i, team := range branch.Teams {
			teams[branch.ID][i] = team.Name
		}
	}
	return teams
}

func (b BranchList) TeamsBranch() map[string]int {
	teams := make(map[string]int)
	for _, branch := range b {
		var builder strings.Builder
		for idx, team := range branch.Teams {
			if idx > 0 {
				builder.WriteString(", ")
			}
			builder.WriteString(team.Name)
		}
		teams[builder.String()] = branch.ID
	}
	return teams
}
