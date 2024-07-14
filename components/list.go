package components

import (
	"context"
	"fmt"

	"mangalib-downloader/components/utils"
	"mangalib-downloader/core"

	"github.com/rivo/tview"
)

type ListPage struct {
	table *tview.Table

	cWrap *utils.ContextWrapper
}

func ShowListPage() {
	listPage := newListPage()

	core.App.TView.SetFocus(listPage.table)
	core.App.PageHolder.AddAndSwitchToPage(utils.ListPageID, listPage.table, true)
}

func newListPage() *ListPage {
	table := newListTable()

	ctx, cancel := context.WithCancel(context.Background())
	listPage := &ListPage{
		table: table,
		cWrap: &utils.ContextWrapper{
			Context: ctx,
			Cancel:  cancel,
		},
	}

	go listPage.setListTable()

	return listPage
}

func newListTable() *tview.Table {
	table := tview.NewTable()

	table.SetSelectable(true, false).
		SetSeparator('|').
		SetTitle("Список манги").
		SetBorder(true)

	return table
}

func (p *ListPage) setListTable() {
	ctx, cancel := p.cWrap.ResetContext()
	defer cancel()
	p.setHandlers(cancel)

	tableTitle := "Популярная манга"
	if core.App.Client.Query != "" {
		tableTitle = "Результаты поиска"
	}

	core.App.TView.QueueUpdateDraw(func() {
		title := tview.NewTableCell("Название").
			SetAlign(tview.AlignCenter).
			SetSelectable(false)

		p.table.SetCell(0, 0, title).
			SetFixed(1, 0)

		p.table.SetTitle(fmt.Sprintf("%s. Загрузка...", tableTitle))
	})

	if p.cWrap.ToCancel(ctx) {
		Logger.WriteLog(ctx.Err().Error())
		return
	}

	data, err := core.App.Client.GetData(ctx)
	if err != nil {
		Logger.WriteLog(err.Error())
		return
	}

	meta := data.Meta
	manga := data.Manga

	if meta.From == 0 {
		// Вывести модалку
		Logger.WriteLog("манга не найдена")
		core.App.Client.Query = ""
		go p.setListTable()
	}

	p.table.SetTitle(fmt.Sprintf("%s. Страница %d (%d-%d)",
		tableTitle, meta.Page, meta.From, meta.To))

	for idx, manga := range manga {
		if p.cWrap.ToCancel(ctx) {
			Logger.WriteLog(ctx.Err().Error())
			return
		}

		manga.RusNameChange()
		title := tview.NewTableCell(manga.RusName).
			SetReference(manga)

		p.table.SetCell(idx+1, 0, title)
	}

	core.App.TView.QueueUpdateDraw(func() {
		p.table.Select(1, 0)
		p.table.ScrollToBeginning()
	})
}
