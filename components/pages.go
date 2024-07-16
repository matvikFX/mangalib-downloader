package components

import (
	"context"
	"time"

	"mangalib-downloader/components/utils"
	"mangalib-downloader/core"
	"mangalib-downloader/logger"
	"mangalib-downloader/models"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	timer *time.Timer

	Logger        = logger.Logger{}
	selectedManga = &models.MangaInfo{}
)

func SetHandlers() {
	core.App.TView.EnableMouse(true)
	core.App.TView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlS: // Поиск по названию
			ShowSearchModal()
		case tcell.KeyCtrlO: // Помощь
			ShowHelpPage()
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

	p.table.SetSelectedFunc(func(row, _ int) {
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
		timer = time.AfterFunc(600*time.Millisecond, func() {
			info, err := core.App.Client.GetInfo(p.cWrap.Context, manga.Slug)
			if err != nil {
				Logger.WriteLog(err.Error())
				return
			}
			selectedManga = info

			infoText := utils.ListInfoText(info)
			core.App.TView.QueueUpdateDraw(func() {
				p.textView.SetTitle("Информация о манге")
				p.textView.SetText(infoText)
			})
		})
	})
}

func (p *MangaPage) setHandlers(ctx context.Context, cancel context.CancelFunc) {
	p.grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape: // Выход со страницы манги
			cancel()
			core.App.Client.Branch = 0
			core.App.PageHolder.RemovePage(utils.MangaPageID)
		case tcell.KeyCtrlD: // Скачивание выделенных
			p.downloadSelected()
		case tcell.KeyCtrlA: // Скачивание всех глав
			core.App.Client.DownloadManga(ctx, selectedManga)
		case tcell.KeyCtrlP: // Выбор ветки перевода
			if len(selectedManga.Branches) > 0 {
				ShowBranchModal(ctx)
			} else {
				ShowModal("NoTranslateBranch",
					"У данной менги нет другой ветки перевода")
			}
		}
		return event
	})

	// Выбор глав для скачивания
	p.table.SetSelectedFunc(func(row, _ int) {
		cell := p.table.GetCell(row, 0)
		if p.selected[row] {
			cell.SetBackgroundColor(tcell.ColorBlack)
			p.selected[row] = false
		} else {
			cell.SetBackgroundColor(tcell.ColorRed)
			p.selected[row] = true
		}
		p.table.SetCell(row, 0, cell)
	})
}

func (p *SearchModal) setHandlers() {
	p.form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape: // Закрытие окна поиска
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

func (p *HelpPage) setHandlers() {
	p.Grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			core.App.PageHolder.RemovePage(utils.HelpPageID)
		}
		return event
	})
}
