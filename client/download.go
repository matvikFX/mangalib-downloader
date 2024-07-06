package client

import (
	"context"
	"io"
	"os"
	"strings"
	"sync"

	"mangalib-downlaoder/models"
)

func (c *MangaLibClient) DownloadManga(manga *models.MangaInfo) {
	ctx := context.Background()

	branches, err := c.GetMangaBranches(ctx, manga.ID)
	if err != nil {
		Logger.WriteLog(err.Error())
	}
	manga.Branches = branches

	// branch := make(map[int]string)
	// for _, b := range manga.Branches {
	// 	teams := make([]string, 0)
	// 	for _, t := range b.Teams {
	// 		teams = append(teams, t.Name)
	// 	}
	// 	branch[b.BranchID] = strings.Join(teams, ",")
	// }

	// Выбор ветки перевода,
	// Если ветка одна выбра нет
	teams := []string{}
	var branchTeams string
	if len(branches) != 0 {
		for _, t := range manga.Branches[0].Teams {
			teams = append(teams, t.Name)
		}
		// объединение переводчиков в ветке
		branchTeams = strings.Join(teams, ",")
	}

	// Получение глав
	chapters, err := c.GetChapters(ctx, manga.Slug)
	if err != nil {
		Logger.WriteLog(err.Error())
	}

	wg := &sync.WaitGroup{}
	for _, ch := range chapters {
		manga.RusName = removeChars(manga.RusName)
		chapPath := c.CreateChapterPath(branchTeams, manga.RusName, ch.Volume, ch.Number, ch.Name)
		if err = os.MkdirAll(chapPath, 0o755); err != nil {
			Logger.WriteLog(err.Error())
		}

		// Скачивание главы
		wg.Add(1)
		go func(ch *models.Chapter) {
			wg.Done()
			c.DownloadChapter(manga.Slug, ch.Branches[0].BranchID, ch.Volume, ch.Number, chapPath)
		}(ch)
	}
	wg.Wait()
}

func (c *MangaLibClient) DownloadChapters(mangaSlug, mangaName string, chapters models.ChapterList) {
	wg := &sync.WaitGroup{}
	for _, ch := range chapters {
		chapPath := c.CreateChapterPath("", mangaName, ch.Volume, ch.Number, ch.Name)
		wg.Add(1)
		go func(vol, num string) {
			defer wg.Done()
			c.DownloadChapter(mangaSlug, 0, vol, num, chapPath)
		}(ch.Volume, ch.Number)
	}
	wg.Wait()
}

func (c *MangaLibClient) DownloadChapter(slug string, branch int, volume, number, chapPath string) {
	// Получение страниц
	chapter, err := c.GetChapter(context.Background(), slug, branch, volume, number)
	if err != nil {
		Logger.WriteLog(err.Error())
	}

	if err = os.MkdirAll(chapPath, 0o755); err != nil {
		Logger.WriteLog(err.Error())
	}

	// Скачивание страниц
	wg := &sync.WaitGroup{}
	for _, p := range chapter.Pages {
		// Создание имени страницы
		pageName := createPageName(p.Slug, p.Image)

		// Если файл скачан, пропускаем
		// if c.CheckExistence(chapPath, pageName) {
		// 	continue
		// }

		// Создание пути для страницы
		pagePath := createPagePath(chapPath, pageName)
		// Если файл скачан, пропускаем
		if _, err := os.Stat(pagePath); !os.IsNotExist(err) {
			continue
		}

		// Скачивание страницы
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			c.downloadPage(pagePath, url)
		}(p.URL)
	}
	wg.Wait()
}

func (c *MangaLibClient) downloadPage(pagePath, pageURL string) {
	// url := "https://img33.imgslib.link/" + pageURL
	url := c.createPageURL(pageURL)
	resp, err := c.client.Get(url)
	if err != nil {
		Logger.WriteLog(err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		Logger.WriteLog(err.Error())
	}

	if err = createFile(body, pagePath); err != nil {
		Logger.WriteLog(err.Error())
	}
}
