package components

import (
	"path/filepath"
	"strings"

	"manga-downloader/components/utils"

	"github.com/rivo/tview"
)

type PathModal struct {
	app *TViewApp

	form  *tview.Form
	modal tview.Primitive
}

func (a *TViewApp) ShowPathModal() {
	pathsModal := newPathModal()
	pathsModal.setHandlers()

	a.App.SetFocus(pathsModal.form)
	a.Pages.AddPage(utils.PathsModalID, pathsModal.modal, true, true)
}

func newPathModal() *PathModal {
	pathModal := &PathModal{}
	pathModal.setForm()

	return pathModal
}

func (p *PathModal) setForm() {
	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewGrid().
			SetColumns(0, width, 0).SetRows(0, height, 0).
			AddItem(p, 1, 1, 1, 1, 0, 0, true)
	}

	dInput := tview.NewInputField()
	dInput.SetLabel(utils.PathDownloadLabel).
		SetText(p.app.Config.DownloadPath)
	dInput.SetAutocompleteFunc(getMatches)

	lInput := tview.NewInputField()
	lInput.SetLabel(utils.PathLogsLabel).
		SetText(p.app.Logger.Path)
	lInput.SetAutocompleteFunc(getMatches)

	form := tview.NewForm()
	form.SetBorder(true).SetTitle("Установить пути")
	form.SetButtonsAlign(tview.AlignCenter)
	form.AddFormItem(dInput).AddFormItem(lInput).
		AddButton("OK", func() {
			// downloadPath := dInput.GetText()
			// logPath := lInput.GetText()

			// if msg := p.app.Config.ChangePath(downloadPath); msg != "" {
			// 	// Show error message
			// 	p.app.ShowModal(utils.DownloaderPathID, msg)
			// }
			//
			// if msg := p.app.Config.ChangePath(logPath); msg != "" {
			// 	// Show error message
			// 	p.app.ShowModal(utils.LoggerPathID, msg)
			// }

			p.app.Config.Save()
			p.app.Pages.RemovePage(utils.PathsModalID)
		}).
		AddButton("Default", func() {
			p.app.Config.Default()
			p.app.Pages.RemovePage(utils.PathsModalID)
		}).
		AddButton("Cancel", func() {
			p.app.Pages.RemovePage(utils.PathsModalID)
		})

	p.form = form
	p.modal = modal(form, 100, 9)
}

func getMatches(currentText string) (entries []string) {
	const hintsNum = 10

	if len(currentText) == 0 {
		return nil
	}

	matchesWithPrefix, err := filepath.Glob(currentText + "*")
	if err != nil {
		return nil
	}

	if len(matchesWithPrefix) == 1 {
		if currentText == matchesWithPrefix[0] {
			return nil
		}
	}

	var matches []string
	for _, match := range matchesWithPrefix {
		dirs := strings.Split(match, "/")
		if strings.HasPrefix(dirs[len(dirs)-1], ".") {
			continue
		}
		matches = append(matches, match)
	}

	if len(matches) > hintsNum {
		matches = matches[:hintsNum]
	}

	return matches
}
