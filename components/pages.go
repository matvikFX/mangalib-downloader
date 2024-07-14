package components

import (
	"context"

	"mangalib-downloader/components/utils"
	"mangalib-downloader/core"
	"mangalib-downloader/logger"
	"mangalib-downloader/models"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var Logger = logger.Logger{}

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

func (p *ListPage) setHandlers(cancel context.CancelFunc) {
	p.table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		reload := func() {
			cancel()
			go p.setListTable()
		}

		switch event.Key() {
		case tcell.KeyEscape: // Обнулить поисковую строку
			if core.App.Client.Query != "" {
				core.App.Client.Query = ""
				reload()
			}
		case tcell.KeyCtrlF: // Предыдущая страница
			core.App.Client.Page++
			reload()

		case tcell.KeyCtrlB: // Следующая страница
			if core.App.Client.Page == 1 {
				// Показать модалку, что ниже 1 пойти нельзя
				break
			}
			core.App.Client.Page--
			reload()
		}
		return event
	})

	p.table.SetSelectedFunc(func(row, column int) {
		manga := p.table.GetCell(row, 0).GetReference().(*models.Manga)

		// Показать модалку с выбором ветки перевода
		ShowBranchModal(p.cWrap.Context, manga.Slug, manga.ID)
	})
}

func (p *MangaPage) setHandlers(cancel context.CancelFunc) {
	p.gird.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape: // Выход со страницы манги
			cancel()
			p.cWrap.Cancel()
			core.App.Client.Branch = 0
			core.App.PageHolder.RemovePage(utils.MangaPageID)
		case tcell.KeyCtrlD: // Скачивание выделенных
			p.downloadSelected(p.selected)
		case tcell.KeyCtrlA: // Скачивание всех глав
			core.App.Client.DownloadManga(p.cWrap.Context, p.manga)
		case tcell.KeyCtrlP: // Выбор ветки перевода
			if len(p.manga.Branches) > 0 {
				ShowBranchModal(p.cWrap.Context, p.manga.Slug, p.manga.ID)
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

func (p *BranchModal) setHandlers() {
	p.form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			core.App.PageHolder.RemovePage(utils.BranchModalID)
		case tcell.KeyEnter:
			SwitchToMangaPage(p.ctx, p.slug, p.id)
			core.App.PageHolder.RemovePage(utils.BranchModalID)
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
