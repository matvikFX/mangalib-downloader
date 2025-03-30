package utils

import (
	"fmt"
	"strings"

	"mangalib-downloader/models"
)

const info = `Называние: %s
Статус: %s
Статус перевода: %s
Тип: %s
Год выпуска: %s
Количество глав: %d

Описание:
%s

Теги: 
%s

Жанры:
%s

Авторы:
%s
Переводчики:
%s`

func InfoText(manga *models.MangaInfo, teamList []string) string {
	ts := manga.GetTags()
	tags := strings.Join(ts, ", ")

	gs := manga.GetGenres()
	genres := strings.Join(gs, ", ")

	teams := strings.Join(teamList, "\n")
	if len(teams) == 0 {
		teams = strings.Join(manga.Branches.GetTeams(), "\n")
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
		tags,
		genres,
		authors, teams)

	return infoText
}
