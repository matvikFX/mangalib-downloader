package components

import (
	"context"
	"fmt"
	"unicode/utf8"

	"mangalib-downloader/components/utils"
	"mangalib-downloader/core"
	"mangalib-downloader/models"

	"github.com/rivo/tview"
)

type ListPage struct {
	grid     *tview.Grid
	textView *tview.TextView
	table    *tview.Table
}

func ShowListPage() {
	listPage := newListPage()

	core.App.TView.SetFocus(listPage.grid)
	core.App.PageHolder.AddAndSwitchToPage(utils.ListPageID, listPage.grid, true)
}

func newListPage() *ListPage {
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
	if core.App.Client.Query != "" {
		tableTitle = "Результаты поиска"
	}

	core.App.TView.QueueUpdateDraw(func() {
		p.table.SetTitle(fmt.Sprintf("%s. Загрузка...", tableTitle))
	})

	data, err := core.App.Client.GetData(ctx)
	if err != nil {
		core.App.Client.Logger.WriteLog(err.Error())
		return
	}

	meta := data.Meta
	manga := data.Manga

	if meta.From == 0 {
		ShowModal(utils.NoMangaID, "Манга не найдена")
		if core.App.Client.Page == 1 {
			core.App.Client.Query = ""
		} else {
			core.App.Client.Page--
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

	core.App.TView.QueueUpdateDraw(func() {
		p.table.Select(0, 0)
		p.table.ScrollToBeginning()
	})
}

func (p *ListPage) getMangaFromCell(row int) *models.Manga {
	return p.table.GetCell(row, 0).GetReference().(*models.Manga)
}
