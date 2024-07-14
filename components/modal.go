package components

import (
	"mangalib-downloader/core"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func ShowModal(id, text string) {
	modal := tview.NewModal()
	modal.SetText(text).
		SetBackgroundColor(tcell.ColorBlack).
		AddButtons([]string{"Ok"}).
		SetFocus(0).
		SetDoneFunc(func(_ int, _ string) {
			core.App.PageHolder.RemovePage(id)
		})

	core.App.TView.SetFocus(modal)
	core.App.PageHolder.AddPage(id, modal, true, true)
}
