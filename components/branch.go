package components

import (
	"context"

	"mangalib-downloader/components/utils"
	"mangalib-downloader/core"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func ShowBranchModal(ctx context.Context) {
	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewGrid().
			SetColumns(0, width, 0).SetRows(0, height, 0).
			AddItem(p, 1, 1, 1, 1, 0, 0, true)
	}

	branches, err := core.App.Client.GetMangaBranches(ctx, selectedManga.ID)
	if err != nil {
		Logger.WriteLog(err.Error())
		return
	}

	if len(branches) == 0 {
		ShowMangaPage(ctx)
		return
	}

	Logger.WriteLog("У манги есть ветки перевода")

	selectedManga.Branches = branches
	teamsBranch := branches.BranchTeams()
	form := newBranchForm(ctx, teamsBranch)

	core.App.TView.SetFocus(form)
	core.App.PageHolder.AddPage(utils.BranchModalID, modal(form, 50, 5), true, true)
}

func newBranchForm(ctx context.Context, teamsBranch map[int]string) *tview.Form {
	form := tview.NewForm()
	form.SetTitle("Выбор ветки переводчиков").SetBorder(true)

	dropDown := tview.NewDropDown().SetLabel(utils.BranchModalLabel)

	for branch, team := range teamsBranch {
		dropDown.AddOption(team, func() {
			core.App.Client.Branch = branch
		})
	}

	dropDown.SetCurrentOption(0)
	form.AddFormItem(dropDown)

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			core.App.PageHolder.RemovePage(utils.BranchModalID)
		case tcell.KeyEnter:
			ShowMangaPage(ctx)
			core.App.PageHolder.RemovePage(utils.BranchModalID)
		}
		return event
	})

	return form
}
