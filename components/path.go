package components

import (
	"log/slog"
	"mangalib-downloader/components/utils"
	"mangalib-downloader/services"

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
	log := slog.With("PathModal", "setForm")

	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewGrid().
			SetColumns(0, width, 0).SetRows(0, height, 0).
			AddItem(p, 1, 1, 1, 1, 0, 0, true)
	}

	// log.Debug("Variables", slog.Group("Paths",
	// 	"DownloadPath", p.app.Config.DownloadPath,
	// 	"LoggerPath", p.app.Config.LoggerPath,
	// 	"BookmarksPath", p.app.Config.BookmarksPath,
	// 	"CbzFormat", p.app.Config.CbzFormat,
	// ))

	dInput := tview.NewInputField()
	dInput.SetLabel(utils.PathDownloadLabel).
		SetText(p.app.Config.DownloadPath)
	dInput.SetAutocompleteFunc(utils.GetMatches)

	lInput := tview.NewInputField()
	lInput.SetLabel(utils.PathLogsLabel).
		SetText(p.app.Config.LoggerPath)
	lInput.SetAutocompleteFunc(utils.GetMatches)

	bInput := tview.NewInputField()
	bInput.SetLabel(utils.PathBookmarksLabel).
		SetText(p.app.Config.BookmarksPath)
	bInput.SetAutocompleteFunc(utils.GetMatches)

	cbzCheckbox := tview.NewCheckbox()
	cbzCheckbox.SetLabel(utils.CBZFormatCheckbox).
		SetChecked(p.app.Config.CbzFormat).
		SetChangedFunc(func(checked bool) {
			p.app.Config.CbzFormat = checked
		})

	form := tview.NewForm()
	form.SetBorder(true).SetTitle("Установить пути")
	form.SetButtonsAlign(tview.AlignCenter)
	form.AddFormItem(dInput).AddFormItem(lInput).AddFormItem(bInput).AddFormItem(cbzCheckbox).
		AddButton("OK", func() {
			downloadPath := dInput.GetText()
			logPath := lInput.GetText()
			bookmarksPath := bInput.GetText()

			if msg := p.app.Config.ChangeDownloadPath(downloadPath); msg != "" {
				p.app.ShowModal(utils.DownloaderPathID, msg)
			}

			if msg := p.app.Config.ChangeLogPath(logPath); msg != "" {
				p.app.ShowModal(utils.LoggerPathID, msg)
			}

			if msg := p.app.Config.ChangeBookmarksPath(bookmarksPath); msg != "" {
				p.app.ShowModal(utils.BookmarksPathID, msg)
			}

			if err := p.app.downloader.ChangeConfig(
				downloadPath, p.app.Config.CbzFormat,
			); err != nil {
				msg := "Error changing downloader config"
				log.Error(msg, "Error", err)
				p.app.ShowModal(utils.BookmarksPathID, msg)
			}

			if err := p.app.Config.Save(); err != nil {
				msg := "Error saving config"
				log.Error(msg, "Error", err)
				p.app.ShowModal(utils.BookmarksPathID, msg)
			}

			log.Info("Config was successfully changed", "Config", p.app.Config)
			p.app.Pages.RemovePage(utils.PathsModalID)
		}).
		AddButton("Default", func() {
			defaultConfig, err := services.DefaultConfig()
			if err != nil {
				p.app.ShowModal(utils.BookmarksPathID, err.Error())
			}

			p.app.Config = defaultConfig
			if err := p.app.downloader.ChangeConfig(
				defaultConfig.DownloadPath, defaultConfig.CbzFormat,
			); err != nil {
				msg := "Error changing downloader config"
				log.Error(msg, "Error", err)
				p.app.ShowModal(utils.BookmarksPathID, msg)
			}

			log.Info("Config set to default", "Config", defaultConfig)
			p.app.Pages.RemovePage(utils.PathsModalID)
		}).
		AddButton("Cancel", func() {
			p.app.Pages.RemovePage(utils.PathsModalID)
		})

	p.form = form
	p.modal = modal(form, 100, 13)
}
