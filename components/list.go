package components

import (
	"context"
	"errors"
	"fmt"

	"mangalib-downlaoder/components/utils"
	"mangalib-downlaoder/core"
	"mangalib-downlaoder/models"

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

		// desc := tview.NewTableCell("Описание").
		// 	SetAlign(tview.AlignCenter).
		// 	SetSelectable(false)

		// tags_genres := tview.NewTableCell("Теги и Жанры").
		// 	SetAlign(tview.AlignCenter).
		// 	SetSelectable(false)

		p.table.SetCell(0, 0, title).
			// SetCell(0, 1, desc).
			// SetCell(0, 1, tags_genres).
			SetFixed(1, 0)

		p.table.SetTitle(fmt.Sprintf("%s. Загрузка...", tableTitle))
	})

	if p.cWrap.ToCancel(ctx) {
		Logger.WriteLog(ctx.Err().Error())
		return
	}

	// list, meta, err := getInfoList(ctx)
	// if err != nil {
	// 	Logger.WriteLog(err.Error())
	// 	return
	// }

	data, err := core.App.Client.GetData(ctx)
	if err != nil {
		Logger.WriteLog(err.Error())
		return
	}

	meta := data.Meta
	manga := data.Manga

	if data.Meta.From == 0 {
		// Вывести модалку
		Logger.WriteLog("манга не найдена")
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
		// titleText := fmt.Sprintf("%-60s", manga.RusName)
		title := tview.NewTableCell(manga.RusName).
			// SetMaxWidth(60).
			SetReference(manga)

		// tagsText := strings.Join(manga.GetTags(), ",")
		// tags := tview.NewTableCell(tagsText).
		// 	SetMaxWitdth(60)

		// descText := fmt.Sprintf("%-140s", manga.Description)
		// desc := tview.NewTableCell(descText).
		// 	SetMaxWidth(140).SetReference(manga)

		p.table.SetCell(idx+1, 0, title)
		// SetCell(idx+1, 1, t_g).
		// SetCell(idx+1, 2, desc)
	}

	core.App.TView.QueueUpdateDraw(func() {
		p.table.Select(1, 0)
		p.table.ScrollToBeginning()
	})
}

func getInfoList(ctx context.Context) (models.MangaInfoList, *models.Meta, error) {
	data, err := core.App.Client.GetData(ctx)
	if err != nil {
		return nil, nil, err
	}

	meta := data.Meta
	mangaList := data.Manga

	// mangaCount := meta.To - meta.From + 1
	if meta.From == 0 {
		return nil, nil, errors.New("манга не найдена")
	}

	var slugs []string
	for _, manga := range mangaList {
		slugs = append(slugs, manga.Slug)
	}

	list := core.App.Client.GetInfoList(ctx, slugs)
	return list, meta, nil
}

func (p *ListPage) setSelected(row, column int) {
	manga := p.table.GetCell(row, 0).GetReference().(*models.Manga)
	info, err := core.App.Client.GetInfo(p.cWrap.Context, manga.Slug)
	if err != nil {
		Logger.WriteLog(err.Error())
		return
	}

	ShowMangaPage(info)
}
