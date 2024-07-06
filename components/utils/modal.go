package utils

import (
	"mangalib-downlaoder/core"

	"github.com/rivo/tview"
)

func ShowModal(id string, modal *tview.Modal) {
	core.App.TView.SetFocus(modal)
	core.App.PageHolder.AddPage(id, modal, true, true)
}

func newModal(id, text string) *tview.Modal {
	modal := tview.NewModal()
	modal.SetText(text).
		AddButtons([]string{"Ok"}).
		SetFocus(0).
		SetDoneFunc(func(_ int, _ string) {
			core.App.PageHolder.RemovePage(id)
		})

	return modal
}
