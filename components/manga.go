package components

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"mangalib-downloader/api"
	"mangalib-downloader/components/utils"
	"mangalib-downloader/downloader"
	"mangalib-downloader/models"

	"github.com/rivo/tview"
)

type MangaPage struct {
	app *TViewApp

	selected map[int]struct{}

	grid     *tview.Grid
	textView *tview.TextView
	table    *tview.Table
}

func (a *TViewApp) ShowMangaPage(ctx context.Context, slug string) {
	mangaPage, err := newMangaPage(ctx, a, slug)
	if err != nil {
		// Show modal
		// Не удалось загрузить страницу произведения
		slog.Error("Error creating manga page", "Error", err)
	}

	a.App.SetFocus(mangaPage.grid)
	a.Pages.AddAndSwitchToPage(utils.MangaPageID, mangaPage.grid, true)
}

func newMangaPage(
	ctx context.Context, app *TViewApp,
	slug string,
) (*MangaPage, error) {
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
		app: app,

		selected: make(map[int]struct{}),

		grid:     grid,
		textView: textView,
		table:    table,
	}

	manga, err := api.GetInfo(ctx, slug, app.branchID)
	if err != nil {
		slog.Error("Error receiving manga info", "Error", err)
		return nil, err
	}

	teams := api.GetBranchTeams(ctx, app.branchID)
	slog.Debug("Branch teams", "Teams", teams)

	go mangaPage.setMangaInfo(manga, strings.Split(teams, ","))
	go mangaPage.setChapters(ctx, manga, teams)

	return mangaPage, nil
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

func (p *MangaPage) setMangaInfo(manga *models.MangaInfo, teams []string) {
	info := utils.InfoText(manga, teams)
	p.app.App.QueueUpdateDraw(func() {
		p.textView.SetText(info)
	})
}

func (p *MangaPage) setChapters(
	ctx context.Context, manga *models.MangaInfo, branchTeams string,
) {
	ctx, cancel := context.WithCancel(ctx)
	p.setHandlers(ctx, cancel, manga.Slug)

	p.app.App.QueueUpdateDraw(func() {
		loading := tview.NewTableCell("Загрузка...").SetSelectable(false)
		p.table.SetCell(1, 2, loading)
		p.table.SetTitle("Загрузка глав...")
	})

	chaps, err := api.GetChapters(ctx, manga.Slug, p.app.branchID)
	if err != nil {
		slog.Error(err.Error())
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

	// branchTeams := branchList.BranchTeams()[p.app.branchID]
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
		chapPath := downloader.CreateChapterPath(
			p.app.Config.DownloadPath,
			branchTeams, manga.RusName,
			ch.Volume, ch.Number, ch.Name)
		if downloader.CheckExistence(chapPath) {
			downloadStatus = "X"
		}
		status := tview.NewTableCell(downloadStatus)

		p.table.SetCell(idx+1, 0, vol)
		p.table.SetCell(idx+1, 1, num)
		p.table.SetCell(idx+1, 2, name)
		p.table.SetCell(idx+1, 3, status)
	}

	p.app.App.QueueUpdateDraw(func() {
		p.table.Select(1, 0)
		p.table.ScrollToBeginning()
	})
}

func (p *MangaPage) downloadSelected(ctx context.Context, manga *models.MangaInfo) {
	var chaps models.ChapterList
	for row := range p.selected {
		chap := p.table.GetCell(row, 0).GetReference().(*models.Chapter)
		if chap == nil {
			return
		}
		chaps = append(chaps, chap)
	}

	if err := p.app.downloader.DownloadChapters(
		ctx, &manga.Manga, chaps, p.app.branchID,
	); err != nil {
		p.app.ShowModal(utils.DownloadFailID, err.Error())
		return
	}

	p.app.ShowModal(utils.DownloadSuccessID,
		"Выбранные главы манги '"+manga.RusName+"' успешно скачаны")

	go p.setChapters(ctx, manga, p.app.GetTeamsByID(manga.Branches))
}
