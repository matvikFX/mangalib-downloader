package utils

import (
	"fmt"

	"mangalib-downloader/models"
)

const info = `Называние: %s
Статус: %s
Статус перевода: %s
Тип: %s
Выпуск: %s
Количество глав: %d

Описание:
%s

Авторы:
%s
Переводчики:
%s
`

func InfoText(manga *models.MangaInfo, teamList []string) string {
	var teams string
	if len(teamList) == 0 {
		for _, team := range manga.Teams {
			teams += fmt.Sprintln(team.Name)
		}
	} else {
		for _, team := range teamList {
			teams += fmt.Sprintln(team)
		}
	}

	var authors string
	for _, author := range manga.Authors {
		authors += fmt.Sprintln(author.Name)
	}

	infoText := fmt.Sprintf(info,
		manga.RusName,
		manga.Status.Label,
		manga.Translate.Label,
		manga.Type.Label,
		manga.ReleaseDate,
		manga.ChapterCount.Uploaded,
		manga.Description,
		authors, teams)

	return infoText
}
