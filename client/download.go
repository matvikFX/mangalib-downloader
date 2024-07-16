package client

import (
	"context"
	"io"
	"os"
	"sync"

	"mangalib-downloader/models"
)

func (c *MangaLibClient) DownloadManga(ctx context.Context, manga *models.MangaInfo) {
	// Получение глав
	chapters, err := c.GetChapters(ctx, manga.Slug)
	if err != nil {
		Logger.WriteLog(err.Error())
	}

	branchTeams := c.GetBranchTeams(ctx, manga.ID)

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
			c.DownloadChapter(ctx, manga.Slug, ch.Volume, ch.Number, chapPath)
		}(ch)
	}
	wg.Wait()
}

func (c *MangaLibClient) DownloadChapters(ctx context.Context, mangaID int, mangaSlug, mangaName string, chapters models.ChapterList) {
	branchTeams := c.GetBranchTeams(ctx, mangaID)

	wg := &sync.WaitGroup{}
	for _, ch := range chapters {
		chapPath := c.CreateChapterPath(branchTeams, mangaName, ch.Volume, ch.Number, ch.Name)
		wg.Add(1)
		go func(vol, num string) {
			defer wg.Done()
			c.DownloadChapter(ctx, mangaSlug, vol, num, chapPath)
		}(ch.Volume, ch.Number)
	}
	wg.Wait()
}

func (c *MangaLibClient) DownloadChapter(ctx context.Context, slug string, volume, number, chapPath string) {
	// Получение страниц
	chapter, err := c.GetChapter(ctx, slug, volume, number)
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
