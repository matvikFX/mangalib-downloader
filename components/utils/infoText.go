package utils

import (
	"fmt"

	"mangalib-downlaoder/models"
)

const info = `Называние: %s
Статус: %s
Статус перевода: %s
Тип: %s
Выпуск: %s
Количество глав: %d

Описание:
%s

Переводчики:
%s
Авторы:
%s`

func InfoText(manga *models.MangaInfo) string {
	var teams string
	for _, team := range manga.Teams {
		teams += fmt.Sprintln(team.Name)
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
		teams, authors)

	return infoText
}
