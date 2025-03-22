package components

import (
	"context"
	"mangalib-downloader/components/utils"

	"github.com/gdamore/tcell/v2"
)

func (p *MangaPage) setHandlers(ctx context.Context, cancel context.CancelFunc) {
	select_change_row_color := func(row int) {
		// Я не знаю почему только так работает выделение нескольких столбцов
		// Пока оставлю так, если найду способ лучше, поменяю
		cols := []int{0, 1, 1, 2, 2}
		for _, col := range cols {
			cell := p.table.GetCell(row, col)
			if p.selected[row] {
				cell.SetBackgroundColor(tcell.ColorBlack)
				delete(p.selected, row)
			} else {
				cell.SetBackgroundColor(tcell.ColorRed)
				p.selected[row] = true
			}
			p.table.SetCell(row, col, cell)
		}
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
				go p.downloadSelected(ctx, p.app.branchID)
			}
		case tcell.KeyCtrlA: // Скачивание всех глав
			go func() {
				p.app.downloader.DownloadManga(ctx, p.app.selectedManga, p.app.branchID)

				// <-p.app.downloader.Downloaded
				p.app.ShowModal(utils.DownloadSuccessID,
					"Манга '"+p.app.selectedManga.RusName+"' успешно скачана")

				go p.setChapters(ctx)
			}()
		case tcell.KeyCtrlT: // Выбор ветки перевода
			branches := p.app.selectedManga.Branches
			slug := p.app.selectedManga.Slug
			if len(p.app.selectedManga.Branches) > 0 {
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
