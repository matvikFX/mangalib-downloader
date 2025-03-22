package components

import (
	"log/slog"
	"manga-downloader/downloader"
	"manga-downloader/models"
	"manga-downloader/services"

	"github.com/rivo/tview"
)

type TViewApp struct {
	Config    *services.Config
	Bookmarks *services.Bookmarks

	App   *tview.Application
	Pages *tview.Pages

	downloader *downloader.Downloader

	query         string
	page          int
	branchID      int
	selectedManga *models.MangaInfo
}

func NewTViewApp(
	cfg *services.Config,
	bookmarks *services.Bookmarks,
	downloader *downloader.Downloader,
) *TViewApp {
	return &TViewApp{
		Config:    cfg,
		Bookmarks: bookmarks,

		App:   tview.NewApplication(),
		Pages: tview.NewPages(),

		downloader: downloader,

		query:    "",
		page:     1,
		branchID: 0,
	}
}

func (a *TViewApp) Start() error {
	log := slog.With("TViewApp", "Start")

	if err := a.Bookmarks.Load(a.Config.BookmarksPath); err != nil {
		log.Error("No bookmarks file detected", "Error", err)
		return err
	}
	log.Debug("Bookmarks loaded", "Bookmarks", a.Bookmarks)

	a.ShowListPage()
	a.setHandlers()

	if err := a.App.SetRoot(a.Pages, true).Run(); err != nil {
		log.Error("Error running application", "Error", err)
		return err
	}

	return nil
}

func (a *TViewApp) Stop() {
	a.App.EnableMouse(false)

	a.App.Sync()
	a.App.Stop()
}
