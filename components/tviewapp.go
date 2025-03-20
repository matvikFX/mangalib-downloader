package components

import (
	"manga-downloader/services"

	"github.com/rivo/tview"
)

type TViewApp struct {
	Config    *services.Config
	Bookmarks *services.Bookmarks
	Logger    *services.Logger

	App   *tview.Application
	Pages *tview.Pages

	query    string
	page     int
	branchID int
}

func NewTViewApp() *TViewApp {
	return &TViewApp{
		Config:    services.NewConfig(),
		Bookmarks: services.NewBookmarks(),
		Logger:    services.NewLogger(),

		App:   tview.NewApplication(),
		Pages: tview.NewPages(),

		query:    "",
		page:     1,
		branchID: 0,
	}
}

func (a *TViewApp) Start() {
	if err := a.Bookmarks.Load(a.Config.BookmarksPath); err != nil {
		a.Logger.Write("No bookmarks file detected")
	}

	a.ShowListPage()
	a.SetHandlers()

	if err := a.App.SetRoot(a.Pages, true).Run(); err != nil {
		// panic(err)
		a.Logger.Write(err)
	}
}

func (a *TViewApp) Stop() {
	a.App.Sync()
	a.App.Stop()
}
