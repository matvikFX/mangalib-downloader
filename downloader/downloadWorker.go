package downloader

import (
	"context"
	"fmt"
	"manga-downloader/api"
	"manga-downloader/models"
)

func (d *Downloader) worker(ctx context.Context,
	chapChan <-chan *models.Chapter,
	manga *models.Manga, branchID int, teams string,
) error {
	log := d.logger.With("Downloader", "worker")

	for {
		select {
		case _ = <-ctx.Done():
			log.Error("Context done", "Error", ctx.Err())
			return ctx.Err()
		case chap, ok := <-chapChan:
			if !ok {
				log.Warn("No chapters received")
				return nil
			}

			chapPath := CreateChapterPath(d.downloadPath, teams, manga.RusName,
				chap.Volume, chap.Number, chap.Name)

			if err := d.DownloadChapter(ctx, manga.Slug,
				chap.Volume, chap.Number,
				branchID, chapPath,
			); err != nil {
				log.Error("Error downloading chapter", "Error", err)
				return err
			}
		}
	}
}

func (d *Downloader) downloadPage(ctx context.Context, pagePath, pageURL string) error {
	log := d.logger.With("Downloader", "downloadPage")

	url := createPageURL(pageURL)
	log.Debug("Pages URL", "url", url)

	img, err := api.ReqImg(ctx, url)
	if err != nil {
		log.Error("Error receiving image", "Error", err)
		return err
	}
	log.Debug("Size of the image", "imgSize", len(img))

	if err = createFile(img, pagePath); err != nil {
		log.Error("Error createing file", "Error", err)
		return err
	}

	log.Debug(fmt.Sprintf("Image %s successfully downloaded", pageURL))
	return nil
}
