package client

import (
	"context"

	"mangalib-downloader/models"
)

func (c *MangaLibClient) GetData(ctx context.Context) (*models.MangaListData, error) {
	jsonResp := &models.MangaListData{}

	var url string
	if c.Query == "" {
		url = c.createListURL()
	} else {
		url = c.createSearchURL(c.Query)
	}

	if err := c.ReqAndDecode(ctx, url, jsonResp); err != nil {
		return nil, err
	}

	return jsonResp, nil
}

func (c *MangaLibClient) GetMeta(ctx context.Context) (*models.Meta, error) {
	jsonResp, err := c.GetData(ctx)
	if err != nil {
		return nil, err
	}

	return jsonResp.Meta, nil
}

func (c *MangaLibClient) GetPopularManga(ctx context.Context) (models.MangaList, error) {
	jsonResp, err := c.GetData(ctx)
	if err != nil {
		return nil, err
	}

	return jsonResp.Manga, nil
}

func (c *MangaLibClient) GetSlugs(ctx context.Context) ([]string, error) {
	data, err := c.GetData(ctx)
	if err != nil {
		return nil, err
	}

	var slugs []string
	for _, manga := range data.Manga {
		slugs = append(slugs, manga.Slug)
	}

	return slugs, nil
}

func (c *MangaLibClient) GetInfo(ctx context.Context, slug string) (*models.MangaInfo, error) {
	mangaInfo := &models.MangaInfoData{}
	url := c.createInfoURL(slug)

	if err := c.ReqAndDecode(ctx, url, mangaInfo); err != nil {
		return nil, err
	}

	mangaInfo.Data.RusNameChange()
	return mangaInfo.Data, nil
}

func (c *MangaLibClient) GetMangaBranches(ctx context.Context, id int) (models.BranchList, error) {
	branches := &models.BranchesData{}
	url := c.createBranchesURL(id)

	if err := c.ReqAndDecode(ctx, url, branches); err != nil {
		return nil, err
	}

	return branches.Data, nil
}

func (c *MangaLibClient) GetChapters(ctx context.Context, slug string) (models.ChapterList, error) {
	chapters := &models.ChaptersData{}
	url := c.createChaptersURL(slug)

	if err := c.ReqAndDecode(ctx, url, chapters); err != nil {
		return nil, err
	}

	if c.Branch != 0 {
		chapList := make(models.ChapterList, 0)
		for _, chap := range chapters.Data {
			for _, br := range chap.Branches {
				if br.BranchID == c.Branch {
					chapList = append(chapList, chap)
				}
			}
		}
		return chapList, nil
	}

	return chapters.Data, nil
}

func (c *MangaLibClient) GetChaptersBranch(ctx context.Context, slug string, branch int) (models.ChapterList, error) {
	chapters := &models.ChaptersData{}
	url := c.createChaptersURL(slug)

	if err := c.ReqAndDecode(ctx, url, chapters); err != nil {
		return nil, err
	}

	var chaps models.ChapterList
	for _, chap := range chapters.Data {
		for _, br := range chap.Branches {
			if br.ID == branch {
				chaps = append(chaps, chap)
			}
		}
	}

	return chaps, nil
}

func (c *MangaLibClient) GetChapter(ctx context.Context, slug string, volume, number string) (*models.Chapter, error) {
	chapter := &models.ChapterData{}
	url := c.createChapterURL(slug, number, volume)

	if err := c.ReqAndDecode(ctx, url, chapter); err != nil {
		return nil, err
	}

	return chapter.Data, nil
}
