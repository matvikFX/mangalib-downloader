package api

import (
	"context"
	"log/slog"
	"manga-downloader/models"
)

func GetData(
	ctx context.Context, query string, page int,
) (*models.MangaListData, error) {
	jsonResp := &models.MangaListData{}

	var url string
	if query == "" {
		url = createListURL(page)
	} else {
		url = createSearchURL(query, page)
	}

	if err := ReqJSON(ctx, url, jsonResp); err != nil {
		return nil, err
	}

	return jsonResp, nil
}

func GetMeta(
	ctx context.Context, query string, page int,
) (*models.Meta, error) {
	jsonResp, err := GetData(ctx, query, page)
	if err != nil {
		return nil, err
	}

	return jsonResp.Meta, nil
}

func GetPopularManga(
	ctx context.Context, query string, page int,
) (models.MangaList, error) {
	jsonResp, err := GetData(ctx, query, page)
	if err != nil {
		return nil, err
	}

	return jsonResp.Manga, nil
}

func GetSlugs(
	ctx context.Context, query string, page int,
) ([]string, error) {
	data, err := GetData(ctx, query, page)
	if err != nil {
		return nil, err
	}

	var slugs []string
	for _, manga := range data.Manga {
		slugs = append(slugs, manga.Slug)
	}

	return slugs, nil
}

func GetInfo(
	ctx context.Context, slug string, branchID int,
) (*models.MangaInfo, error) {
	mangaInfo := &models.MangaInfoData{}
	url := createInfoURL(slug, branchID)

	if err := ReqJSON(ctx, url, mangaInfo); err != nil {
		return nil, err
	}

	mangaInfo.Data.RusNameChange()
	return mangaInfo.Data, nil
}

func GetMangaBranches(
	ctx context.Context, id int,
) (models.BranchList, error) {
	branches := &models.BranchesData{}
	url := createBranchesURL(id)

	if err := ReqJSON(ctx, url, branches); err != nil {
		return nil, err
	}

	return branches.Data, nil
}

func GetChapters(
	ctx context.Context, slug string, branchID int,
) (models.ChapterList, error) {
	chapters := &models.ChaptersData{}
	url := createChaptersURL(slug)

	if err := ReqJSON(ctx, url, chapters); err != nil {
		return nil, err
	}

	if branchID != 0 {
		chapList := make(models.ChapterList, 0)
		for _, chap := range chapters.Data {
			for _, br := range chap.Branches {
				if br.BranchID == branchID {
					chapList = append(chapList, chap)
				}
			}
		}
		return chapList, nil
	}

	return chapters.Data, nil
}

// По какой-то причине не используется
func GetChaptersBranch(
	ctx context.Context, slug string, branchID int,
) (models.ChapterList, error) {
	chapters := &models.ChaptersData{}
	url := createChaptersURL(slug)

	if err := ReqJSON(ctx, url, chapters); err != nil {
		return nil, err
	}

	var chaps models.ChapterList
	for _, chap := range chapters.Data {
		for _, branch := range chap.Branches {
			if branch.ID == branchID {
				chaps = append(chaps, chap)
			}
		}
	}

	return chaps, nil
}

func GetChapter(
	ctx context.Context, slug string, volume, number string, branchID int,
) (*models.Chapter, error) {
	chapter := &models.ChapterData{}
	url := createChapterURL(slug, number, volume, branchID)

	if err := ReqJSON(ctx, url, chapter); err != nil {
		return nil, err
	}

	return chapter.Data, nil
}

func GetBranchTeams(ctx context.Context, branchID int) string {
	branchTeams := make(map[int]string)
	if branchID != 0 {
		branches, err := GetMangaBranches(ctx, branchID)
		if err != nil {
			slog.Debug("Error receiving branches", "Error", err)
		}

		branchTeams = branches.BranchTeams()
	}

	return branchTeams[branchID]
}
