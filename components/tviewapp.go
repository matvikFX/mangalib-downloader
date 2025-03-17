package components

import (
	"manga-downloader/downloader"
	"manga-downloader/services"

	"github.com/rivo/tview"
)

type TViewApp struct {
	Downloader *downloader.MangaLibDownloader

	Config    *services.Config
	Bookmarks *services.Bookmarks
	Logger    *services.Logger

	App   *tview.Application
	Pages *tview.Pages
}

func NewTViewApp() *TViewApp {
	return &TViewApp{
		Downloader: downloader.NewClient(),

		App:   tview.NewApplication(),
		Pages: tview.NewPages(),
	}
}

func (a *TViewApp) Start() error {
	if err := a.Bookmarks.Load(); err != nil {
		a.Logger.Write("No bookmarks file detected")
	}

	a.ShowListPage()
	a.SetHandlers()

	return nil
}

func (a *TViewApp) Stop() {
	a.App.Sync()
	a.App.Stop()
}
