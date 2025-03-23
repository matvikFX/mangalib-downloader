package api

import (
	"context"
	"log/slog"
	"mangalib-downloader/models"
)

func GetData(
	ctx context.Context, query string, page int,
) (*models.MangaListData, error) {
	log := slog.With("API", "GetData")

	var url string
	if query == "" {
		url = createListURL(page)
	} else {
		url = createSearchURL(query, page)
	}

	jsonResp := &models.MangaListData{}
	if err := ReqJSON(ctx, url, jsonResp); err != nil {
		log.Error("Error receiving data", "Error", err)
		return nil, err
	}

	return jsonResp, nil
}

func GetMeta(
	ctx context.Context, query string, page int,
) (*models.Meta, error) {
	log := slog.With("API", "GetMeta")

	jsonResp, err := GetData(ctx, query, page)
	if err != nil {
		log.Error("Error receiving meta", "Error", err)
		return nil, err
	}

	return jsonResp.Meta, nil
}

func GetPopularManga(
	ctx context.Context, query string, page int,
) (models.MangaList, error) {
	log := slog.With("API", "GetPopularManga")

	jsonResp, err := GetData(ctx, query, page)
	if err != nil {
		log.Error("Error receiving popular manga", "Error", err)
		return nil, err
	}

	return jsonResp.Manga, nil
}

func GetSlugs(
	ctx context.Context, query string, page int,
) ([]string, error) {
	log := slog.With("API", "GetSlugs")

	data, err := GetData(ctx, query, page)
	if err != nil {
		log.Error("Error receiving manga info", "Error", err)
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
	log := slog.With("API", "GetInfo")

	mangaInfo := &models.MangaInfoData{}
	url := createInfoURL(slug, branchID)

	if err := ReqJSON(ctx, url, mangaInfo); err != nil {
		log.Error("Error receiving manga info", "Error", err)
		return nil, err
	}

	mangaInfo.Data.RusNameChange()
	return mangaInfo.Data, nil
}

func GetMangaBranches(
	ctx context.Context, id int,
) (models.BranchList, error) {
	log := slog.With("API", "GetMangaBranches")

	branches := &models.BranchesData{}
	url := createBranchesURL(id)

	if err := ReqJSON(ctx, url, branches); err != nil {
		log.Error("Error receiving branches", "Error", err)
		return nil, err
	}

	return branches.Data, nil
}

func GetChapters(
	ctx context.Context, slug string, branchID int,
) (models.ChapterList, error) {
	log := slog.With("API", "GetChapters")

	chapters := &models.ChaptersData{}
	url := createChaptersURL(slug)

	if err := ReqJSON(ctx, url, chapters); err != nil {
		log.Error("Error receiving chapters", "Error", err)
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
	log := slog.With("API", "GetChaptersBranch")

	chapters := &models.ChaptersData{}
	url := createChaptersURL(slug)

	if err := ReqJSON(ctx, url, chapters); err != nil {
		log.Error("Error receiving chapters", "Error", err)
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
	log := slog.With("API", "GetChapter")

	chapter := &models.ChapterData{}
	url := createChapterURL(slug, number, volume, branchID)

	if err := ReqJSON(ctx, url, chapter); err != nil {
		log.Error("Error receiving chapter", "Error", err)
		return nil, err
	}

	return chapter.Data, nil
}

func GetBranchTeams(ctx context.Context, branchID int) string {
	log := slog.With("API", "GetBranchTeams")

	branchTeams := make(map[int]string)
	if branchID != 0 {
		branches, err := GetMangaBranches(ctx, branchID)
		if err != nil {
			log.Error("Error receiving branches", "Error", err)
		}

		branchTeams = branches.BranchTeams()
	}

	return branchTeams[branchID]
}
