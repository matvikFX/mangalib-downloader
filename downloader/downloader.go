package downloader

import (
	"context"
	"os"
	"sync"

	"manga-downloader/api"
	"manga-downloader/models"
	"manga-downloader/services"
)

const workerNum = 4

var Path = DefaultDownloadPath()

type MangaLibDownloader struct {
	Logger *services.Logger

	Downloaded   chan struct{}
	DownloadPath string
}

func NewClient() *MangaLibDownloader {
	return &MangaLibDownloader{
		Logger: services.NewLogger(),

		Downloaded:   make(chan struct{}, 1),
		DownloadPath: DefaultDownloadPath(),
	}
}

func (c *MangaLibDownloader) DownloadManga(
	ctx context.Context, manga *models.MangaInfo, branchID int,
) {
	chapters, err := api.GetChapters(ctx, manga.Slug, branchID)
	if err != nil {
		c.Logger.Write(err.Error())
		return
	}

	c.DownloadChapters(ctx, manga.Manga, chapters, branchID)
}

func (c *MangaLibDownloader) DownloadChapters(ctx context.Context,
	manga models.Manga, chapters models.ChapterList, branchID int,
) {
	wg := &sync.WaitGroup{}
	branchTeams := api.GetBranchTeams(ctx, manga.ID)
	chapChan := make(chan *models.Chapter, workerNum)

	go func() {
		for _, chap := range chapters {
			chapChan <- chap
		}
		close(chapChan)
	}()

	for range workerNum {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.downloader(ctx, chapChan, manga, branchID, branchTeams)
		}()
	}

	go func() {
		wg.Wait()
		c.Downloaded <- struct{}{}
	}()
}

func (c *MangaLibDownloader) DownloadChapter(ctx context.Context,
	slug string, volume, number string, branchID int, chapPath string,
) {
	// Получение страниц
	chapter, err := api.GetChapter(ctx, slug, volume, number, branchID)
	if err != nil {
		c.Logger.Write(err.Error())
		return
	}

	if err = os.MkdirAll(chapPath, 0o755); err != nil {
		c.Logger.Write(err.Error())
		return
	}

	// Скачивание страниц
	wg := &sync.WaitGroup{}
	for _, p := range chapter.Pages {
		// Создание имени страницы
		pageName := createPageName(p.Slug, p.Image)
		// Создание пути для страницы
		pagePath := createPagePath(chapPath, pageName)

		// Если файл скачан, пропускаем
		if CheckExistence(pagePath) {
			continue
		}

		// Скачивание страницы
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			c.downloadPage(ctx, pagePath, url)
		}(p.URL)
	}
	wg.Wait()
}

func (c *MangaLibDownloader) downloader(ctx context.Context,
	chapChan <-chan *models.Chapter,
	manga models.Manga, branchID int, teams string,
) {
	for {
		select {
		case <-ctx.Done():
			c.Logger.Write(ctx.Err().Error())
			return
		case chap, ok := <-chapChan:
			if !ok {
				return
			}

			chapPath := CreateChapterPath(c.DownloadPath, teams, manga.RusName,
				chap.Volume, chap.Number, chap.Name)

			if err := os.MkdirAll(chapPath, 0o755); err != nil {
				c.Logger.Write(err.Error())
			}

			c.DownloadChapter(ctx, manga.Slug, chap.Volume, chap.Number, branchID, chapPath)
		}
	}
}

func (c *MangaLibDownloader) downloadPage(ctx context.Context, pagePath, pageURL string) {
	url := createPageURL(pageURL)
	img, err := api.ReqImg(ctx, url)
	if err != nil {
		c.Logger.Write(err.Error())
	}

	if err = createFile(img, pagePath); err != nil {
		c.Logger.Write(err.Error())
	}
}
