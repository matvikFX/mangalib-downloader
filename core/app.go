package core

import (
	"log"

	"mangalib-downloader/client"

	"github.com/rivo/tview"
)

var App *MangaApp

type MangaApp struct {
	Client *client.MangaLibClient

	TView      *tview.Application
	PageHolder *tview.Pages
}

func NewApp() *MangaApp {
	return &MangaApp{
		Client:     client.NewClient(),
		TView:      tview.NewApplication(),
		PageHolder: tview.NewPages(),
	}
}

func (m *MangaApp) Init() {
	log.Println("Initializing app")
	m.TView.SetRoot(m.PageHolder, true).SetFocus(m.PageHolder)
}

func (m *MangaApp) Close() {
	m.TView.Sync()
	m.TView.Stop()

	log.Println("Application successfully closed")
}
