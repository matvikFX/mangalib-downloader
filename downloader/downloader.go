package downloader

import (
	"context"
	"fmt"
	"log/slog"
	"manga-downloader/api"
	"manga-downloader/models"

	"golang.org/x/sync/errgroup"
)

type Downloader struct {
	downloadPath string
	cbzFormat    bool
}

func New(downloadPath string, cbzFormat bool) *Downloader {
	return &Downloader{
		downloadPath: downloadPath,
		cbzFormat:    cbzFormat,
	}
}

func (d *Downloader) DownloadManga(ctx context.Context,
	manga *models.MangaInfo, branchID int,
) error {
	log := slog.With("Downloader", "DownloadManga")

	chapters, err := api.GetChapters(ctx, manga.Slug, branchID)
	if err != nil {
		log.Error("Error receiving chapters", "Error", err)
		return err
	}
	log.Debug("Received chapters", "chaptersLen", len(chapters))

	if err := d.DownloadChapters(ctx, &manga.Manga, chapters, branchID); err != nil {
		log.Error("Error downloading chapters", "Error", err)
		return err
	}
	log.Info("Manga successfully downloaded", "manga", removeChars(manga.RusName))

	cbzPath := fmt.Sprintf("%s/%s", d.downloadPath, removeChars(manga.RusName))
	if d.cbzFormat {
		if err := CreateCBZArchive(cbzPath); err != nil {
			log.Error("Error converting to cbz format", "Error", err)
			return err
		}
	}
	log.Info("Manga added to CBZ archive")

	return nil
}

func (d *Downloader) DownloadChapters(ctx context.Context,
	manga *models.Manga, chapters models.ChapterList, branchID int,
) error {
	log := slog.With("Downloader", "DownloadChapters")

	branchTeams := api.GetBranchTeams(ctx, manga.ID)
	// chapChan := make(chan *models.Chapter, workerNum)
	chapChan := make(chan *models.Chapter)
	log.Debug("Received teams", "branchTeams", branchTeams)

	go func() {
		for _, chap := range chapters {
			chapChan <- chap
		}
		close(chapChan)
	}()

	g, ctx := errgroup.WithContext(ctx)
	for workerIdx := range workerNum {
		g.Go(func() error {
			log.Info(fmt.Sprintf("Starting worker %d", workerIdx))
			if err := d.worker(ctx, chapChan, manga, branchID, branchTeams); err != nil {
				log.Error(fmt.Sprintf("Error in worker %d", workerIdx))
				return err
			}

			log.Info(fmt.Sprintf("Worker %d stopped", workerIdx))
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}
	log.Info("Chapters successfully downloaded", "chapters", downloadedChapters(chapters))

	cbzPath := fmt.Sprintf("%s/%s", d.downloadPath, removeChars(manga.RusName))
	if d.cbzFormat {
		if err := CreateCBZArchive(cbzPath); err != nil {
			log.Error("Error converting to cbz format", "Error", err)
			return err
		}
	}
	log.Info("Manga added to CBZ archive")

	return nil
}

func (d *Downloader) DownloadChapter(ctx context.Context,
	slug string, volume, number string, branchID int, chapPath string,
) error {
	log := slog.With("Downloader", "DownloadChapters")

	// Получение страниц
	chapter, err := api.GetChapter(ctx, slug, volume, number, branchID)
	if err != nil {
		log.Error("Error receiving chapter", "Error", err)
		return err
	}

	// Скачивание страниц
	g, ctx := errgroup.WithContext(ctx)
	for _, p := range chapter.Pages {
		// Создание имени страницы
		pageName := createPageName(p.Number, p.Image)
		// Создание пути для страницы
		pagePath := createPagePath(chapPath, pageName)
		// log.Debug("Absolute chapter path", "pagePath", pagePath)

		// Если файл скачан, пропускаем
		if CheckExistence(pagePath) {
			log.Warn("Chapter already downloaded")
			continue
		}

		// Скачивание страницы
		g.Go(func() error {
			if err := d.downloadPage(ctx, pagePath, p.URL); err != nil {
				log.Error(
					fmt.Sprintf("Error downloading page %d", p.Number),
					"Error", err,
				)
				return err
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	log.Info(fmt.Sprintf(
		"Chapter %s-%s successfully downloaded",
		chapter.Volume, chapter.Number,
	))
	return nil
}
