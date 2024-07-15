package components

import (
	"context"

	"mangalib-downloader/components/utils"
	"mangalib-downloader/core"
	"mangalib-downloader/models"

	"github.com/rivo/tview"
)

type BranchModal struct {
	ctx context.Context

	id   int
	slug string

	form  *tview.Form
	modal tview.Primitive
}

func ShowBranchModal(ctx context.Context, slug string, id int) {
	br, err := core.App.Client.GetMangaBranches(ctx, id)
	if err != nil {
		Logger.WriteLog(err.Error())
		return
	}

	if len(br) == 0 {
		SwitchToMangaPage(ctx, slug, id)
		return
	}

	branchModal := newBranchModal(ctx, slug, id, br)
	branchModal.setHandlers()

	core.App.TView.SetFocus(branchModal.form)
	core.App.PageHolder.AddPage(utils.BranchModalID, branchModal.modal, true, true)
}

func newBranchModal(ctx context.Context, slug string, id int, branches models.BranchList) *BranchModal {
	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewGrid().
			SetColumns(0, width, 0).SetRows(0, height, 0).
			AddItem(p, 1, 1, 1, 1, 0, 0, true)
	}

	form := tview.NewForm()
	form.SetTitle("Выбор ветки переводчиков").SetBorder(true)

	branchModel := &BranchModal{
		ctx: ctx,

		form:  form,
		modal: modal(form, 100, 5),
	}
	branchModel.setDropDown(branches)

	return branchModel
}

func (p *BranchModal) setDropDown(branches models.BranchList) {
	tBr := branches.TeamsBranch()
	dropDown := tview.NewDropDown().SetLabel(utils.BranchModalLabel)

	for team, branch := range tBr {
		dropDown.AddOption(team, func() {
			core.App.Client.Branch = branch
		})
	}
	dropDown.SetCurrentOption(0)

	p.form.AddFormItem(dropDown)
}
