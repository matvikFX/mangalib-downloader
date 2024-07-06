package components

import (
	"context"
	"fmt"

	"mangalib-downlaoder/components/utils"
	"mangalib-downlaoder/core"
	"mangalib-downlaoder/models"

	"github.com/rivo/tview"
)

type MangaPage struct {
	manga *models.MangaInfo

	selected map[int]*models.Chapter

	gird     *tview.Grid
	textView *tview.TextView
	table    *tview.Table

	cWrap *utils.ContextWrapper
}

func ShowMangaPage(manga *models.MangaInfo) {
	mangaPage := newMangaPage(manga)

	core.App.TView.SetFocus(mangaPage.gird)
	core.App.PageHolder.AddAndSwitchToPage(utils.MangaPageID, mangaPage.gird, true)
}

func newMangaPage(manga *models.MangaInfo) *MangaPage {
	textView := tview.NewTextView()
	textView.SetWrap(true).SetWordWrap(true).
		SetTitle("Информация").SetBorder(true)

	table := newInfoTable()

	grid := tview.NewGrid()
	grid.SetRows(-1).SetColumns(-1, -1, -1, -1, -1, -1, -1, -1, -1).
		SetTitle("Информация о главе").SetBorder(true)

	grid.AddItem(textView, 0, 0, 1, 3, 0, 0, false).
		AddItem(table, 0, 3, 1, 6, 0, 0, true)

	ctx, cancel := context.WithCancel(context.Background())
	mangaPage := &MangaPage{
		manga: manga,

		gird:     grid,
		textView: textView,
		table:    table,

		cWrap: &utils.ContextWrapper{
			Context: ctx,
			Cancel:  cancel,
		},
	}

	go mangaPage.setMangaInfo()
	go mangaPage.setChapters()

	return mangaPage
}

func newInfoTable() *tview.Table {
	table := tview.NewTable()

	vol := tview.NewTableCell("Том").
		SetMaxWidth(3).
		SetAlign(tview.AlignCenter).
		SetSelectable(false)
	num := tview.NewTableCell("Номер").
		SetMaxWidth(5).
		SetAlign(tview.AlignCenter).
		SetSelectable(false)
	name := tview.NewTableCell("Название").
		SetMaxWidth(40).
		SetAlign(tview.AlignCenter).
		SetSelectable(false)
	downloadStatus := tview.NewTableCell("Состояние загрузки").
		SetAlign(tview.AlignCenter).
		SetSelectable(false)

	table.SetCell(0, 0, vol).
		SetCell(0, 1, num).
		SetCell(0, 2, name).
		SetCell(0, 3, downloadStatus).
		SetFixed(1, 0)

	table.SetSelectable(true, false).
		SetSeparator('|').
		SetTitle("Главы").
		SetBorder(true)

	return table
}

func (p *MangaPage) setMangaInfo() {
	info := utils.InfoText(p.manga)

	core.App.TView.QueueUpdateDraw(func() {
		p.textView.SetText(info)
	})
}

func (p *MangaPage) setChapters() {
	ctx, cancel := p.cWrap.ResetContext()
	defer cancel()
	p.setHandlers(cancel)

	core.App.TView.QueueUpdateDraw(func() {
		loading := tview.NewTableCell("Загрузка...").SetSelectable(false)
		p.table.SetCell(1, 1, loading)
		p.table.SetTitle("Загрузка глав...")
	})

	if p.cWrap.ToCancel(ctx) {
		Logger.WriteLog(ctx.Err().Error())
		return
	}

	chaps, err := core.App.Client.GetChapters(ctx, p.manga.Slug)
	if err != nil {
		Logger.WriteLog(err.Error())
		return
	}

	if len(chaps) == 0 {
		core.App.TView.QueueUpdateDraw(func() {
			// noRes := tview.NewTableCell("Не удалось найти ни одну главу").SetSelectable(false)
			// p.table.SetCell(1, 1, noRes)
			p.table.SetTitle("Не удалось найти ни одну главу")
		})
		return
	}

	p.table.SetTitle("Главы")
	for idx, ch := range chaps {
		if p.cWrap.ToCancel(ctx) {
			Logger.WriteLog(ctx.Err().Error())
			return
		}

		vol := tview.NewTableCell(
			fmt.Sprintf("%-3s", ch.Volume)).
			SetMaxWidth(5).SetReference(ch)
		num := tview.NewTableCell(
			fmt.Sprintf("%-5s", ch.Number)).
			SetMaxWidth(5).SetReference(ch)
		name := tview.NewTableCell(
			fmt.Sprintf("%-40s", ch.Name)).
			SetMaxWidth(40).SetReference(ch)

		var downloadStatus string
		chapPath := core.App.Client.CreateChapterPath("", p.manga.RusName, ch.Volume, ch.Number, ch.Name)
		if core.App.Client.CheckExistence(chapPath) {
			downloadStatus = "X"
		}
		download := tview.NewTableCell(downloadStatus)

		p.table.SetCell(idx+1, 0, vol)
		p.table.SetCell(idx+1, 1, num)
		p.table.SetCell(idx+1, 2, name)
		p.table.SetCell(idx+1, 3, download)
	}

	core.App.TView.QueueUpdateDraw(func() {
		p.table.Select(1, 0)
		p.table.ScrollToBeginning()
	})
}
