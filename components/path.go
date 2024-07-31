package components

import (
	"path/filepath"
	"strings"

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
		LogsPath:     core.App.Client.Logger.Path,
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

	dInput := tview.NewInputField()
	dInput.SetLabel(utils.PathDownloadLabel).SetText(core.App.Client.DownloadPath)

	lInput := tview.NewInputField()
	lInput.SetLabel(utils.PathLogsLabel).SetText(core.App.Client.Logger.Path)

	dInput.SetAutocompleteFunc(getMatches)
	lInput.SetAutocompleteFunc(getMatches)

	form := tview.NewForm()
	form.SetBorder(true).SetTitle("Установить пути")
	form.SetButtonsAlign(tview.AlignCenter)
	form.AddFormItem(dInput).AddFormItem(lInput).
		AddButton("OK", func() {
			downloadPath := dInput.GetText()
			logPath := lInput.GetText()

			core.App.Client.ChangePath(downloadPath)
			core.App.Client.Logger.ChangePath(logPath)

			core.App.SaveConfig()

			core.App.PageHolder.RemovePage(utils.PathsModalID)
		}).
		AddButton("Default", func() {
			core.App.DefaultConfig()
			core.App.PageHolder.RemovePage(utils.PathsModalID)
		}).
		AddButton("Cancel", func() {
			core.App.PageHolder.RemovePage(utils.PathsModalID)
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
