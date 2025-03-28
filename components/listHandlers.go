package components

import (
	"context"
	"time"

	"mangalib-downloader/api"
	"mangalib-downloader/components/utils"
	"mangalib-downloader/models"

	"github.com/gdamore/tcell/v2"
)

var timer *time.Timer

func (p *ListPage) setHandlers(ctx context.Context, cancel context.CancelFunc) {
	p.table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		reload := func() {
			cancel()
			go p.setListTable()
		}

		switch event.Key() {
		case tcell.KeyEscape: // Обнулить поисковую строку и вернуться на первую страницу
			if p.app.query != "" || p.app.page != 1 {
				p.app.query = ""
				p.app.page = 1
				reload()
			}
		case tcell.KeyCtrlF: // Предыдущая страница
			p.app.page++
			reload()
		case tcell.KeyCtrlB: // Следующая страница
			if p.app.page == 1 {
				p.app.ShowModal(utils.NoMangaID, "Ниже первой страницы опуститься нельзя")
				break
			}
			p.app.page--
			reload()
		}

		return event
	})

	p.table.SetSelectedFunc(func(row, column int) {
		manga := p.getMangaFromCell(row)
		p.app.selectedManga = &models.MangaInfo{
			Manga: *manga,
		}

		branches, err := api.GetMangaBranches(ctx, p.app.selectedManga.ID)
		if err != nil {
			return
		}

		if len(branches) == 0 {
			p.app.ShowMangaPage(ctx, p.app.selectedManga.Slug, 0)
			return
		} else {
			p.app.ShowBranchModal(ctx, p.app.selectedManga.Slug, branches)
		}
	})

	p.table.SetSelectionChangedFunc(func(row, column int) {
		manga := p.getMangaFromCell(row)
		p.app.selectedManga = &models.MangaInfo{
			Manga: *manga,
		}

		p.textView.SetTitle("Загрузка информации о манге...")
		p.textView.SetText("")

		if timer != nil {
			timer.Stop()
		}

		timer = time.AfterFunc(600*time.Millisecond, func() {
			mangaInfo, err := loadMangaInfo(ctx, p.app.selectedManga.Slug, p.app.branchID)
			if err != nil {
				return
			}
			p.app.selectedManga = mangaInfo

			infoText := utils.InfoText(mangaInfo, nil)
			p.app.App.QueueUpdateDraw(func() {
				p.textView.SetTitle("Информация о манге")
				p.textView.SetText(infoText)
			})
		})
	})
}

func loadMangaInfo(ctx context.Context, slug string, branchID int) (*models.MangaInfo, error) {
	info, err := api.GetInfo(ctx, slug, branchID)
	if err != nil {
		return nil, err
	}

	branches, err := api.GetMangaBranches(ctx, info.ID)
	if err != nil {
		return nil, err
	}

	info.Branches = branches
	return info, nil
}
