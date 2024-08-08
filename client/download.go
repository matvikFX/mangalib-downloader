package client

import (
	"context"
	"os"
	"sync"

	"mangalib-downloader/models"
)

const workerNum = 4

func (c *MangaLibClient) downloader(ctx context.Context,
	reciever <-chan *models.Chapter,
	manga models.Manga, teams string,
) {
	for {
		select {
		case <-ctx.Done():
			c.Logger.WriteLog(ctx.Err().Error())
			return
		case chap, ok := <-reciever:
			if !ok {
				return
			}

			chapPath := c.CreateChapterPath(teams, manga.RusName,
				chap.Volume, chap.Number, chap.Name)

			if err := os.MkdirAll(chapPath, 0o755); err != nil {
				c.Logger.WriteLog(err.Error())
			}

			c.DownloadChapter(ctx, manga.Slug, chap.Volume, chap.Number, chapPath)
		}
	}
}

func (c *MangaLibClient) DownloadManga(ctx context.Context, manga *models.MangaInfo) {
	chapters, err := c.GetChapters(ctx, manga.Slug)
	if err != nil {
		c.Logger.WriteLog(err.Error())
		return
	}

	c.DownloadChapters(ctx, manga.Manga, chapters)
}

func (c *MangaLibClient) DownloadChapters(ctx context.Context,
	manga models.Manga, chapters models.ChapterList,
) {
	wg := &sync.WaitGroup{}
	branchTeams := c.GetBranchTeams(ctx, manga.ID)
	chapChan := make(chan *models.Chapter, len(chapters))

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
			c.downloader(ctx, chapChan, manga, branchTeams)
		}()
	}

	go func() {
		wg.Wait()
		c.Downloaded <- struct{}{}
	}()
}

func (c *MangaLibClient) DownloadChapter(ctx context.Context,
	slug string, volume, number string, chapPath string,
) {
	// Получение страниц
	chapter, err := c.GetChapter(ctx, slug, volume, number)
	if err != nil {
		c.Logger.WriteLog(err.Error())
		return
	}

	if err = os.MkdirAll(chapPath, 0o755); err != nil {
		c.Logger.WriteLog(err.Error())
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
		if c.CheckExistence(pagePath) {
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

func (c *MangaLibClient) downloadPage(ctx context.Context, pagePath, pageURL string) {
	url := c.createPageURL(pageURL)
	img, err := c.ReqImg(ctx, url)
	if err != nil {
		c.Logger.WriteLog(err.Error())
	}

	if err = createFile(img, pagePath); err != nil {
		c.Logger.WriteLog(err.Error())
	}
}
