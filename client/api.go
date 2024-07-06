package client

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"mangalib-downlaoder/logger"
	"mangalib-downlaoder/models"
)

var Logger = logger.Logger{}

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

	return chapters.Data, nil
}

func (c *MangaLibClient) GetChapter(ctx context.Context, slug string, branchID int, volume, number string) (*models.Chapter, error) {
	chapter := &models.ChapterData{}
	url := c.createChapterURL(slug, branchID, number, volume)

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
		manga.RusName = c.removeChars(manga.RusName)
		chapPath := c.CreateChapterPath(branchTeams, manga.RusName, ch.Volume, ch.Number, ch.Name)
		if err = os.MkdirAll(chapPath, os.ModeDir); err != nil {
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

func (c *MangaLibClient) DownloadChapter(slug string, branch int, volume, number, chapPath string) {
	// Получение страниц
	chapter, err := c.GetChapter(context.Background(), slug, branch, volume, number)
	if err != nil {
		Logger.WriteLog(err.Error())
	}

	if err = os.MkdirAll(chapPath, os.ModeDir); err != nil {
		Logger.WriteLog(err.Error())
	}

	// Скачивание страниц
	wg := &sync.WaitGroup{}
	for _, p := range chapter.Pages {
		// Создание имени страницы
		pageName := c.createPageName(p.Slug, p.Image)

		// Если файл скачан, пропускаем
		// if c.CheckExistence(chapPath, pageName) {
		// 	continue
		// }

		// Создание пути для страницы
		pagePath := c.createPagePath(chapPath, pageName)
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
		fmt.Println("Error getting response from main server")
		Logger.WriteLog(err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading bytes")
		Logger.WriteLog(err.Error())
	}

	if err = c.createFile(body, pagePath); err != nil {
		fmt.Println("Error creating file")
		Logger.WriteLog(err.Error())
	}
}
