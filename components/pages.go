package components

import (
	"context"
	"time"

	"mangalib-downloader/components/utils"
	"mangalib-downloader/core"
	"mangalib-downloader/models"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	timer *time.Timer

	selectedManga = &models.MangaInfo{}
)

func SetHandlers() {
	core.App.TView.EnableMouse(true)
	core.App.TView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'P': // Пути
			ShowPathModal()
		case 'H': // Помощь
			ShowHelpPage()
		}

		switch event.Key() {
		case tcell.KeyCtrlS: // Поиск по названию
			ShowSearchModal()
		case tcell.KeyCtrlC: // Завершение рабоыт
			core.App.TView.Stop()
		}

		return event
	})
}

func (p *ListPage) setHandlers(ctx context.Context, cancel context.CancelFunc) {
	p.table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		reload := func() {
			cancel()
			go p.setListTable()
		}

		if event.Rune() == ' ' {
			if timer != nil {
				timer.Stop()
			}
			ShowBranchModal(ctx)
		}

		switch event.Key() {
		case tcell.KeyEscape: // Обнулить поисковую строку и вернуться на первую страницу
			if core.App.Client.Query != "" || core.App.Client.Page != 1 {
				core.App.Client.Query = ""
				core.App.Client.Page = 1
				reload()
			}
		case tcell.KeyCtrlF: // Предыдущая страница
			core.App.Client.Page++
			reload()
		case tcell.KeyCtrlB: // Следующая страница
			if core.App.Client.Page == 1 {
				// Показать модалку, что ниже 1 пойти нельзя
				ShowModal(utils.NoMangaID, "Ниже первой страницы опуститься нельзя")
				break
			}
			core.App.Client.Page--
			reload()
		}
		return event
	})

	p.table.SetSelectedFunc(func(_, _ int) {
		if timer != nil {
			timer.Stop()
		}
		ShowBranchModal(ctx)
	})

	p.table.SetSelectionChangedFunc(func(row, column int) {
		manga := p.getMangaFromCell(row)
		selectedManga = &models.MangaInfo{
			Manga: *manga,
		}

		p.textView.SetTitle("Загрузка информации о манге...")
		p.textView.SetText("")

		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(800*time.Millisecond, func() {
			info, err := core.App.Client.GetInfo(ctx, manga.Slug)
			if err != nil {
				core.App.Client.Logger.WriteLog(err.Error())
				return
			}
			selectedManga = info

			infoText := utils.InfoText(info, nil)
			core.App.TView.QueueUpdateDraw(func() {
				p.textView.SetTitle("Информация о манге")
				p.textView.SetText(infoText)
			})
		})
	})
}

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
			core.App.Client.Branch = 0
			timer.Reset(1 * time.Second)
			core.App.PageHolder.RemovePage(utils.MangaPageID)
			cancel()
		case tcell.KeyCtrlD: // Скачивание выделенных
			if len(p.selected) != 0 {
				go p.downloadSelected(ctx)
			}
		case tcell.KeyCtrlA: // Скачивание всех глав
			go func() {
				core.App.Client.DownloadManga(ctx, selectedManga)

				<-core.App.Client.Downloaded
				ShowModal(utils.DownloadSuccessID,
					"Манга '"+selectedManga.RusName+"' успешно скачана")

				go p.setChapters(ctx)
			}()
		case tcell.KeyCtrlP: // Выбор ветки перевода
			if len(selectedManga.Branches) > 0 {
				ShowBranchModal(ctx)
			} else {
				ShowModal(utils.NoBranchesID,
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

func (p *SearchModal) setHandlers() {
	p.form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape: // Закрытие страницы поиска
			core.App.PageHolder.RemovePage(utils.SearchModalID)
		case tcell.KeyEnter:
			searchInput := p.form.GetFormItemByLabel(utils.SearchModalLabel).(*tview.InputField)
			formText := searchInput.GetText()
			core.App.Client.Query = formText

			searchInput.SetText("")
			core.App.PageHolder.RemovePage(utils.SearchModalID)
			ShowListPage()
		}
		return event
	})
}

func (p *PathModal) setHandlers() {
	p.form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			core.App.PageHolder.RemovePage(utils.PathsModalID)
		}
		return event
	})
}

func (p *HelpPage) setHandlers() {
	p.Grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			core.App.PageHolder.RemovePage(utils.HelpPageID)
		}
		return event
	})
}
