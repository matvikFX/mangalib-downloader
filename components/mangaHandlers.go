package components

import (
	"context"
	"log/slog"

	"mangalib-downloader/api"
	"mangalib-downloader/components/utils"

	"github.com/gdamore/tcell/v2"
)

func (p *MangaPage) setHandlers(ctx context.Context, cancel context.CancelFunc, slug string) {
	select_change_row_color := func(row int) {
		// Я не знаю почему только так работает выделение нескольких столбцов
		// Пока оставлю так, если найду способ лучше, поменяю
		cols := []int{0, 1, 1, 2, 2}
		for _, col := range cols {
			cell := p.table.GetCell(row, col)
			if _, ok := p.selected[row]; ok {
				cell.SetBackgroundColor(tcell.ColorBlack)
				delete(p.selected, row)
			} else {
				cell.SetBackgroundColor(tcell.ColorRed)
				p.selected[row] = struct{}{}
			}
			p.table.SetCell(row, col, cell)
		}
	}

	manga, err := api.GetInfo(ctx, slug, p.app.branchID)
	if err != nil {
		slog.Error("Error receiving manga info", "Error", err)
		return
	}

	p.grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == ' ' {
			row, _ := p.table.GetSelection()
			select_change_row_color(row)
		}

		switch event.Key() {
		case tcell.KeyEscape: // Выход со страницы манги
			p.app.branchID = 0
			// timer.Reset(1 * time.Second)
			p.app.Pages.RemovePage(utils.MangaPageID)
			cancel()
		case tcell.KeyCtrlD: // Скачивание выделенных
			if len(p.selected) != 0 {
				go p.downloadSelected(ctx, manga)
			}
		case tcell.KeyCtrlA: // Скачивание всех глав
			go func() {
				p.app.downloader.DownloadManga(ctx, manga, p.app.branchID)

				p.app.ShowModal(utils.DownloadSuccessID,
					"Манга '"+manga.RusName+"' успешно скачана")

				go p.setChapters(ctx, manga, p.app.GetTeamsByID(manga.Branches))
			}()
		case tcell.KeyCtrlT: // Выбор ветки перевода
			branches := manga.Branches
			slug := manga.Slug
			if len(manga.Branches) > 0 {
				p.app.ShowBranchModal(ctx, slug, branches)
			} else {
				p.app.ShowModal(utils.NoBranchesID,
					"У данной менги нет другой ветки перевода")
			}
		}
		return event
	})

	// Выбор глав для скачивания
	p.table.SetSelectedFunc(func(row, _ int) {
		select_change_row_color(row)
	})
}
