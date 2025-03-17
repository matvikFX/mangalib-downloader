package components

import (
	"context"

	"manga-downloader/components/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *TViewApp) ShowBranchModal(ctx context.Context) {
	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewGrid().
			SetColumns(0, width, 0).SetRows(0, height, 0).
			AddItem(p, 1, 1, 1, 1, 0, 0, true)
	}

	branches, err := a.Downloader.GetMangaBranches(ctx, selectedManga.ID)
	if err != nil {
		a.Logger.Write(err.Error())
		return
	}

	if len(branches) == 0 {
		a.ShowMangaPage(ctx)
		return
	}

	selectedManga.Branches = branches
	teamsBranch := branches.BranchTeams()
	form := a.newBranchForm(ctx, teamsBranch)

	a.App.SetFocus(form)
	a.Pages.AddPage(utils.BranchModalID, modal(form, 50, 5), true, true)
}

func (a *TViewApp) newBranchForm(ctx context.Context, teamsBranch map[int]string) *tview.Form {
	form := tview.NewForm()
	form.SetTitle("Выбор ветки переводчиков").SetBorder(true)

	dropDown := tview.NewDropDown().SetLabel(utils.BranchModalLabel)

	for branch, team := range teamsBranch {
		dropDown.AddOption(team, func() {
			a.Downloader.Branch = branch
		})
	}

	dropDown.SetCurrentOption(0)
	form.AddFormItem(dropDown)

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			a.Pages.RemovePage(utils.BranchModalID)
		case tcell.KeyEnter:
			a.ShowMangaPage(ctx)
			a.Pages.RemovePage(utils.BranchModalID)
		}
		return event
	})

	return form
}
