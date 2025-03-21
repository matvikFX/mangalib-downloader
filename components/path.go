package components

import (
	"log/slog"
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
	pathsModal := newPathModal(a)
	pathsModal.setHandlers()

	a.App.SetFocus(pathsModal.form)
	a.Pages.AddPage(utils.PathsModalID, pathsModal.modal, true, true)
}

func newPathModal(app *TViewApp) *PathModal {
	pathModal := &PathModal{
		app: app,
	}
	pathModal.setForm()

	return pathModal
}

func (p *PathModal) setForm() {
	log := p.app.Logger.With("PathModal", "setForm")

	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewGrid().
			SetColumns(0, width, 0).SetRows(0, height, 0).
			AddItem(p, 1, 1, 1, 1, 0, 0, true)
	}

	log.Debug("Variables", slog.Group("Paths",
		"DownloadPath", p.app.Config.DownloadPath,
		"LoggerPath", p.app.Config.LoggerPath,
		"BookmarksPath", p.app.Config.BookmarksPath,
	))

	dInput := tview.NewInputField()
	dInput.SetLabel(utils.PathDownloadLabel).
		SetText(p.app.Config.DownloadPath)
	dInput.SetAutocompleteFunc(getMatches)

	lInput := tview.NewInputField()
	lInput.SetLabel(utils.PathLogsLabel).
		SetText(p.app.Config.LoggerPath)
	lInput.SetAutocompleteFunc(getMatches)

	bInput := tview.NewInputField()
	bInput.SetLabel(utils.PathBookmarksLabel).
		SetText(p.app.Config.BookmarksPath)
	bInput.SetAutocompleteFunc(getMatches)

	form := tview.NewForm()
	form.SetBorder(true).SetTitle("Установить пути")
	form.SetButtonsAlign(tview.AlignCenter)
	form.AddFormItem(dInput).AddFormItem(lInput).AddFormItem(bInput).
		AddButton("OK", func() {
			downloadPath := dInput.GetText()
			logPath := lInput.GetText()
			bookmarksPath := bInput.GetText()

			if msg := p.app.Config.ChangeDownloadPath(downloadPath); msg != "" {
				// Show error message
				p.app.ShowModal(utils.DownloaderPathID, msg)
			}

			if msg := p.app.Config.ChangeLogPath(logPath); msg != "" {
				// Show error message
				p.app.ShowModal(utils.LoggerPathID, msg)
			}

			if msg := p.app.Config.ChangeBookmarkPath(bookmarksPath); msg != "" {
				// Show error message
				p.app.ShowModal(utils.BookmarksPathID, msg)
			}

			p.app.Config.Save()
			p.app.Pages.RemovePage(utils.PathsModalID)
		}).
		AddButton("Default", func() {
			p.app.Config.Default()
			p.app.Pages.RemovePage(utils.PathsModalID)
		}).
		AddButton("Cancel", func() {
			p.app.Pages.RemovePage(utils.PathsModalID)
		}).
		AddCheckbox("cbz format", true, func(checked bool) {
			p.app.Config.CbzFormat = checked
		})

	p.form = form
	// old value: 9
	p.modal = modal(form, 100, 13)
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
