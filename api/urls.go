package api

import (
	"net/url"
	"strconv"
)

/*
1 - Manga
2 - Yaoi/Yuri (для скачивания нужен аккаунт)
3 - Ranobe
4 - Hentai (для скачивания нужен аккаунт)
5 - Anime
*/

const (
	// Оснвная ссылка на API MangaLib
	MangaLibURL = "https://api.lib.social/api/manga/"
	TeamURL     = "https://api.lib.social/api/teams/"
	BranchesURL = "https://api.lib.social/api/branches/"

	FirstURL      = "https://img2.mixlib.me"
	SecondURL     = "https://img4.imgslib.link"   // Работает
	CompressedURL = "https://img33.imgslib.link/" // Работает
	DownloadURL   = "https://img4.imgslib.org"
	// Чтобы получить список определенной команды,
	// надо к поиску добавить
	// targer_id="TeamID"&targer_model=team
)

func createSearchURL(name string, page int) string {
	queryParams := url.Values{}
	queryParams.Add("fields[]", "rate_avg")
	queryParams.Add("fields[]", "rate")
	queryParams.Add("fields[]", "releaseDate")
	queryParams.Add("site_id[]", "1")
	queryParams.Add("q", name)
	queryParams.Add("page", strconv.Itoa(page))

	baseURL, _ := url.Parse(MangaLibURL)
	baseURL.RawQuery = queryParams.Encode()

	return baseURL.String()
}

func createListURL(page int) string {
	queryParams := url.Values{}
	queryParams.Add("fields[]", "rate_avg")
	queryParams.Add("fields[]", "rate")
	queryParams.Add("fields[]", "releaseDate")
	queryParams.Add("site_id[]", "1")
	queryParams.Add("page", strconv.Itoa(page))

	baseURL, _ := url.Parse(MangaLibURL)
	baseURL.RawQuery = queryParams.Encode()

	return baseURL.String()
}

func createInfoURL(slug string, branchID int) string {
	queryParams := url.Values{}
	queryParams.Add("fields[]", "summary")
	queryParams.Add("fields[]", "releaseDate")
	queryParams.Add("fields[]", "views")
	queryParams.Add("fields[]", "genres")
	queryParams.Add("fields[]", "tags")
	queryParams.Add("fields[]", "teams")
	queryParams.Add("fields[]", "chap_count")
	queryParams.Add("fields[]", "authors")
	queryParams.Add("fields[]", "status_id")
	queryParams.Add("branch", strconv.Itoa(branchID))

	baseURL, _ := url.Parse(MangaLibURL + slug)
	baseURL.RawQuery = queryParams.Encode()

	return baseURL.String()
}

func createChaptersURL(slug string) string {
	return MangaLibURL + slug + "/chapters"
}

func createPageURL(image string) string {
	return CompressedURL + image
}

func createBranchesURL(id int) string {
	return BranchesURL + strconv.Itoa(id)
}

func createChapterURL(
	slug string, number, volume string, branchID int,
) string {
	queryParams := url.Values{}
	if branchID != 0 {
		queryParams.Add("branch_id", strconv.Itoa(branchID))
	}
	queryParams.Add("number", number)
	queryParams.Add("volume", volume)

	chapters := createChaptersURL(slug)
	chapter := chapters[:len(chapters)-1]
	baseURL, _ := url.Parse(chapter)
	baseURL.RawQuery = queryParams.Encode()

	return baseURL.String()
}
