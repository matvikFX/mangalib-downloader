package components

import (
	"context"
	"fmt"

	"mangalib-downloader/components/utils"
	"mangalib-downloader/core"
	"mangalib-downloader/models"

	"github.com/rivo/tview"
)

type MangaPage struct {
	selected map[int]bool

	grid     *tview.Grid
	textView *tview.TextView
	table    *tview.Table
}

func ShowMangaPage(ctx context.Context) {
	if selectedManga.Description == "" {
		info, err := core.App.Client.GetInfo(ctx, selectedManga.Slug)
		if err != nil {
			core.App.Client.Logger.WriteLog(err.Error())
			return
		}

		if len(selectedManga.Branches) != 0 {
			info.Branches = selectedManga.Branches
		}

		selectedManga = info
	}

	mangaPage := newMangaPage(ctx)
	core.App.TView.SetFocus(mangaPage.grid)
	core.App.PageHolder.AddAndSwitchToPage(utils.MangaPageID, mangaPage.grid, true)
}

func newMangaPage(ctx context.Context) *MangaPage {
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
	info := utils.InfoText(selectedManga, teams[core.App.Client.Branch])

	core.App.TView.QueueUpdateDraw(func() {
		p.textView.SetText(info)
	})
}

func (p *MangaPage) setChapters(parentCtx context.Context) {
	ctx, cancel := context.WithCancel(parentCtx)
	p.setHandlers(ctx, cancel)

	core.App.TView.QueueUpdateDraw(func() {
		loading := tview.NewTableCell("Загрузка...").SetSelectable(false)
		p.table.SetCell(1, 2, loading)
		p.table.SetTitle("Загрузка глав...")
	})

	chaps, err := core.App.Client.GetChapters(ctx, selectedManga.Slug)
	if err != nil {
		core.App.Client.Logger.WriteLog(err.Error())
		return
	}

	if len(chaps) == 0 {
		core.App.TView.QueueUpdateDraw(func() {
			noRes := tview.NewTableCell("Не удалось найти ни одну главу").
				SetSelectable(false)
			p.table.SetCell(1, 2, noRes)
		})
		return
	}

	branchTeams := selectedManga.Branches.BranchTeams()[core.App.Client.Branch]
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

		var downloadStatus string
		chapPath := core.App.Client.CreateChapterPath(
			branchTeams, selectedManga.RusName,
			ch.Volume, ch.Number, ch.Name)
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

func (p *MangaPage) downloadSelected(ctx context.Context) {
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

	core.App.Client.DownloadChapters(ctx, selectedManga.Manga, chaps)

	<-core.App.Client.Downloaded
	ShowModal(utils.DownloadSuccessID,
		"Выбранные главы манги '"+selectedManga.RusName+"' успешно скачаны")

	go p.setChapters(ctx)
}
