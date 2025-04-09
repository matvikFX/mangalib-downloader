package components

import (
	"context"

	"mangalib-downloader/components/utils"
	"mangalib-downloader/models"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *TViewApp) ShowBranchModal(
	ctx context.Context, slug string, branches models.BranchList,
) {
	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewGrid().
			SetColumns(0, width, 0).SetRows(0, height, 0).
			AddItem(p, 1, 1, 1, 1, 0, 0, true)
	}

	form := a.newBranchForm(ctx, slug, branches)

	a.App.SetFocus(form)
	a.Pages.AddPage(utils.BranchModalID, modal(form, 50, 5), true, true)
}

func (a *TViewApp) newBranchForm(
	ctx context.Context, slug string, branches models.BranchList,
) *tview.Form {
	form := tview.NewForm()
	form.SetTitle("Выбор ветки переводчиков").SetBorder(true)

	dropDown := tview.NewDropDown().SetLabel(utils.BranchModalLabel)
	for branch, team := range branches.BranchTeams() {
		dropDown.AddOption(team, func() {
			a.branchID = branch
		})
	}

	dropDown.SetCurrentOption(0)
	form.AddFormItem(dropDown)

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			a.query = ""
			a.Pages.RemovePage(utils.BranchModalID)
		case tcell.KeyEnter:
			a.ShowMangaPage(ctx, slug)
			a.Pages.RemovePage(utils.BranchModalID)
		}
		return event
	})

	return form
}

func (a *TViewApp) GetTeamsByID(branchList models.BranchList) string {
	return branchList.BranchTeams()[a.branchID]
}
