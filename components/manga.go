package components

import (
	"context"
	"fmt"

	"manga-downloader/api"
	"manga-downloader/components/utils"
	"manga-downloader/downloader"
	"manga-downloader/models"

	"github.com/rivo/tview"
)

type MangaPage struct {
	app        *TViewApp
	downloader *downloader.MangaLibDownloader

	selected map[int]bool

	grid     *tview.Grid
	textView *tview.TextView
	table    *tview.Table
}

func (a *TViewApp) ShowMangaPage(ctx context.Context) {
	if selectedManga.Description == "" {
		info, err := api.GetInfo(ctx, selectedManga.Slug, a.branchID)
		if err != nil {
			a.Logger.Write(err.Error())
			return
		}

		if len(selectedManga.Branches) != 0 {
			info.Branches = selectedManga.Branches
		}

		selectedManga = info
	}

	mangaPage := newMangaPage(ctx, a)
	a.App.SetFocus(mangaPage.grid)
	a.Pages.AddAndSwitchToPage(utils.MangaPageID, mangaPage.grid, true)
}

func newMangaPage(ctx context.Context, app *TViewApp) *MangaPage {
	textView := tview.NewTextView()
	textView.SetWrap(true).SetWordWrap(true).
		SetTitle("Информация").SetBorder(true)

	table := newInfoTable()

	grid := tview.NewGrid()
	grid.SetRows(-1).SetColumns(-1, -1, -1, -1, -1, -1, -1, -1, -1).
		SetTitle("Информация о главе").SetBorder(true)

	grid.AddItem(textView, 0, 0, 1, 3, 0, 0, false).
		AddItem(table, 0, 3, 1, 6, 0, 0, true)

	mangaPage := &MangaPage{
		app:        app,
		downloader: downloader.NewClient(),

		selected: make(map[int]bool),

		grid:     grid,
		textView: textView,
		table:    table,
	}

	go mangaPage.setMangaInfo()
	go mangaPage.setChapters(ctx)

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
	teams := selectedManga.Branches.BranchTeamList()
	info := utils.InfoText(selectedManga, teams[p.app.branchID])

	p.app.App.QueueUpdateDraw(func() {
		p.textView.SetText(info)
	})
}

func (p *MangaPage) setChapters(parentCtx context.Context) {
	ctx, cancel := context.WithCancel(parentCtx)
	p.setHandlers(ctx, cancel)

	p.app.App.QueueUpdateDraw(func() {
		loading := tview.NewTableCell("Загрузка...").SetSelectable(false)
		p.table.SetCell(1, 2, loading)
		p.table.SetTitle("Загрузка глав...")
	})

	chaps, err := api.GetChapters(ctx, selectedManga.Slug, p.app.branchID)
	if err != nil {
		p.app.Logger.Write(err.Error())
		return
	}

	if len(chaps) == 0 {
		p.app.App.QueueUpdateDraw(func() {
			noRes := tview.NewTableCell("Не удалось найти ни одну главу").
				SetSelectable(false)
			p.table.SetCell(1, 2, noRes)
		})
		return
	}

	// branchTeams := selectedManga.Branches.BranchTeams()[p.app.branchID]
	p.table.SetTitle("Главы")
	for idx, ch := range chaps {
		vol := tview.NewTableCell(
			fmt.Sprintf("%-3s", ch.Volume)).
			SetMaxWidth(5).SetReference(ch)
		num := tview.NewTableCell(
			fmt.Sprintf("%-5s", ch.Number)).
			SetMaxWidth(5).SetReference(ch)
		name := tview.NewTableCell(
			fmt.Sprintf("%-40s", ch.Name)).
			SetMaxWidth(40).SetReference(ch)

		// Пока в разработке, так скажем
		//
		// var downloadStatus string
		// chapPath := downloader.CreateChapterPath(
		// 	downloader.Path,
		// 	branchTeams, selectedManga.RusName,
		// 	ch.Volume, ch.Number, ch.Name)
		// if downloader.CheckExistence(chapPath) {
		// 	downloadStatus = "X"
		// }
		// status := tview.NewTableCell(downloadStatus)

		p.table.SetCell(idx+1, 0, vol)
		p.table.SetCell(idx+1, 1, num)
		p.table.SetCell(idx+1, 2, name)
		// p.table.SetCell(idx+1, 3, status)
	}

	p.app.App.QueueUpdateDraw(func() {
		p.table.Select(1, 0)
		p.table.ScrollToBeginning()
	})
}

func (p *MangaPage) downloadSelected(ctx context.Context, branchID int) {
	var chaps models.ChapterList
	for row, selected := range p.selected {
		if !selected {
			continue
		}

		chap := p.table.GetCell(row, 0).GetReference().(*models.Chapter)
		if chap == nil {
			return
		}
		chaps = append(chaps, chap)
	}

	p.downloader.DownloadChapters(ctx, selectedManga.Manga, chaps, branchID)

	<-p.downloader.Downloaded
	p.app.ShowModal(utils.DownloadSuccessID,
		"Выбранные главы манги '"+selectedManga.RusName+"' успешно скачаны")

	go p.setChapters(ctx)
}
