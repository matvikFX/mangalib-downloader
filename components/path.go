package components

import (
	"mangalib-downloader/components/utils"
	"mangalib-downloader/core"

	"github.com/rivo/tview"
)

type PathModal struct {
	DownloadPath string
	LogsPath     string

	form  *tview.Form
	modal tview.Primitive
}

func ShowPathModal() {
	pathsModal := newPathModal()
	pathsModal.setHandlers()

	core.App.TView.SetFocus(pathsModal.form)
	core.App.PageHolder.AddPage(utils.PathsModalID, pathsModal.modal, true, true)
}

func newPathModal() *PathModal {
	pathModal := &PathModal{
		DownloadPath: core.App.Client.DownloadPath,
		LogsPath:     Logger.Path,
	}
	pathModal.setForm()

	return pathModal
}

func (p *PathModal) setForm() {
	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewGrid().
			SetColumns(0, width, 0).SetRows(0, height, 0).
			AddItem(p, 1, 1, 1, 1, 0, 0, true)
	}

	form := tview.NewForm()
	form.SetButtonsAlign(tview.AlignCenter).
		SetTitle("Указание путей").
		SetBorder(true)
	form.AddInputField(utils.PathDownloadLabel, core.App.Client.DownloadPath, 68, nil, func(text string) {
		p.DownloadPath = text
	}).AddInputField(utils.PathLogsLabel, Logger.Path, 68, nil, func(text string) {
		p.LogsPath = text
	}).
		AddButton("Ok", p.changePaths).
		AddButton("Cancel", p.closePage)

	p.form = form
	p.modal = modal(form, 100, 9)
}

func (p *PathModal) changePaths() {
	core.App.Client.ChangePath(p.DownloadPath)
	Logger.ChangePath(p.LogsPath)
}

func (p *PathModal) closePage() {
	core.App.PageHolder.RemovePage(utils.PathsModalID)
}
