package client

import (
	"net/url"
	"strconv"
)

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

func (c *MangaLibClient) createSearchURL(name string) string {
	queryParams := url.Values{}
	queryParams.Add("fields[]", "rate_avg")
	queryParams.Add("fields[]", "rate")
	queryParams.Add("fields[]", "releaseDate")
	queryParams.Add("site_id[]", "1")
	queryParams.Add("q", name)
	queryParams.Add("page", strconv.Itoa(c.Page))

	baseURL, _ := url.Parse(MangaLibURL)
	baseURL.RawQuery = queryParams.Encode()

	return baseURL.String()
}

func (c *MangaLibClient) createListURL() string {
	queryParams := url.Values{}
	queryParams.Add("fields[]", "rate_avg")
	queryParams.Add("fields[]", "rate")
	queryParams.Add("fields[]", "releaseDate")
	queryParams.Add("site_id[]", "1")
	queryParams.Add("page", strconv.Itoa(c.Page))

	baseURL, _ := url.Parse(MangaLibURL)
	baseURL.RawQuery = queryParams.Encode()

	return baseURL.String()
}

func (c *MangaLibClient) createInfoURL(slug string) string {
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
	queryParams.Add("branch", strconv.Itoa(c.Branch))

	baseURL, _ := url.Parse(MangaLibURL + slug)
	baseURL.RawQuery = queryParams.Encode()

	return baseURL.String()
}

func (c *MangaLibClient) createChaptersURL(slug string) string {
	return MangaLibURL + slug + "/chapters"
}

func (c *MangaLibClient) createPageURL(image string) string {
	return SecondURL + image
}

func (c *MangaLibClient) createBranchesURL(id int) string {
	return BranchesURL + strconv.Itoa(id)
}

func (c *MangaLibClient) createChapterURL(slug string, number, volume string) string {
	queryParams := url.Values{}
	if c.Branch != 0 {
		queryParams.Add("branch_id", strconv.Itoa(c.Branch))
	}
	queryParams.Add("number", number)
	queryParams.Add("volume", volume)

	chapters := c.createChaptersURL(slug)
	chapter := chapters[:len(chapters)-1]
	baseURL, _ := url.Parse(chapter)
	baseURL.RawQuery = queryParams.Encode()

	return baseURL.String()
}
