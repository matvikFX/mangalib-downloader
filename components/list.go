package components

import (
	"context"
	"fmt"
	"unicode/utf8"

	"manga-downloader/components/utils"
	"manga-downloader/models"

	"github.com/rivo/tview"
)

type ListPage struct {
	app *TViewApp

	grid     *tview.Grid
	textView *tview.TextView
	table    *tview.Table
}

func (a *TViewApp) ShowListPage() {
	listPage := newListPage(a)

	a.App.SetFocus(listPage.grid)
	a.Pages.AddAndSwitchToPage(utils.ListPageID, listPage.grid, true)
}

func newListPage(app *TViewApp) *ListPage {
	textView := tview.NewTextView()
	textView.SetWrap(true).SetWordWrap(true).
		SetTitle("Информация о манге").SetBorder(true)

	table := tview.NewTable()
	table.SetSelectable(true, false).
		SetSeparator('|').
		SetBorder(true)

	grid := tview.NewGrid()
	grid.SetRows(-1).SetColumns(-1, -1, -1, -1, -1, -1, -1, -1, -1).
		SetTitle("Список манги").SetBorder(true)

	grid.AddItem(table, 0, 0, 1, 3, 0, 0, true).
		AddItem(textView, 0, 3, 1, 6, 0, 0, false)

	listPage := &ListPage{
		app: app,

		grid:     grid,
		textView: textView,
		table:    table,
	}

	go listPage.setListTable()

	return listPage
}

func (p *ListPage) setListTable() {
	ctx, cancel := context.WithCancel(context.Background())
	p.setHandlers(ctx, cancel)

	tableTitle := "Популярная манга"
	if p.app.Downloader.Query != "" {
		tableTitle = "Результаты поиска"
	}

	p.app.App.QueueUpdateDraw(func() {
		p.table.SetTitle(fmt.Sprintf("%s. Загрузка...", tableTitle))
	})

	data, err := p.app.Downloader.GetData(ctx)
	if err != nil {
		p.app.Logger.Write(err.Error())
		return
	}

	meta := data.Meta
	manga := data.Manga

	if meta.From == 0 {
		p.app.ShowModal(utils.NoMangaID, "Манга не найдена")
		if p.app.Downloader.Page == 1 {
			p.app.Downloader.Query = ""
		} else {
			p.app.Downloader.Page--
		}
		go p.setListTable()
		return
	}

	p.table.SetTitle(fmt.Sprintf("%s. Страница %d (%d-%d)",
		tableTitle, meta.Page, meta.From, meta.To))

	for idx, manga := range manga {
		manga.RusNameChange()
		if utf8.RuneCountInString(manga.RusName) > 60 {
			runes := []rune(manga.RusName)
			manga.RusName = string(runes[:60])
		}
		title := tview.NewTableCell(
			fmt.Sprintf("%-60s", manga.RusName)).
			SetMaxWidth(60).SetReference(manga)

		p.table.SetCell(idx, 0, title)
	}

	p.app.App.QueueUpdateDraw(func() {
		p.table.Select(0, 0)
		p.table.ScrollToBeginning()
	})
}

func (p *ListPage) getMangaFromCell(row int) *models.Manga {
	return p.table.GetCell(row, 0).GetReference().(*models.Manga)
}
