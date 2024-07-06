package components

import (
	"context"

	"mangalib-downlaoder/components/utils"
	"mangalib-downlaoder/core"
	"mangalib-downlaoder/models"

	"github.com/rivo/tview"
)

type MangaPage struct {
	manga *models.MangaInfo

	selected map[int]bool

	Grid     *tview.Grid
	TextView *tview.TextView
	Table    *tview.Table

	ctxWrap *utils.ContextWrapper
}

func ShowMangaPage(manga *models.MangaInfo) {
	mangaPage := newMangaPage(manga)

	core.App.TView.SetFocus(mangaPage.Grid)
	core.App.PageHolder.AddAndSwitchToPage(utils.MangaPageID, mangaPage.Grid, true)
}

func newMangaPage(manga *models.MangaInfo) *MangaPage {
	textView := tview.NewTextView()
	textView.SetWrap(true).SetWordWrap(true).
		SetTitle("Информация").SetBorder(true)

	table := newTable()

	grid := tview.NewGrid()
	grid.SetRows(-1).SetColumns(-1, -1, -1, -1, -1, -1, -1, -1, -1).
		SetTitle("Информация о главе").SetBorder(true)

	grid.AddItem(textView, 0, 0, 1, 3, 0, 0, false).
		AddItem(table, 0, 3, 1, 6, 0, 0, true)

	ctx, cancel := context.WithCancel(context.Background())
	mangaPage := &MangaPage{
		manga: manga,

		Grid:     grid,
		TextView: textView,
		Table:    table,

		ctxWrap: &utils.ContextWrapper{
			Context: ctx,
			Cancel:  cancel,
		},
	}

	go mangaPage.setMangaInfo()

	return mangaPage
}

func newTable() *tview.Table {
	table := tview.NewTable()

	vol := tview.NewTableCell("Том").
		SetMaxWidth(3).
		SetAlign(tview.AlignCenter).
		SetSelectable(false)
	num := tview.NewTableCell("Том").
		SetMaxWidth(3).
		SetAlign(tview.AlignCenter).
		SetSelectable(false)
	name := tview.NewTableCell("Том").
		SetMaxWidth(3).
		SetAlign(tview.AlignCenter).
		SetSelectable(false)
	downloadStatus := tview.NewTableCell("Том").
		SetMaxWidth(3).
		SetAlign(tview.AlignCenter).
		SetSelectable(false)

	table.SetCell(0, 0, vol).
		SetCell(0, 1, num).
		SetCell(0, 2, name).
		SetCell(0, 3, downloadStatus).
		SetFixed(0, 1)

	table.SetSelectable(true, false).
		SetSeparator('|').
		SetTitle("Главы").
		SetBorder(true)

	return table
}

func (p *MangaPage) setMangaInfo() {
	info := utils.InfoText(p.manga)

	core.App.TView.QueueUpdateDraw(func() {
		p.TextView.SetText(info)
	})
}
