package client

import (
	"context"
	"io"
	"runtime"
	"sync"
	"time"

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

	chapList := make(models.ChapterList, 0)
	if c.Branch != 0 {
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

// Не думаю, что понадобится
func (c *MangaLibClient) GetPageBytes(pageURL string) ([]byte, error) {
	url := c.createPageURL(pageURL)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Поличение списка без горутин
func (c *MangaLibClient) GetListInfo(ctx context.Context, slugs []string) models.MangaInfoList {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var list models.MangaInfoList
	for _, slug := range slugs {
		info, err := c.GetInfo(ctx, slug)
		if err != nil {
			Logger.WriteLog(err.Error())
		}

		list = append(list, info)
	}

	return list
}

var (
	cpuCores = runtime.NumCPU()
	timeout  = 10 * time.Second
	reqLimit = 5
)

func (c *MangaLibClient) GetInfoList(ctx context.Context, slugs []string) models.MangaInfoList {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg := &sync.WaitGroup{}
	slugChan := make(chan string, reqLimit)
	infoChan := make(chan *models.MangaInfo, reqLimit)

	// Чтение из одного канала и запись в другой
	for range cpuCores {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.infoWriter(ctx, slugChan, infoChan)
		}()
	}

	// Запись названий в канал
	go func() {
		for idx, slug := range slugs {
			// Приостанвока горутины при заполнении канала
			if (idx+1)%10 == 0 {
				time.Sleep(500 * time.Millisecond)
			}
			slugChan <- slug
		}
		close(slugChan)
	}()

	// Ожидание завершения горутин и закрытие канала
	go func() {
		wg.Wait()
		close(infoChan)
	}()

	// Получение информации из канала
	var list models.MangaInfoList
	for info := range infoChan {
		list = append(list, info)
	}

	return list
}

// Получение данных из одного канала и запись в другой
func (c *MangaLibClient) infoWriter(ctx context.Context, slugChan <-chan string, infoChan chan<- *models.MangaInfo) {
	for {
		select {
		case <-ctx.Done():
			return
		case slug, open := <-slugChan:
			if !open {
				return
			}

			info, err := c.GetInfo(ctx, slug)
			if err != nil {
				Logger.WriteLog(err.Error())
			}
			infoChan <- info
		}
	}
}
