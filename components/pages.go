package components

import (
	"context"

	"mangalib-downlaoder/components/utils"
	"mangalib-downlaoder/core"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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

func (p *MangaPage) setHandlers(cancel context.CancelFunc) {
	p.Grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape: // Выход со страницы манги
			cancel()
			core.App.PageHolder.RemovePage(utils.MangaPageID)
		case tcell.KeyCtrlD: // Скачивание выделенных
			func() {}()
		case tcell.KeyCtrlA: // Скачивание всех глав
			// core.App.Client.DownloadManga(p.manga)
		}
		return event
	})

	p.Table.SetSelectedFunc(func(row, _ int) {
		chapRef := p.Table.GetCell(row, 0).GetReference()
		if chapRef == nil {
			return
		}

		// Выбор главы для скачивания
		// if ch, ok := chapRef.(*models.Chapter); ok {
		// 	selected[ch] = struct{}
		// }
	})
}

func (p *SearchModal) setHandlers() {
	p.form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape: // Закрытие окна поиска
			core.App.PageHolder.RemovePage(utils.SearchModalID)
		case tcell.KeyEnter:
			searchInput := p.form.GetFormItemByLabel("Название").(*tview.InputField)
			formText := searchInput.GetText()
			core.App.Client.Query = formText

			searchInput.SetText("")
			core.App.PageHolder.RemovePage(utils.SearchModalID)
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
